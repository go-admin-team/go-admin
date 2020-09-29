package sysjob

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/admin/service"
	"go-admin/common/apis"
	"go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/tools"
	"go-admin/tools/app"
)

type SysJob struct {
	apis.Api
}

// RemoveJobForService 调用service实现
func (e *SysJob) RemoveJobForService(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%s", msgID, err)
		app.Error(c, 500, err, "")
		return
	}
	var v dto.GeneralDelDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.MsgID = msgID
	s.Orm = db
	err = s.RemoveJob(&v)
	if err != nil {
		app.Error(c, 500, err, "")
		return
	}
	app.OK(c, nil, s.Msg)
}

// StartJobForService 启动job service实现
func (e *SysJob) StartJobForService(c *gin.Context) {
	msgID := tools.GenerateMsgIDFromContext(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%s", msgID, err)
		app.Error(c, 500, err, "")
		return
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%s", msgID, err)
		app.Error(c, 422, err, "参数验证失败")
		return
	}
	s := service.SysJob{}
	s.Orm = db
	s.MsgID = msgID
	err = s.StartJob(&v)
	if err != nil {
		app.Error(c, 500, err, "")
		return
	}
	app.OK(c, nil, s.Msg)
}
