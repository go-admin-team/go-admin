package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"
)

type Api struct {
}

// GetLogger 获取上下文提供的日志
func (e *Api) GetLogger(c *gin.Context) *logger.Logger {
	return api.GetRequestLogger(c)
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
	response.Error(c, code, err, msg)
}

// OK 通常成功数据处理
func (e *Api) OK(c *gin.Context, data interface{}, msg string) {
	response.OK(c, data, msg)
}

// PageOK 分页数据处理
func (e *Api) PageOK(c *gin.Context, result interface{}, count int, pageIndex int, pageSize int, msg string) {
	response.PageOK(c, result, count, pageIndex, pageSize, msg)
}

// Custom 兼容函数
func (e *Api) Custom(c *gin.Context, data gin.H) {
	response.Custum(c, data)
}
