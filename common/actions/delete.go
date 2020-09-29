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

// DeleteAction 通用删除动作
func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := tools.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		//删除操作
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
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
		).Where(req.GetId()).Delete(object)
		if db.Error != nil {
			log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
			app.Error(c, http.StatusInternalServerError, err, "删除失败")
			return
		}
		if db.RowsAffected == 0 {
			app.Error(c, http.StatusForbidden, nil, "无权删除该数据")
			return
		}
		app.OK(c, object.GetId(), "删除成功")
		c.Next()
	}
}
