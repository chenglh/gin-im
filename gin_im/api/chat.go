package api

import (
	"IM/gin_im/global"
	"IM/gin_im/model"
	"IM/gin_im/service"
	"IM/gin_im/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 连接信息
type Node struct {
	Conn          *websocket.Conn // 连接资源
	Addr          string          // 客户端地址
	FirstTime     int64           // 首次连接时间
	HeartbeatTime int64           // 心跳时间
	LoginTime     int64           // 登录时间
	DataQueue     chan []byte     // 消息[并发转串行，io资源，存在竞争关系，正在写数据，又发送一个过来，数据会乱。]
	GroupSets     set.Interface   // 好友/群
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 客户端的映射关系
var clientMap map[int64]*Node = make(map[int64]*Node)

// 读写锁
var rwLock sync.RWMutex

// HTTP 协议升级为 WebSocket 协议
func Chat(ctx *gin.Context) {
	// 1 检验接入是否合法
	uid := ctx.Query("userid")
	token := ctx.Query("token")

	// 解析得到用户ID
	userId, _ := strconv.ParseInt(uid, 10, 64)

	// 1.1 判断token是否合法
	userService := service.UserService{}
	userInfo, err := userService.GetUserByUserId(userId)
	if err != nil || userInfo.UserId == 0 || userInfo.Token != token {
		util.RespFail(ctx, -1, "websocket验证不通过")
		return
	}

	// 2 升级webSocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		// 如果前端js报错，websocket: request origin not allowed by Upgrader.CheckOrigin 是这里的跨域问题
		return true
	}}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("升级websocket失败", err.Error())
		return
	}

	// 3 获取连接 conn
	nowTime := carbon.Now().Timestamp()
	node := &Node{
		Conn:          conn,                       // 客户端连接
		Addr:          conn.RemoteAddr().String(), // 客户端地址
		HeartbeatTime: nowTime,                    // 心跳时间
		LoginTime:     nowTime,                    // 登录时间
		DataQueue:     make(chan []byte, 50),      // 消息队列
		GroupSets:     set.New(set.ThreadSafe),    // 线程安全
	}

	// 获取用户全部群Id
	contactService := service.ContactService{}
	communityList, _ := contactService.CommunityList(int(userId))
	for _, community := range communityList {
		node.GroupSets.Add(community.Id)
	}

	// 通过用户ID与Node绑定
	rwLock.Lock()
	clientMap[userId] = node
	rwLock.Unlock()

	fmt.Printf("client:%+v\n", clientMap)

	// TODO 收发信息逻辑
	go sendProc(node) // 服务端转发送信息给B客户端
	go recvProc(node) // 服务端接收A客户端信息 -> 单聊 / 群聊(B端)

	// 测试连接(字符串)
	sendMsg(userId, []byte("hello,world!----->"))
}

func MessageList(ctx *gin.Context) {
	// 获取请求参数
	userIdA, _ := strconv.Atoi(ctx.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(ctx.PostForm("userIdB"))
	start, _ := strconv.Atoi(ctx.PostForm("start"))
	end, _ := strconv.Atoi(ctx.PostForm("end"))
	msgType, _ := strconv.Atoi(ctx.PostForm("msgType"))
	if msgType == 0 {
		msgType = 1
	}
	isRev, _ := strconv.ParseBool(ctx.PostForm("isRev"))

	contactService := service.ContactService{}
	list, total := contactService.GetMessageList(int64(userIdA), int64(userIdB), int64(start), int64(end), int64(msgType), isRev)
	util.RespListOK(ctx, 0, int(total), list)
}

// 协程发送信息
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Println("数据发送出错：", err)
				return
			}
			fmt.Println("数据发送成功")
			fmt.Printf("数据:%+v\n", string(data))
		}
	}
}

// 协程接收信息
func recvProc(node *Node) {
	for {
		// WebSocket协议区分文本消息和二进制消息，使用TextMessage和BinaryMessage常量来标识。
		// 文本消息必须是UTF-8编码。conn.ReadMessage()->WriteMessage(TextMessage, data)
		// 二进制使用 conn.NextReader()->NextWriter(BinaryMessage)
		// messageType, data, err := node.Conn.ReadMessage() //文本读取，messageType=TextMessage
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println("接收数据出错,", err.Error())
			return
		}

		// A向B发送文字消息 {userid, dstid, cmd(单聊|群聊), media, content}
		// 方法一：使用TCP协议
		// dispatch(data)

		// 方法二：使用UDP协议，广播透传到局域网
		broad(data)
	}
}

func dispatch(data []byte) {
	message := model.Message{}
	if err := json.Unmarshal(data, &message); err != nil {
		fmt.Println("解析json出错,", err)
		return
	}

	fmt.Printf("接收到数据:%+v\n", string(data))
	fmt.Printf("解析到数据:%+v\n", message)
	switch message.MsgType {
	case 1: // 单聊
		//sendMsg(message.Dstobj, data)
		rwLock.RLock()
		defer rwLock.RUnlock()
		if nodeB, ok := clientMap[message.Dstobj]; ok {
			nodeB.DataQueue <- data // 把消息写入对端管道
		} else {
			fmt.Println("sendMsg failed：发送给", message.Dstobj)
		}
		// 保存数据
		messageListSave(message, data, 1)
	case 2: // 群聊
		// 每次都要遍历 clientMap,
		for _, v := range clientMap {
			if v.GroupSets.Has(message.Dstobj) {
				// 自己发的信息，排除自己，即不发送给自己
				//if message.Ownerid != userId {
				v.DataQueue <- data
				//}
			}
		}
		// 保存数据
		messageListSave(message, data, 2)
	case 3: // 心跳
		// TODO 一般什么都不做即可
	default:
		return
	}
}

// 用于存放发送要广播的数据
var updSendChan chan []byte = make(chan []byte, 1024)

func broad(data []byte) {
	updSendChan <- data
	fmt.Println("broad message received")
}

// udp协程发送数据(把单机的消息透传到局域网)
func udpSendProc() {
	log.Println("updSendProc start -->")

	// 使用udp协拨号(*内网IP段，第三位的1是本段网络段的，根据具体来填写)
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 255), // IPv4的内网网段，255是子网掩码
		Port: 3000,                       // 自定义TCP port
	})
	defer udpConn.Close()
	if err != nil {
		log.Println("[updSendProc]连接udp出错：", err.Error())
		return
	}

	// 通过udp的conn，把数据广播到局域网,updConn.Write()
	for {
		select {
		case data := <-updSendChan:
			fmt.Println("udpSendProc  data :", string(data))
			_, err := udpConn.Write(data)
			if err != nil {
				log.Println("广播消息出错：", err.Error())
				return
			}
		}
	}
}

// udp接收并处理是否传发(监听局域网内的数据)
func udpRecvProc() {
	log.Println("updRecvProc start -->")

	// 监听upd广播端口
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, // 即 0.0.0.0
		Port: 3000,
	})
	defer udpConn.Close()
	if err != nil {
		log.Println("[udpRecvProc]连接udp出错：", err.Error())
		return
	}

	// 处理端口发过来的数据
	for {
		var buffer [512]byte
		n, err := udpConn.Read(buffer[0:])
		fmt.Println("读取数据：", n)
		if err != nil {
			log.Println("读取内网数据出错：", err.Error())
			return
		}
		fmt.Println("udpRecvProc  data :", string(buffer[0:n]))

		// 接收到广播内容
		dispatch(buffer[0:n])
	}
	log.Println("[udpRecvProc] stop!!!")
}

// 发送信息
func sendMsg(userId int64, msg []byte) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if nodeB, ok := clientMap[userId]; ok {
		nodeB.DataQueue <- msg
	} else {
		fmt.Println("sendMsg failed：发送给", userId)
	}
}

// 保存消息
func messageListSave(message model.Message, data []byte, msgType int32) {
	switch msgType {
	case 1: // 单聊
		var key string
		if message.Ownerid > message.Dstobj {
			key = fmt.Sprintf("single_msg_%v_%v", message.Ownerid, message.Dstobj)
		} else {
			key = fmt.Sprintf("single_msg_%v_%v", message.Dstobj, message.Ownerid)
		}
		// 写入缓存
		ctx := context.Background()
		// 	ZREVRANGE key start stop [WITHSCORES]：返回有序集中指定区间内的成员，通过索引，分数从高到低
		//result, err := global.Rdb.ZRevRange(ctx, key, 0, -1).Result()
		//if err != nil {
		//	fmt.Println(err)
		//}
		//score := float64(cap(result)) + 1
		zcard, _ := global.Rdb.ZCard(ctx, key).Result() // zcard key：获取有序集合的成员数
		//fmt.Println("score:", score)
		fmt.Println("zcard:", zcard)
		score := float64(zcard) + 1

		res, e := global.Rdb.ZAdd(ctx, key, &redis.Z{score, string(data)}).Result()
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(res, "单聊写入redis")
	case 2: // 群聊
		var key string
		key = fmt.Sprintf("qunliao_msg_%v", message.Dstobj)
		// 写入缓存
		ctx := context.Background()
		result, err := global.Rdb.ZRevRange(ctx, key, 0, -1).Result()
		if err != nil {
			fmt.Println(err)
		}
		score := float64(cap(result)) + 1
		res, e := global.Rdb.ZAdd(ctx, key, &redis.Z{score, string(data)}).Result()
		if e != nil {
			fmt.Println(e)
		}
		fmt.Println(res, "群聊写入redis")
	}
}
