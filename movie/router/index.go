package router

import (
	"IM/movie/api"
	"github.com/gin-gonic/gin"
)

func InitIndexRouter(Router *gin.RouterGroup) {
	indexRouter := Router.Group("index")
	{
		// 首页
		indexRouter.GET("/", api.Index)
		indexRouter.GET("/index", api.Index)

		// 电影座
		indexRouter.GET("/movie", api.Movie)
		indexRouter.GET("/payment", api.Payment)
	}
}
