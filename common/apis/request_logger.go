package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/pkg/logger"
	"go-admin/tools"
	"go-admin/tools/app"
	"strings"
)

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) *logger.Logger {
	requestId := tools.GenerateMsgIDFromContext(c)
	log := app.Runtime.GetLogger().Fields(map[string]interface{}{
		strings.ToLower(tools.TrafficKey): requestId,
	})
	return &logger.Logger{Logger: log}
}
