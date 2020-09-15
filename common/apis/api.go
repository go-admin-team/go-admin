package apis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Api struct {
}

// GetOrm 获取orm连接
func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New("db connect not exist")
	}
	switch idb.(type) {
	case *gorm.DB:
		//新增操作
		return idb.(*gorm.DB), nil
	default:
		return nil, errors.New("db connect not exist")
	}
}
