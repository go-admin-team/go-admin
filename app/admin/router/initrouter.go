package router

import (
	"os"

	"go-admin/app/admin/middleware"
	"go-admin/app/admin/middleware/handler"
	"go-admin/common/global"
	"go-admin/common/log"
	common "go-admin/common/middleware"
	_ "go-admin/pkg/jwtauth"
	"go-admin/tools"
	"go-admin/tools/config"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	var r *gin.Engine
	h := global.Cfg.GetEngine()
	if h == nil {
		h = gin.New()
		global.Cfg.SetEngine(h)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}
	if config.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}
	r.Use(middleware.WithContextDb(middleware.GetGormFromConfig(global.Cfg)))

	r.Use(common.Sentinel())
	middleware.InitMiddleware(r)
	// the jwt middleware
	var err error
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)

	//return r
}
