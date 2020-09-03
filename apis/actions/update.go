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

			//数据权限检查
			p := getPermissionFromContext(c)

			db = db.WithContext(c).Scopes(
				Permission(object.TableName(), p),
			).Updates(object)
			tools.HasError(db.Error, "更新失败", 500)
			if db.RowsAffected == 0 {
				err = errors.New("无权更新该数据")
				tools.HasError(err, "", 403)
			}
			app.OK(c, object.GetId(), "更新成功")
			c.Next()
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
	}
}
