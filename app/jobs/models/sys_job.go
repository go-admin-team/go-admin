package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type SysJob struct {
	JobId          int    `json:"jobId" gorm:"primaryKey;autoIncrement"` // 编码
	JobName        string `json:"jobName" gorm:"size:255;"`              // 名称
	JobGroup       string `json:"jobGroup" gorm:"size:255;"`             // 任务分组
	JobType        int    `json:"jobType" gorm:"size:1;"`                // 任务类型
	CronExpression string `json:"cronExpression" gorm:"size:255;"`       // cron表达式
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;"`         // 调用目标
	Args           string `json:"args" gorm:"size:255;"`                 // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy" gorm:"size:255;"`        // 执行策略
	Concurrent     int    `json:"concurrent" gorm:"size:1;"`             // 是否并发
	Status         int    `json:"status" gorm:"size:1;"`                 // 状态
	EntryId        int    `json:"entry_id" gorm:"size:11;"`              // job启动时返回的id
	models.ControlBy
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (SysJob) TableName() string {
	return "sys_job"
}

func (e *SysJob) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysJob) GetId() interface{} {
	return e.JobId
}

func (e *SysJob) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

func (e *SysJob) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

func (e *SysJob) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("status = ?", 2).Find(list).Error
}

// 更新SysJob
func (e *SysJob) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *SysJob) RemoveAllEntryID(tx *gorm.DB) (update SysJob, err error) {
	if err = tx.Table(e.TableName()).Where("entry_id > ?", 0).Update("entry_id", 0).Error; err != nil {
		return
	}
	return
}
