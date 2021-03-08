package middleware

import (
	sentinelPlugin "github.com/alibaba/sentinel-golang/adapter/gin"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/gin-gonic/gin"

	log "github.com/go-admin-team/go-admin-core/logger"
)

// Sentinel 限流
func Sentinel() gin.HandlerFunc {
	if _, err := system.LoadRules([]*system.Rule{
		{
			MetricType:   system.InboundQPS,
			TriggerCount: 200,
			Strategy:     system.BBR,
		},
	}); err != nil {
		log.Fatalf("Unexpected error: %+v", err)
	}
	return sentinelPlugin.SentinelMiddleware()
}
