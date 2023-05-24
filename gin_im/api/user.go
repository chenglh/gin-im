package api

import (
	"IM/gin_im/model"
	"IM/gin_im/service"
	"IM/gin_im/util"
	"crypto/md5"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"html/template"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 登录模板
func GetLogin(ctx *gin.Context) {
	idx, err := template.ParseFiles(
		"view/user/login.html",
		"view/chat/head.html",
	)
	if err != nil {
		panic(err)
	}
	idx.Execute(ctx.Writer, "login")
	return
}

// 登录校验
func PostLogin(ctx *gin.Context) {
	// 获取数据
	mobile := ctx.Request.FormValue("mobile")
	password := ctx.Request.FormValue("password")

	// 调用服务
	userService := service.UserService{}
	userInfo, err := userService.GetUserByMobile(mobile)

	msg := ""
	if err != nil {
		msg = err.Error()
	} else if userInfo.UserId == 0 {
		msg = "当前用户不存在"
	} else if password != userInfo.Password {
		msg = "登录密码不正确"
	}

	if msg != "" {
		util.RespFail(ctx, -1, msg)
		return
	}

	// 返回结果
	data := map[string]interface{}{
		"userid": userInfo.UserId,
		"token":  userInfo.Token,
	}
	util.RespOK(ctx, 0, "请求成功", data)
	return
}

// 注册模板
func GetRegister(ctx *gin.Context) {
	idx, err := template.ParseFiles(
		"view/user/register.html",
		"view/chat/head.html",
	)
	if err != nil {
		panic(err)
	}
	idx.Execute(ctx.Writer, "login")
	return
}

// 用户注册
func PostRegister(ctx *gin.Context) {
	// 数据迁移
	// global.DB.AutoMigrate(&model.User{})

	// 获取数据
	mobile := ctx.Request.FormValue("mobile")
	password := ctx.Request.FormValue("password")
	confirm := ctx.Request.FormValue("confirm")

	msg := ""
	if password != confirm {
		msg = "密码与确认密码不一致"
	} else if mobile == "" {
		msg = "手机号码需要填写"
	}

	if msg != "" {
		util.RespFail(ctx, -1, msg)
		return
	}

	// 调用服务
	userService := service.UserService{}
	userInfo, err := userService.GetUserByMobile(mobile)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	} else if userInfo.UserId > 0 {
		util.RespFail(ctx, -1, "当前手机号已被注册")
		return
	}

	// 手机号码中间4位替换 ****
	re, _ := regexp.Compile("(\\d{3})(\\d{4})(\\d{4})")
	nickName := re.ReplaceAllString(mobile, "$1****$3")

	userInfo.Mobile = mobile
	userInfo.Password = password
	userInfo.Token = string(fmt.Sprintf("%x", md5.Sum([]byte(password))))
	userInfo.CreatedAt = carbon.Now().Carbon2Time()
	userInfo.UpdatedAt = carbon.Now().Carbon2Time()
	userInfo.Status = 0
	userInfo.NickName = nickName

	_, err = userService.CreateUser(userInfo)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	data := map[string]interface{}{
		"userid": userInfo.UserId,
		"token":  userInfo.Token,
	}
	util.RespOK(ctx, 0, "注册成功", data)
	return
}

// 好友列表
func FriendList(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Request.FormValue("userid"))

	contactService := service.ContactService{}
	friendList, err := contactService.FriendList(int64(userId))
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	data := []map[string]interface{}{}
	for _, friend := range friendList {
		userInfo := make(map[string]interface{})
		userInfo["userid"] = friend.UserId
		userInfo["headUrl"] = friend.HeadUrl
		userInfo["nickName"] = friend.NickName
		userInfo["gender"] = friend.Gender
		data = append(data, userInfo)
	}

	util.RespListOK(ctx, 0, len(data), data)
	return
}

// 添加好友
func Friend(ctx *gin.Context) {
	// 数据迁移
	// global.DB.AutoMigrate(&model.Contact{})

	// 获取值
	userId, _ := strconv.Atoi(ctx.Request.FormValue("userid"))
	dstobj, _ := strconv.Atoi(ctx.Request.FormValue("dstobj"))

	// 初始判断
	if userId == dstobj {
		util.RespFail(ctx, -1, "不能添加自己")
		return
	}

	contactService := service.ContactService{}
	err := contactService.CreateFriend(userId, dstobj)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	util.RespOK(ctx, 0, "添加成功", nil)
	return
}

// 群聊列表
func CommunityList(ctx *gin.Context) {
	ownerId, _ := strconv.Atoi(ctx.Request.FormValue("ownerId"))

	contactService := service.ContactService{}
	communityList, err := contactService.CommunityList(ownerId)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	data := []map[string]interface{}{}
	total := len(communityList)
	for _, community := range communityList {
		comInfo := make(map[string]interface{})
		comInfo["name"] = community.Name
		comInfo["id"] = community.Id
		comInfo["icon"] = community.Icon
		comInfo["memo"] = community.Memo

		data = append(data, comInfo)
	}
	util.RespListOK(ctx, 0, total, data)
	return
}

// gin数据绑定
type CommunityForm struct {
	//支持各种类型的输入，互不影响，不同的类型可以定义不一样的名字
	OwnerId int64  `form:"ownerId" json:"ownerId" binding:"required,min=0" msg:"用户ID不能为空"`
	Icon    string `form:"icon" json:"icon" binding:"required,min=5" msg:"群头像不能为空"`
	Cate    int    `form:"cate" json:"cate" binding:"required,min=0" msg:"请正确填写类型"`
	Name    string `form:"name" json:"name" binding:"required" msg:"群名称必须填写"`
	Memo    string `form:"memo" json:"memo" binding:"required" msg:"群描述需要填写"`
}

// 创建群聊
func CreateCommunity(ctx *gin.Context) {
	// 数据迁移
	//global.DB.AutoMigrate(&model.Community{})

	// 1、gin的绑定数据
	communityForm := &CommunityForm{}

	// 绑定表单到结构体，shouldBind = 应该绑定成这样，如果不这样就抛异常
	if err := ctx.ShouldBind(communityForm); err != nil {
		util.RespFail(ctx, -1, util.GetValidMsg(err, communityForm))
		return
	}

	// 这里最好使用事务来处理以下两个步骤

	// 组装数据(群)
	nowTime := carbon.Now().Carbon2Time()
	community := model.Community{
		Name:      communityForm.Name,
		Ownerid:   communityForm.OwnerId,
		Icon:      communityForm.Icon,
		Cate:      communityForm.Cate,
		Memo:      communityForm.Memo,
		CreatedAt: &nowTime,
		UpdatedAt: &nowTime,
	}

	// 群主与群绑定
	var contact model.Contact
	contact.CreatedAt = &nowTime
	contact.UpdatedAt = &nowTime
	contact.Ownerid = communityForm.OwnerId
	contact.Dstobj = 0 // 插入群记录时再填充
	contact.Cate = 1   // 群聊类型

	var contactService = service.ContactService{}
	cid, err := contactService.CreateCommunity(community, contact)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	fmt.Printf("Client: %+v\n", clientMap)

	// 添加群时，需要维护用户的 GroupSets
	rwLock.Lock()
	if node, ok := clientMap[contact.Ownerid]; ok {
		fmt.Println("创建群ID ->", cid)
		node.GroupSets.Add(cid)
	}
	rwLock.Unlock()

	util.RespOK(ctx, 0, "添加成功", nil)
	return
}

// 加群操作
func AddCommunity(ctx *gin.Context) {
	// 获取值
	ownerId, _ := strconv.Atoi(ctx.Request.FormValue("ownerid"))
	dstobj, _ := strconv.Atoi(ctx.Request.FormValue("dstobj"))

	// 调用服务
	contactService := service.ContactService{}
	err := contactService.AddCommunity(ownerId, dstobj)
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	// 添加群时，需要维护用户的 GroupSets
	rwLock.Lock()
	if node, ok := clientMap[int64(ownerId)]; ok {
		//fmt.Println("添加群ID，更新groupSet ->", dstobj)
		node.GroupSets.Add(dstobj)
	}
	rwLock.Unlock()

	util.RespOK(ctx, 0, "加群成功", nil)
	return
}

// 上传文件方式：
// 1.上传本地【文件复制与移动】
// 2.上传云服务器
func Upload(ctx *gin.Context) {
	// 获取图片文件
	filetype := ctx.Request.FormValue("filetype") // 如果前端指定filetype
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	var url string
	/** 上传本地服务器 */
	if url, err = uploadFileToLocal(file, header, filetype); err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}

	/** 上传OSS
	if url, err = uploadFileToOss(file, header); err != nil {
		util.RespFail(ctx, -1, err.Error())
		return
	}*/

	util.RespOK(ctx, 0, "上传成功", url)
	return
}

// 上传到本地服务器
func uploadFileToLocal(file multipart.File, header *multipart.FileHeader, filetype string) (string, error) {
	suffix := ".png"            // 图片后缀
	fileName := header.Filename // 文件名称
	if tem := strings.Split(fileName, "."); len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}

	// 如果前端指定filetype，文件后缀
	if len(filetype) > 0 {
		suffix = filetype
	}

	// 初始化文件名
	newFileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)

	// 创建目录及文件
	dstFile, err := os.Create("./asset/upload/" + newFileName)
	if err != nil {
		return "", err
	}

	// 拷贝文件
	if _, err = io.Copy(dstFile, file); err != nil {
		return "", err
	}

	url := fmt.Sprintf("/asset/upload/%s", newFileName)
	return url, nil
}

// 上传到云服务器
func uploadFileToOss(file multipart.File, header *multipart.FileHeader) (string, error) {
	suffix := ".png"            // 图片后缀
	fileName := header.Filename // 文件名称
	if tem := strings.Split(fileName, "."); len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}

	// 初始化文件名
	newFileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)

	endpoint := "oss-cn-hangzhou.aliyuncs.com"
	accessKeyId := "LTAI5tNCXPJwS3MstKoKgixh"
	accessKeySecret := "YhHE8OyCMsqfjwOnxQ1oO7paYlDjVHX"
	bucketName := "ginchat"

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		//os.Exit(-1)
		return "", err
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		//os.Exit(-1)
		return "", err
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bucket.PutObject(newFileName, file)
	if err != nil {
		fmt.Println("Error:", err)
		//os.Exit(-1)
		return "", err
	}

	url := "http://" + bucketName + "." + endpoint + "/" + newFileName
	return url, nil
}

func FindById(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Request.FormValue("id"))

	userService := service.UserService{}
	userInfo, _ := userService.GetUserByUserId(int64(userId))
	if userInfo.UserId <= 0 {
		util.RespFail(ctx, -1, "参数错误")
		return
	}

	data := make(map[string]interface{})
	data["userid"] = userInfo.UserId
	data["headUrl"] = userInfo.HeadUrl
	data["nickName"] = userInfo.NickName
	data["gender"] = userInfo.Gender

	util.RespOK(ctx, 0, "查询成功", userInfo)
	return
}

// 更新资料
func PostUpdate(ctx *gin.Context) {
	// 判断用户是否存在
	userId, _ := strconv.Atoi(ctx.Request.FormValue("id"))
	icon := ctx.Request.FormValue("icon")
	nickName := ctx.Request.FormValue("name")

	if userId <= 0 || (icon == "undefined" && nickName == "undefined") {
		util.RespFail(ctx, -1, "参数错误")
		return
	}

	// 判断用户是否存在
	userService := service.UserService{}
	userInfo, _ := userService.GetUserByUserId(int64(userId))
	if userInfo.UserId <= 0 {
		util.RespFail(ctx, -1, "参数错误")
		return
	}

	// 更新用户信息
	if nickName != "undefined" {
		userInfo.NickName = nickName
	}
	if icon != "undefined" {
		userInfo.HeadUrl = icon
	}

	if err := userService.UpdateUser(userInfo); err != nil {
		util.RespFail(ctx, -1, "参数错误")
		return
	}

	util.RespOK(ctx, 0, "修改成功", nil)
	return
}
