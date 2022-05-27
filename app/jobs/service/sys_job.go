package service

import (
	"github.com/pkg/errors"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/robfig/cron/v3"

	"go-admin/app/jobs"
	"go-admin/app/jobs/models"
	"go-admin/common/dto"
)

type SysJob struct {
	service.Service
	Cron *cron.Cron
}

// RemoveJob 删除job
func (e *SysJob) RemoveJob(c *dto.GeneralDelDto) error {
	var err error
	var data models.SysJob
	err = e.Orm.Table(data.TableName()).First(&data, c.Id).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	cn := jobs.Remove(e.Cron, data.EntryId)

	select {
	case res := <-cn:
		if res {
			err = e.Orm.Table(data.TableName()).Where("entry_id = ?", data.EntryId).Update("entry_id", 0).Error
			if err != nil {
				e.Log.Errorf("db error: %s", err)
			}
			return err
		}
	case <-time.After(time.Second * 1):
		e.Msg = "操作超时！"
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
		e.Log.Errorf("db error: %s", err)
		return err
	}

	if data.Status == 1 {
		err = errors.New("当前Job是关闭状态不能被启动，请先启用。")
		return err
	}

	if data.JobType == 1 {
		var j = &jobs.HttpJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		data.EntryId, err = jobs.AddJob(e.Cron, j)
		if err != nil {
			e.Log.Errorf("jobs AddJob[HttpJob] error: %s", err)
		}
	} else {
		var j = &jobs.ExecJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		j.Args = data.Args
		data.EntryId, err = jobs.AddJob(e.Cron, j)
		if err != nil {
			e.Log.Errorf("jobs AddJob[ExecJob] error: %s", err)
		}
	}
	if err != nil {
		return err
	}

	err = e.Orm.Table(data.TableName()).Where(c.Id).Updates(&data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
	}
	return err
}
