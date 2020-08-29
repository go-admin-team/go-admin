package sysjob

import (
	"time"

	"github.com/gin-gonic/gin"

	"go-admin/dto"
	"go-admin/jobs"
	"go-admin/models"
	"go-admin/tools"
	"go-admin/tools/app"
	"go-admin/tools/app/msg"
)

func GetSysJobList(c *gin.Context) {
	var data models.SysJob
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize,err = tools.StringToInt(size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex,err = tools.StringToInt(index)
	}

	var v dto.SysJobSearch
	err = c.Bind(&v)
	tools.HasError(err, "数据解析失败", 422)

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex, v)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

func GetSysJob(c *gin.Context) {
	var data models.SysJob
	var v tools.GeneralGetDto
	err := c.BindUri(&v)
	tools.HasError(err, "", 500)
	data.JobId, _ = tools.StringToInt(v.Id)
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

func InsertSysJob(c *gin.Context) {
	var data models.SysJob
	err := c.ShouldBindJSON(&data)
	data.CreateBy = tools.GetUserIdStr(c)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

func UpdateSysJob(c *gin.Context) {
	var data models.SysJob
	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "数据解析失败", -1)
	data.UpdateBy = tools.GetUserIdStr(c)
	result, err := data.Update(data.JobId)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

func DeleteSysJob(c *gin.Context) {
	var data models.SysJob

	var v tools.GeneralDelDto
	err := c.BindUri(&v)
	tools.HasError(err, "", 500)
	data.UpdateBy = tools.GetUserIdStr(c)
	IDS := tools.IdsStrToIdsIntGroupStr(v.Id)
	_, err = data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}

func RemoveJob(c *gin.Context) {
	var data models.SysJob
	var v tools.GeneralDelDto
	err := c.BindUri(&v)
	tools.HasError(err, "", 500)
	data.JobId, _ = tools.StringToInt(v.Id)
	result, err := data.Get()
	tools.HasError(err, "", 500)
	cn := jobs.Remove(result.EntryId)

	select {
	case res := <-cn:
		if res {
			_, _ = data.RemoveEntryID(result.EntryId)
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
	data.JobId, _ = tools.StringToInt(v.Id)
	result, err := data.Get()
	tools.HasError(err, "", 500)
	if result.JobType == 1 {
		var j = &jobs.HttpJob{}
		j.InvokeTarget = result.InvokeTarget
		j.CronExpression = result.CronExpression
		j.JobId = result.JobId
		j.Name = result.JobName
		data.EntryId, err = jobs.AddJob(j)
	} else {
		var j = &jobs.ExecJob{}
		j.InvokeTarget = result.InvokeTarget
		j.CronExpression = result.CronExpression
		j.JobId = result.JobId
		j.Name = result.JobName
		j.Args = result.Args
		data.EntryId, err = jobs.AddJob(j)
	}

	tools.HasError(err, "", 500)
	_, err = data.Update(data.JobId)
	tools.HasError(err, "", 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
