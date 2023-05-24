package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespOK(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
	return
}

func RespListOK(ctx *gin.Context, code int, total int, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":  code,
		"rows":  data,
		"total": total,
	})
	return
}

func RespFail(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
	})
	return
}
