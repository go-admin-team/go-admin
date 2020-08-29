package actions

import (
	"errors"
	"go-admin/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/model"
)

// IndexAction 通用查询动作
func IndexAction(m model.ActiveRecord, d dto.Dtor) gin.HandlerFunc {
	return func(c *gin.Context) {
		object := m.Generate()
		req := d.Generate()
		var err error
		idb, exist := c.Get("db")
		if !exist {
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		list := object.GenerateList()
		var count int64
		switch idb.(type) {
		case *gorm.DB:
			//新增操作
			db := idb.(*gorm.DB)
			err = c.Bind(req)
			tools.HasError(err, "参数验证失败", 422)
			err = req.Validate()
			tools.HasError(err, "参数验证失败", 422)
			p, err := newDataPermission(db, tools.GetUserId(c))
			err = db.WithContext(c).Model(object).
				Scopes(
					tools.MakeCondition(req.GetNeedSearch()),
					tools.Paginate(req.GetPageSize(), req.GetPageIndex()),
					Permission(object.TableName(), p),
				).
				Find(list).Limit(-1).Offset(-1).
				Count(&count).Error
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tools.HasError(err, "查询失败", 500)
			}
		default:
			err = errors.New("db connect not exist")
			tools.HasError(err, "", 500)
		}
		app.PageOK(c, list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
		c.Next()
	}
}
