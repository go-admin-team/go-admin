package actions

import (
	"errors"
	"github.com/gin-gonic/gin"
	dto2 "go-admin/common/dto"
	"go-admin/common/models"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
)

// CreateAction 通用新增动作
func CreateAction(control dto2.Control) gin.HandlerFunc {
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
			var object models.ActiveRecord
			object, err = req.GenerateM()
			tools.HasError(err, "模型生成失败", 422)
			object.SetCreateBy(tools.GetUserIdUint(c))
			err = db.WithContext(c).Create(object).Error
			tools.HasError(err, "创建失败", 500)
			app.OK(c, object.GetId(), "创建成功")
			c.Next()
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
	}
}
