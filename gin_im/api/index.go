package api

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

// 应用首页
func Index(ctx *gin.Context) {
	templates := []string{
		"view/chat/index.html",
		"view/chat/head.html",
		"view/chat/foot.html",
		"view/chat/tabmenu.html",
		"view/chat/concat.html",
		"view/chat/group.html",
		"view/chat/profile.html",
		"view/chat/createcom.html",
		"view/chat/userinfo.html",
		"view/chat/main.html",
	}
	idx, err := template.ParseFiles(templates...)
	if err != nil {
		panic(err)
	}
	//idx.ExecuteTemplate() 这个方法有三个参数
	idx.Execute(ctx.Writer, "index")
}
