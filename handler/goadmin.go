package handler

import (
	"github.com/gin-gonic/gin"
	"go-admin/tools/app"
)

func HelloGoAdmin(c *gin.Context) {
	app.OK(c,"欢迎使用go-admin 中后台脚手架","")
}