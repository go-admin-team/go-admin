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

// CreateAction 通用新增动作
func CreateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := tools.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := tools.GenerateMsgIDFromContext(c)
		//新增操作
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
		object.SetCreateBy(tools.GetUserIdUint(c))
		err = db.WithContext(c).Create(object).Error
		if err != nil {
			log.Errorf("MsgID[%s] Create error: %s", msgID, err)
			app.Error(c, http.StatusInternalServerError, err, "创建失败")
			return
		}
		app.OK(c, object.GetId(), "创建成功")
		c.Next()
	}
}
