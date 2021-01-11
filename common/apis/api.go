package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-admin/common/log"
	"gorm.io/gorm"
	"net/http"

	"go-admin/common/models"
	"go-admin/tools"
)

type Api struct {
}

// GetOrm 获取Orm DB
func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	return tools.GetOrm(c)
}

// Error 通常错误数据处理
func (e *Api) Error(c *gin.Context, code int, err error, msg string) {
	var res models.Response
	if err != nil {
		res.Msg = err.Error()
	}
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	Return := res.ReturnError(code)
	setResult(c, Return, res.RequestId, code)
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

// OK 通常成功数据处理
func (e *Api) OK(c *gin.Context, data interface{}, msg string) {
	var res models.Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	Return := res.ReturnOK()
	setResult(c, Return, res.RequestId, 200)
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

// PageOK 分页数据处理
func (e *Api) PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, msg string) {
	var res models.Page
	res.List = result
	res.Count = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	e.OK(c, res, msg)
}

// Custom 兼容函数
func (e *Api) Custom(c *gin.Context, data gin.H) {
	msgID := tools.GenerateMsgIDFromContext(c)
	Return := data
	setResult(c, Return, msgID, 200)
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

func setResult(c *gin.Context, Return interface{}, msgID string, status int) {
	var jsonStr []byte
	jsonStr, err := json.Marshal(Return)
	if err != nil {
		log.Debugf("MsgID[%s] setResult error: %#v", msgID, err.Error())
	}
	c.Set("result", string(jsonStr))
	c.Set("status", status)
}
