package actions

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/model"
)

// DeleteAction 通用删除动作
func DeleteAction(m model.ActiveRecord) gin.HandlerFunc {
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
			var generalDelDto tools.GeneralDelDto
			err = c.Bind(&generalDelDto)
			tools.HasError(err, "参数验证失败", 422)
			err = c.BindUri(&generalDelDto)
			tools.HasError(err, "参数验证失败", 422)
			err = db.WithContext(c).Where(generalDelDto.GetIds()).Delete(object).Error
			tools.HasError(err, "更新失败", 500)
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		app.OK(c, object.GetId(), "更新成功")
		c.Next()
	}
}
