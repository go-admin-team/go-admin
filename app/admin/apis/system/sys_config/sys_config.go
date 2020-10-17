package sys_config

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/apis"
	"go-admin/common/log"
	"go-admin/tools"
	"go-admin/tools/app"
)

type SysConfig struct {
	apis.Api
}

// GetSysConfigByKEYForService 根据Key获取SysConfig的Service
func (e *SysConfig) GetSysConfigByKEYForService(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%s", msgID, err)
		app.Error(c, 500, err, "")
		return
	}
	var v dto.SysConfigControl
	err = c.Bind(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	s := service.SysConfig{}
	s.MsgID = msgID
	s.Orm = db
	err = s.GetSysConfigByKEY(&v)
	if err != nil {
		app.Error(c, 500, err, "")
		return
	}
	app.OK(c, v, s.Msg)
}
