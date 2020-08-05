package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/handler"
	"go-admin/middleware"
	_ "go-admin/pkg/jwtauth"
	"go-admin/tools"
	config2 "go-admin/tools/config"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.GinEngine == nil {
		r = gin.New()
	} else {
		r = global.GinEngine
	}
	if config2.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	middleware.InitMiddleware(r)
	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)

	return r
}
