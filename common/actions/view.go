package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/common/apis"
	"go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/models"
	"go-admin/tools"
	"go-admin/tools/app"
)

// ViewAction 通用详情动作
func ViewAction(control dto.Control, params ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := apis.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		//查看详情
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			app.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			app.Error(c, http.StatusInternalServerError, err, "模型生成失败")
			return
		}

		//数据权限检查
		p := getPermissionFromContext(c)

		err = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).First(object).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.Error(c, http.StatusNotFound, nil, "查看对象不存在或无权查看")
			return
		}
		if err != nil {
			log.Errorf("MsgID[%s] Create error: %#v", msgID, err)
			app.Error(c, http.StatusInternalServerError, err, "查看失败")
			return
		}
		app.OK(c, object, "查看成功")
		c.Next()
	}
}
