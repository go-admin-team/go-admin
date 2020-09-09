package actions

import (
	"errors"
	dto2 "go-admin/common/dto"
	"go-admin/common/models"
	"go-admin/tools/app"
	"gopkg.in/ffmt.v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
)

// ViewAction 通用详情动作
func ViewAction(control dto2.Control) gin.HandlerFunc {
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
			var object models.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 500)

			//数据权限检查
			p := getPermissionFromContext(c)

			ffmt.P(object)
			err = db.WithContext(c).Scopes(
				Permission(object.TableName(), p),
			).Where(req.GetId()).First(object).Error
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
