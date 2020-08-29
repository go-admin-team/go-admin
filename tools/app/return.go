package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 失败数据处理
func Error(c *gin.Context, code int, err error, msg string) {
	var res Response
	res.Msg = err.Error()
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnError(code))
}

// 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}

// 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, msg string) {
	var res Page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, msg)
}

// 兼容函数
func Custum(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, data)
}
