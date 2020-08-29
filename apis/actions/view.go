package actions

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/model"
)

// ViewAction 通用详情动作
func ViewAction(m model.ActiveRecord) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := m.Generate()
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		switch idb.(type) {
		case *gorm.DB:
			//新增操作
			db := idb.(*gorm.DB)
			var generalGetDto tools.GeneralGetDto
			err = c.BindUri(&generalGetDto)
			tools.HasError(err, "参数验证失败", 422)
			err = db.WithContext(c).Where(generalGetDto.Id).First(object).Error
			tools.HasError(err, "查看失败", 500)
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		app.OK(c, object, "查看成功")
		c.Next()
	}
}
