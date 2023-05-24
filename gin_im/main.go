package main

import (
	"IM/gin_im/initzlize"
)

func main() {
	// 数据库连接
	initzlize.InitDB()
	initzlize.InitRedis()

	// 初始化路由
	router := initzlize.InitRouter()

	// 静态资源
	router.Static("/asset", "asset/")
	router.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	router.LoadHTMLGlob("./view/**/*")

	// 启动http服务
	if err := router.Run(":8081"); err != nil {
		panic("启动gin服务失败")
	}
}
