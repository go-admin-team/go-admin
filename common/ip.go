package common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetClientIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		return "127.0.0.1"
	}
	return ip
}
