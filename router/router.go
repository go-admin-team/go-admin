package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
	_ "go-admin/pkg/jwtauth"
	"go-admin/tools"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	// 日志处理
	r.Use(middleware.LoggerToFile())
	// 自定义错误处理
	r.Use(middleware.CustomError)
	// NoCache is a middleware function that appends headers
	r.Use(middleware.NoCache)
	// 跨域处理
	r.Use(middleware.Options)
	// Secure is a middleware function that appends security
	r.Use(middleware.Secure)
	// Set X-Request-Id header
	r.Use(middleware.RequestId())
	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由是有演示代码
	InitExamplesRouter(r, authMiddleware)

	return r
}
