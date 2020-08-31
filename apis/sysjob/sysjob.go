package sysjob

import (
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/jobs"
	"go-admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
)

func RemoveJob(c *gin.Context) {
	var data models.SysJob
	var v tools.GeneralDelDto
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
	var v tools.GeneralGetDto
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
