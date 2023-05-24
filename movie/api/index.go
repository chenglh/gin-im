package api

import (
	"IM/movie/global"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// 客户端的映射关系
var clientMap map[int64]*Node = make(map[int64]*Node)

// 读写锁
var rwLock sync.RWMutex

// http://127.0.0.1:8081/v1/index/?userid=12000&movie=tjxs
// 电影详情模板加载
func Index(ctx *gin.Context) {
	templates := []string{
		"view/movie/index.html",
	}
	idx, err := template.ParseFiles(templates...)
	if err != nil {
		panic(err)
	}
	//idx.ExecuteTemplate() 这个方法有三个参数
	idx.Execute(ctx.Writer, "index")
}

// js触发ws请求，升级websocket，并把当前电影票座位数据返回给当前用户
func Movie(ctx *gin.Context) {
	// 接收用户ID
	uid := ctx.Query("userid")
	movieName := ctx.Query("movie")
	userid, _ := strconv.ParseInt(uid, 10, 64)

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		// 如果前端js报错，websocket: request origin not allowed by Upgrader.CheckOrigin 是这里的跨域问题
		return true
	}}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("升级websocket失败", err.Error())
		return
	}

	nowTime := carbon.Now().Timestamp()
	node := &Node{
		Conn:          conn,                       // 客户端连接
		Addr:          conn.RemoteAddr().String(), // 客户端地址
		HeartbeatTime: nowTime,                    // 心跳时间
		LoginTime:     nowTime,                    // 登录时间
		DataQueue:     make(chan []byte, 50),      // 消息队列
	}

	// 通过用户ID与Node绑定
	rwLock.Lock()
	clientMap[userid] = node
	rwLock.Unlock()

	// 开启消息发送的协程
	go sendProc(node)

	// 获取数据返回
	if movieName != "" {
		result := make(map[string]interface{})
		result["data"], _ = global.Rdb.LRange(ctx, movieName, 0, -1).Result()
		result["movie"] = movieName

		data, _ := json.Marshal(result)
		fmt.Println(result["data"], "页面详情数据")

		// 把电影售票结果发送给当前用户
		sendMsg(userid, data)
	}
}

// 支付成功把当前电影票座位数据推送给所有客户端，客户端根据电影名来判断是否是判断页面的电影
func Payment(ctx *gin.Context) {
	//uid := ctx.Query("userid")
	movieName := ctx.Query("movie")
	ids := ctx.Query("ids")
	//userid, _ := strconv.ParseInt(uid, 10, 64)

	if movieName == "" || ids == "" {
		ctx.JSON(200, gin.H{
			"message": "参数缺失",
		})
		return
	}

	// 锁
	rwLock.Lock()
	defer rwLock.Unlock()

	var movieSeat = []string{}
	ids = strings.Replace(ids, "cart-item-", "", -1)
	splitArr := strings.Split(ids, "|")
	for _, split := range splitArr {
		movieSeat = append(movieSeat, split)
	}

	// 从数据表中取出数据，判断是否有交集
	result, _ := global.Rdb.LRange(ctx, movieName, 0, -1).Result()

	if len(result) > 0 {
		if iSeat := intersect(result, movieSeat); len(iSeat) > 0 {
			ctx.JSON(200, gin.H{
				"message": "有票已经被抢购啦，" + strings.Join(iSeat, "|"),
			})
			return
		}
	}

	// 持久化数据后
	for _, seat := range movieSeat {
		global.Rdb.LPush(ctx, movieName, seat)
	}

	// fmt.Println(movieSeat)
	// 重新获取数据，推送其他客户端
	responseData := make(map[string]interface{})
	responseData["data"], _ = global.Rdb.LRange(ctx, movieName, 0, -1).Result()
	responseData["movie"] = movieName
	data, _ := json.Marshal(responseData)

	for _, v := range clientMap {
		v.DataQueue <- data
	}

	ctx.JSON(200, gin.H{
		"message": "支付成功",
	})
	return
}

// 发送信息，把消息放入数据列表
func sendMsg(userid int64, msg []byte) {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if nodeB, ok := clientMap[userid]; ok {
		nodeB.DataQueue <- msg // 把数据写入队列中
	} else {
		fmt.Println("sendMsg failed：发送给", userid)
	}
}

// 执行真正的数据发送
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

func intersect(slice1 []string, slice2 []string) []string { // 取两个切片的交集
	m := make(map[string]int)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}

	return n
}
