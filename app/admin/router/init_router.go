package router

import (
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"

	"go-admin/app/admin/middleware"
	"go-admin/app/admin/middleware/handler"
	common "go-admin/common/middleware"

	"github.com/go-admin-team/go-admin-core/sdk/config"
)

// InitRouter 路由初始化，不要怀疑，这里用到了
func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		h = gin.New()
		sdk.Runtime.SetEngine(h)
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

	r.Use(common.Sentinel()).
		Use(common.RequestId(pkg.TrafficKey))
	middleware.InitMiddleware(r)
	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)
}
