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

// ViewAction 通用详情动作
func ViewAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		switch idb.(type) {
		case *gorm.DB:
			//查看详情
			db := idb.(*gorm.DB)
			req := control.Generate()
			err = req.Bind(c)
			tools.HasError(err, "参数验证失败", 422)
			var object model.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 500)
			var p = new(dataPermission)
			if userId := tools.GetUserIdStr(c); userId != "" {
				p, err = newDataPermission(db, userId)
				tools.HasError(err, "权限范围鉴定错误", 500)
			}
			err = db.WithContext(c).Scopes(
				Permission(object.TableName(), p),
			).First(object).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tools.HasError(err, "查看失败", 404)
			}
			tools.HasError(err, "查看失败", 500)
			app.OK(c, object, "查看成功")
			c.Next()
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
	}
}
