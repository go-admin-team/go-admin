package apis

import (
	"github.com/gin-gonic/gin"
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
	c.AbortWithStatusJSON(http.StatusOK, res.ReturnError(code))
}

// OK 通常成功数据处理
func (e *Api) OK(c *gin.Context, data interface{}, msg string) {
	var res models.Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	res.RequestId = tools.GenerateMsgIDFromContext(c)
	c.AbortWithStatusJSON(http.StatusOK, res.ReturnOK())
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
	c.AbortWithStatusJSON(http.StatusOK, data)
}
