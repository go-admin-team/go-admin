package router

import (
	"github.com/go-admin-team/go-admin-core/sdk"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/tools/transfer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Monitor() {
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
	}
	//开发环境启动监控指标
	r.GET("/metrics", transfer.Handler(promhttp.Handler()))
	//健康检查
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
