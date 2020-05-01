package middleware

import (
	"github.com/gin-gonic/gin"
	config2 "go-admin/tools/config"
	"net/http"
)

func DemoEvn() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config2.ApplicationConfig.Mode == "demo" {
			method := c.Request.Method
			if method == "GET" || method == "OPTIONS" || c.Request.RequestURI == "/login" || c.Request.RequestURI == "/api/v1/logout" {
				c.Next()
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 500,
					"msg":  config2.ApplicationConfig.DemoMsg,
				})
				c.Abort()
				return
			}
		} else {
			c.Next()
		}
	}
}
