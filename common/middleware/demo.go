package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func DemoEvn() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" ||
			method == "OPTIONS" ||
			strings.Contains(c.Request.RequestURI, "/login") ||
			strings.Contains(c.Request.RequestURI, "/logout") {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "谢谢您的参与，但为了大家更好的体验，所以本次提交就算了吧！\U0001F600\U0001F600\U0001F600",
			})
			c.Abort()
			return
		}

	}
}
