package initzlize

import (
	"IM/gin_im/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	defaultRouter := gin.Default()

	// 初始化用户路由
	apiGroup := defaultRouter.Group("/v1")
	router.InitUserRouter(apiGroup)  // 用户相关
	router.InitIndexRouter(apiGroup) // 首页展示
	router.InitChatRouter(apiGroup)  // 聊天相关

	return defaultRouter
}
