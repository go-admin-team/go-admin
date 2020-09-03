package actions

import (
	"errors"
	dto2 "go-admin/common/dto"
	"go-admin/common/models"
	"go-admin/tools/app"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
)

// DeleteAction 通用删除动作
func DeleteAction(control dto2.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		switch idb.(type) {
		case *gorm.DB:
			//删除操作
			db := idb.(*gorm.DB)
			req := control.Generate()
			err = req.Bind(c)
			tools.HasError(err, "参数验证失败", 422)
			var object models.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 500)

			object.SetUpdateBy(tools.GetUserIdUint(c))

			//数据权限检查
			p := getPermissionFromContext(c)

			db = db.WithContext(c).Scopes(
				Permission(object.TableName(), p),
			).Where(req.GetId()).Delete(object)
			tools.HasError(db.Error, "删除失败", 500)
			if db.RowsAffected == 0 {
				err = errors.New("无权删除该数据")
				tools.HasError(err, "", 403)
			}
			app.OK(c, object.GetId(), "删除成功")
			c.Next()
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
	}
}
