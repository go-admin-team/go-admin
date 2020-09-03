package actions

import (
	"errors"
	"go-admin/service/dto"
	"go-admin/tools/app"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/model"
)

// DeleteAction 通用删除动作
func DeleteAction(control dto.Control) gin.HandlerFunc {
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
			var object model.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 500)

			//数据权限检查
			object.SetUpdateBy(tools.GetUserIdStr(c))
			var p = new(dataPermission)
			if userId := tools.GetUserIdStr(c); userId != "" {
				p, err = newDataPermission(db, userId)
				tools.HasError(err, "权限范围鉴定错误", 500)
			}
			db = db.WithContext(c).Scopes(
				Permission(object.TableName(), p),
			).Delete(object)
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
