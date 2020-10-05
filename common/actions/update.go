package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/common/models"
	"go-admin/tools"
	"go-admin/tools/app"
)

// UpdateAction 通用更新动作
func UpdateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := tools.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		req := control.Generate()
		//更新操作
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
		object.SetUpdateBy(tools.GetUserIdUint(c))

		//数据权限检查
		p := GetPermissionFromContext(c)

		db = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Updates(object)
		if db.Error != nil {
			log.Errorf("MsgID[%s] Update error: %s", msgID, err)
			app.Error(c, http.StatusInternalServerError, err, "更新失败")
			return
		}
		if db.RowsAffected == 0 {
			app.Error(c, http.StatusForbidden, nil, "无权更新该数据")
			return
		}
		app.OK(c, object.GetId(), "更新成功")
		c.Next()
	}
}
