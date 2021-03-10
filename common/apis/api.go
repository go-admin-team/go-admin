package apis

import (
	"encoding/json"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"gorm.io/gorm"

	"go-admin/common/models"
)

type Api struct {
}

// GetLogger 获取上下文提供的日志
func (e *Api) GetLogger(c *gin.Context) *logger.Logger {
	return GetRequestLogger(c)
}

// GetOrm 获取Orm DB
func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	db, err := pkg.GetOrm(c)
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "数据库连接获取失败")
		return nil, err
	}
	return db, nil
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
	res.RequestId = pkg.GenerateMsgIDFromContext(c)
	Return := res.ReturnError(code)
	e.setResult(c, Return, res.RequestId, code)
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

// OK 通常成功数据处理
func (e *Api) OK(c *gin.Context, data interface{}, msg string) {
	var res models.Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = pkg.GenerateMsgIDFromContext(c)
	Return := res.ReturnOK()
	e.setResult(c, Return, res.RequestId, 200)
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
	msgID := pkg.GenerateMsgIDFromContext(c)
	Return := data
	e.setResult(c, Return, msgID, 200)
	c.AbortWithStatusJSON(http.StatusOK, Return)
}

func (e *Api) setResult(c *gin.Context, Return interface{}, msgID string, status int) {
	requestLogger := e.GetLogger(c)
	jsonStr, err := json.Marshal(Return)
	if err != nil {
		requestLogger.Debugf("setResult error: %#v", msgID, err.Error())
	}
	c.Set("result", string(jsonStr))
	c.Set("status", status)
}
