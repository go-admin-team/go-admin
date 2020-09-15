package sysjob

import (
	"go-admin/app/admin/service"
	"go-admin/common/apis"
	"go-admin/common/dto"
	"gorm.io/gorm"
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/app/admin/models"
	"go-admin/app/jobs"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
)

type SysJob struct {
	apis.Api
}

// RemoveJobForService 调用service实现
func (e *SysJob) RemoveJobForService(c *gin.Context) {
	db, err := e.GetOrm(c)
	if err != nil {
		app.Error(c, 500, err, "")
	}
	var v dto.GeneralDelDto
	err = c.BindUri(&v)
	if err != nil {
		app.Error(c, 422, err, "参数验证失败")
	}
	s := service.SysJob{}
	s.Orm = db
	err = s.RemoveJob(&v)
	if err != nil {
		app.Error(c, 500, err, "")
	}
	app.OK(c, nil, s.Msg)
}

// StartJobForService 启动job service实现
func (e *SysJob) StartJobForService(c *gin.Context) {
	db, err := e.GetOrm(c)
	if err != nil {
		app.Error(c, 500, err, "")
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		app.Error(c, 422, err, "参数验证失败")
	}
	s := service.SysJob{}
	s.Orm = db
	err = s.StartJob(&v)
	if err != nil {
		app.Error(c, 500, err, "")
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

// StartJobForService 启动job service实现
func StartJobForService(c *gin.Context) {
	var err error
	idb, exist := c.Get("db")
	if !exist {
		app.Error(c, 500, nil, "db connect not exist")
	}
	switch idb.(type) {
	case *gorm.DB:
		//新增操作
		db := idb.(*gorm.DB)
		var v dto.GeneralGetDto
		err = c.BindUri(&v)
		if err != nil {
			app.Error(c, 422, err, "参数验证失败")
		}
		s := service.SysJob{}
		s.Orm = db
		err = s.StartJob(&v)
		if err != nil {
			app.Error(c, 500, err, "")
		}
		app.OK(c, nil, s.Msg)
	default:
		app.Error(c, 500, nil, "db connect not exist")
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
