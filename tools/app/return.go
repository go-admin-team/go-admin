package app

import (
	"encoding/json"
	"go-admin/common/log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-admin/tools"
)

// 失败数据处理
func Error(c *gin.Context, code int, err error, msg string) {
	msgID := tools.GenerateMsgIDFromContext(c)
	var res Response
	if err != nil {
		res.Msg = err.Error()
	}
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	Return := res.ReturnError(code)
	var jsonStr []byte
	jsonStr, err = json.Marshal(Return)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	c.Set("result",string(jsonStr))
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

// 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	msgID := tools.GenerateMsgIDFromContext(c)
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	Return := res.ReturnOK()
	var jsonStr []byte
	jsonStr, err := json.Marshal(Return)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	c.Set("result",string(jsonStr))
	c.AbortWithStatusJSON(http.StatusOK, Return)
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
	msgID := tools.GenerateMsgIDFromContext(c)
	Return := data
	var jsonStr []byte
	jsonStr, err := json.Marshal(Return)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	c.Set("result",string(jsonStr))
	c.AbortWithStatusJSON(http.StatusOK, Return)
}
