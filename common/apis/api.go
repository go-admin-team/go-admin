package apis

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"
)

type Api struct {
	Context *gin.Context
}

// GetLogger 获取上下文提供的日志
func (e Api) GetLogger() *logger.Logger {
	return api.GetRequestLogger(e.Context)
}

// GetOrm 获取Orm DB
func (e Api) GetOrm() (*gorm.DB, error) {
	db, err := pkg.GetOrm(e.Context)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return nil, err
	}
	return db, nil
}

// Error 通常错误数据处理
func (e Api) Error(code int, err error, msg string) {
	response.Error(e.Context, code, err, msg)
}

// OK 通常成功数据处理
func (e Api) OK(data interface{}, msg string) {
	response.OK(e.Context, data, msg)
}

// PageOK 分页数据处理
func (e Api) PageOK(result interface{}, count int, pageIndex int, pageSize int, msg string) {
	response.PageOK(e.Context, result, count, pageIndex, pageSize, msg)
}

// Custom 兼容函数
func (e Api) Custom(data gin.H) {
	response.Custum(e.Context, data)
}

// Bind 绑定参数，⚠️注意：uri绑定，请传nil
func (e Api) Bind(c interface{}, bindings ...binding.Binding) error {
	var err error
	for i := range bindings {
		switch bindings[i] {
		case binding.JSON:
			err = e.Context.ShouldBindWith(c, binding.JSON)
		case binding.XML:
			err = e.Context.ShouldBindWith(c, binding.XML)
		case binding.Form:
			err = e.Context.ShouldBindWith(c, binding.Form)
		case binding.Query:
			err = e.Context.ShouldBindWith(c, binding.Query)
		case binding.FormPost:
			err = e.Context.ShouldBindWith(c, binding.FormPost)
		case binding.FormMultipart:
			err = e.Context.ShouldBindWith(c, binding.FormMultipart)
		case binding.ProtoBuf:
			err = e.Context.ShouldBindWith(c, binding.ProtoBuf)
		case binding.MsgPack:
			err = e.Context.ShouldBindWith(c, binding.MsgPack)
		case binding.YAML:
			err = e.Context.ShouldBindWith(c, binding.YAML)
		case binding.Header:
			err = e.Context.ShouldBindWith(c, binding.Header)
		default:
			err = e.Context.ShouldBindUri(c)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
