package common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetClientIP(c *gin.Context) string {
	// 优先从 X-Forwarded-For 获取 IP
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" || strings.Contains(ip, "127.0.0.1") {
		// 如果为空或为本地地址，则尝试从 X-Real-IP 获取
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		// 如果仍为空，则使用 RemoteIP
		ip = c.RemoteIP()
	}
	if ip == "" || ip == "127.0.0.1" {
		// 如果仍为空或为本地地址，则使用 ClientIP
		ip = c.ClientIP()
	}
	if ip == "" {
		// 最后兜底为本地地址
		ip = "127.0.0.1"
	}
	return ip
}
