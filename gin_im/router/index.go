package router

import (
	"IM/gin_im/api"
	"github.com/gin-gonic/gin"
)

func InitIndexRouter(Router *gin.RouterGroup) {
	indexRouter := Router.Group("index")
	{
		// 首页
		indexRouter.GET("/", api.Index)
		indexRouter.GET("/index", api.Index)
	}
}
