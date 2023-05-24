package initzlize

import (
	"IM/movie/router"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	defaultRouter := gin.Default()

	// 初始化用户路由
	apiGroup := defaultRouter.Group("/v1")
	router.InitIndexRouter(apiGroup) // 首页展示

	return defaultRouter
}
