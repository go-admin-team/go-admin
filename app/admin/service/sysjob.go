package service

import (
	"go-admin/app/admin/models"
	"go-admin/app/jobs"
	"go-admin/common/dto"
	"go-admin/common/service"
	"go-admin/tools/app/msg"
	"time"
)

type SysJob struct {
	service.Service
}

// RemoveJob 删除job
func (e *SysJob) RemoveJob(c *dto.GeneralDelDto) error {
	var err error
	var data models.SysJob
	data.JobId = c.Id
	err = e.Orm.Table(data.TableName()).First(&data).Error
	if err != nil {
		return err
	}
	cn := jobs.Remove(data.EntryId)

	select {
	case res := <-cn:
		if res {
			err = e.Orm.Table(data.TableName()).Where("entry_id = ?", data.EntryId).Update("entry_id", 0).Error
			return err
		}
	case <-time.After(time.Second * 1):
		e.Msg = msg.TimeOut
		return nil
	}
	return nil
}

// StartJob 启动任务
func (e *SysJob) StartJob(c *dto.GeneralGetDto) error {
	var data models.SysJob
	var err error
	err = e.Orm.Table(data.TableName()).First(&data, c.Id).Error
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}

	err = e.Orm.Table(data.TableName()).Where(c.Id).Updates(&data).Error
	return err
}
