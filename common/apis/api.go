package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Api struct {
}

func (e *Api) GetOrm(c *gin.Context) (*gorm.DB, error) {
	return GetOrm(c)
}

// GetOrm 获取orm连接
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	msgID := GenerateMsgIDFromContext(c)
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

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	var msgID string
	data, ok := c.Get("msgID")
	if !ok {
		msgID = uuid.New().String()
		c.Set("msgID", msgID)
		return msgID
	}
	msgID = cast.ToString(data)
	return msgID
}
