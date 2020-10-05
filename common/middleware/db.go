package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithContextDb(dbMap map[string]*gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db, ok := dbMap["*"]; ok {
			c.Set("db", db)
		} else {
			c.Set("db", dbMap[c.Request.Host])
		}
		c.Next()
	}
}
