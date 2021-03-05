package middleware

import (
	"github.com/gin-gonic/gin"

	"go-admin/common/global"
)

func WithContextDb(c *gin.Context) {
	c.Set("db", global.Runtime.GetDbByKey(c.Request.Host))
	c.Next()
}
