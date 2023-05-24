package router

import (
	"IM/gin_im/api"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		// 用户登录
		userRouter.GET("/login", api.GetLogin)         // 登录页面
		userRouter.POST("/login", api.PostLogin)       // 用户登录
		userRouter.GET("/register", api.GetRegister)   // 注册页面
		userRouter.POST("/register", api.PostRegister) // 用户注册
		userRouter.POST("/update", api.PostUpdate)     // 更新资料
		userRouter.POST("/find", api.FindById)         // 根据ID查询

		// 好友列表
		userRouter.POST("/friend", api.Friend)          // 添加好友
		userRouter.POST("/friend/list", api.FriendList) // 好友列表

		// 添加群
		userRouter.POST("/community", api.CreateCommunity)    // 创建群
		userRouter.POST("/community/add", api.AddCommunity)   // 申请入群
		userRouter.POST("/community/list", api.CommunityList) // 群列表

		// 上传文件(图片/mp3等)
		userRouter.POST("/upload", api.Upload)
	}
}
