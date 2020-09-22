package apis

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
)

type Api struct {
}

func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	return GetOrm(c)
}

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	msgID := tools.GenerateMsgIDFromContext(c)
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New(fmt.Sprintf("msgID[%s], db connect not exist", msgID))
	}
	switch idb.(type) {
	case *gorm.DB:
		//新增操作
		return idb.(*gorm.DB), nil
	default:
		return nil, errors.New(fmt.Sprintf("msgID[%s], db connect not exist", msgID))
	}
}
