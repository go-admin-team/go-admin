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

// CreateAction 通用新增动作
func CreateAction(control dto.Control) gin.HandlerFunc {
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
			//新增操作
			db := idb.(*gorm.DB)
			err = req.Bind(c)
			tools.HasError(err, "参数验证失败", 422)
			var object model.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 422)
			object.SetCreateBy(tools.GetUserIdStr(c))
			err = db.WithContext(c).Create(object).Error
			tools.HasError(err, "创建失败", 500)
			app.OK(c, object.GetId(), "创建成功")
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		c.Next()
	}
}
