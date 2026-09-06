package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/v2/tools/transfer"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-admin/common/health"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, RegisterMonitorRouter)
}

// readyTimeout bounds the whole probe. It has to stay under whatever period
// the orchestrator polls on, or a slow dependency turns a readiness check into
// a queue of readiness checks.
const readyTimeout = 2 * time.Second

// HealthPath and ReadyPath are the two probe routes, relative to APIPrefix.
//
// Exported for the same reason as the prefix: the rate limiter has to be told
// to skip them, and it is installed in a package that cannot import this one.
const (
	HealthPath = "/health"
	ReadyPath  = "/ready"
)

// RegisterMonitorRouter mounts the metrics endpoint and the two probes on v1.
//
// Exported so that a test can put the real probes on a server of its own. The
// alternative - a test that re-implements the handler it means to check - is
// how a probe comes to be asserted against a copy of itself.
//
// 无需认证的路由代码
func RegisterMonitorRouter(v1 *gin.RouterGroup) {
	v1.GET("/metrics", transfer.Handler(promhttp.Handler()))

	// 健康检查（存活）
	//
	// Stays a bare 200 on purpose. This is the answer to "should I restart
	// you", and a process whose database is unreachable does not want
	// restarting - that turns one outage into a crash loop and throws away the
	// connection pool, the cache and every in-flight request along the way.
	v1.GET(HealthPath, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 就绪检查
	//
	// The answer to "should I send you requests". It fails while a dependency
	// is unreachable, and from the moment shutdown begins. Nothing waits on
	// that second answer today, so it is readable rather than actionable -
	// the package comment in common/health says what it would take.
	v1.GET(ReadyPath, func(c *gin.Context) {
		if health.Draining() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "draining",
				"checks": []health.Check{},
			})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), readyTimeout)
		defer cancel()

		checks := health.Ready(ctx)
		status := http.StatusOK
		if !health.Healthy(checks) {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, gin.H{"status": http.StatusText(status), "checks": checks})
	})

}
