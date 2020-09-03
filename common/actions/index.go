package actions

import (
	"errors"
	dto2 "go-admin/common/dto"
	"go-admin/common/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
)

// IndexAction 通用查询动作
func IndexAction(m models.ActiveRecord, d dto2.Index, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		list := f()
		object := m.Generate()
		req := d.Generate()
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		var count int64
		switch idb.(type) {
		case *gorm.DB:
			//查询列表
			db := idb.(*gorm.DB)
			err = c.Bind(req)
			tools.HasError(err, "参数验证失败", 422)
			err = req.Bind(c)
			tools.HasError(err, "参数验证失败", 422)

			//数据权限检查
			p := getPermissionFromContext(c)

			err = db.WithContext(c).Model(object).
				Scopes(
					dto2.MakeCondition(req.GetNeedSearch()),
					dto2.Paginate(req.GetPageSize(), req.GetPageIndex()),
					Permission(object.TableName(), p),
				).
				Find(list).Limit(-1).Offset(-1).
				Count(&count).Error
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tools.HasError(err, "查询失败", 500)
			}
			app.PageOK(c, list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
			c.Next()
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
	}
}
