package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/logger"
	"strings"
)

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) *logger.Logger {
	requestId := pkg.GenerateMsgIDFromContext(c)
	log := sdk.Runtime.GetLogger().Fields(map[string]interface{}{
		strings.ToLower(pkg.TrafficKey): requestId,
	})
	return &logger.Logger{Logger: log}
}
