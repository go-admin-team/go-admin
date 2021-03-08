package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/tools/app"
)

func WithContextDb(c *gin.Context) {
	c.Set("db", app.Runtime.GetDbByKey(c.Request.Host))
	c.Next()
}
