package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/common/middleware"
)

func InitMiddleware(r *gin.Engine) {
	// 日志处理
	r.Use(LoggerToFile())
	// 自定义错误处理
	r.Use(CustomError)
	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// 跨域处理
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	// 链路追踪
	r.Use(middleware.Trace())
}
