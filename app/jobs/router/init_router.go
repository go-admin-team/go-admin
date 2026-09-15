package router

import (
	//"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/v2/logger"
	"github.com/go-admin-team/go-admin-core/v2/sdk"
	common "go-admin/common/middleware"
)

// InitRouter 路由初始化，不要怀疑，这里用到了
func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		log.Fatal("not found engine...")
		os.Exit(-1)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}

	// the jwt middleware: shared instance InitMiddleware built at startup,
	// not one built here per module (see common/middleware.GetAuthMiddleware).
	authMiddleware := common.GetAuthMiddleware()

	// 注册业务路由
	initRouter(r, authMiddleware)
}
