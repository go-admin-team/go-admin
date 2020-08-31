package actions

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/model"
)

// UpdateAction 通用更新动作
func UpdateAction(m model.ActiveRecord) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.SetUpdateBy(tools.GetUserIdStr(c))
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
			err = c.Bind(object)
			tools.HasError(err, "参数验证失败", 422)
			err = db.WithContext(c).Updates(object).Error
			tools.HasError(err, "更新失败", 500)
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		app.OK(c, object.GetId(), "更新成功")
		c.Next()
	}
}
