package apis

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
)

type Api struct {
}

func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	return tools.GetOrm(c)
}
