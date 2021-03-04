package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/google/uuid"
)

// RequestIdLogger 自动增加requestId
func RequestIdLogger(trafficKey, loggerKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader(trafficKey)
		if requestId == "" {
			requestId = c.GetHeader(strings.ToLower(trafficKey))
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Header(trafficKey, requestId)
		c.Set(loggerKey, logger.DefaultLogger.Fields(map[string]interface{}{
			strings.ToLower(trafficKey): requestId,
		}))
		c.Next()
	}
}
