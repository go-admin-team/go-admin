package sysjob

import (
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/jobs"
	"go-admin/common/apis"
	"go-admin/common/dto"
	"go-admin/common/log"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
)

type SysJob struct {
	apis.Api
}

// RemoveJobForService 调用service实现
func (e *SysJob) RemoveJobForService(c *gin.Context) {
	msgID := apis.GenerateMsgIDFromContext(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%#v", msgID, err)
		app.Error(c, 500, err, "")
		return
	}
	var v dto.GeneralDelDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%#v", msgID, err)
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
	msgID := apis.GenerateMsgIDFromContext(c)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Errorf("msgID[%s] error:%#v", msgID, err)
		app.Error(c, 500, err, "")
		return
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		log.Errorf("msgID[%s] 参数验证错误, error:%#v", msgID, err)
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

func RemoveJob(c *gin.Context) {
	var data models.SysJob
	var v dto.GeneralDelDto
	err := c.BindUri(&v)
	tools.HasError(err, "", 500)
	data.JobId = v.Id
	err = data.Get(data.JobId)
	tools.HasError(err, "", 500)
	cn := jobs.Remove(data.EntryId)

	select {
	case res := <-cn:
		if res {
			_, _ = data.RemoveEntryID(data.EntryId)
			app.OK(c, nil, msg.DeletedSuccess)
		}
	case <-time.After(time.Second * 1):
		app.OK(c, nil, msg.TimeOut)
	}

}

func StartJob(c *gin.Context) {
	var data models.SysJob
	var v dto.GeneralGetDto
	err := c.BindUri(&v)
	tools.HasError(err, "", 500)
	data.JobId = v.Id
	err = data.Get(data.JobId)
	tools.HasError(err, "", 500)
	if data.JobType == 1 {
		var j = &jobs.HttpJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		data.EntryId, err = jobs.AddJob(j)
	} else {
		var j = &jobs.ExecJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		j.Args = data.Args
		data.EntryId, err = jobs.AddJob(j)
	}

	tools.HasError(err, "", 500)
	err = data.Update(data.JobId)
	tools.HasError(err, "", 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
