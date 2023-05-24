package router

import (
	"IM/gin_im/api"
	"github.com/gin-gonic/gin"
)

func InitChatRouter(Router *gin.RouterGroup) {
	chatRouter := Router.Group("chat")
	{
		// 登录后跳转聊天页面js的websocket
		// ws://127.0.0.1:8082/v1/chat?userId=1&token=xxx
		chatRouter.GET("", api.Chat)

		// 聊天记录
		chatRouter.POST("/message", api.MessageList)
	}
}
