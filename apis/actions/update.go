package actions

import (
	"errors"
	"go-admin/service/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/model"
)

// UpdateAction 通用更新动作
func UpdateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := control.Generate()
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		switch idb.(type) {
		case *gorm.DB:
			//更新操作
			db := idb.(*gorm.DB)
			err = req.Bind(c)
			tools.HasError(err, "参数验证失败", 422)
			var object model.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "参数验证失败", 422)
			object.SetUpdateBy(tools.GetUserIdStr(c))
			err = db.WithContext(c).Updates(object).Error
			tools.HasError(err, "更新失败", 500)
			app.OK(c, object.GetId(), "更新成功")
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		c.Next()
	}
}
