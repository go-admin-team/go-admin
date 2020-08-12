package models

import (
	orm "go-admin/global"
	"go-admin/tools"
)

type SysJob struct {
	JobId          int    `json:"jobId" gorm:"primary_key;AUTO_INCREMENT"` // 编码
	JobName        string `json:"jobName" gorm:"size:255;"`                // 名称
	JobGroup       string `json:"jobGroup" gorm:"size:255;"`               // 任务分组
	JobType        int    `json:"jobType" gorm:"size:1;"`                  // 任务类型
	CronExpression string `json:"cronExpression" gorm:"size:255;"`         // cron表达式
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;"`           // 调用目标
	MisfirePolicy  int    `json:"misfirePolicy" gorm:"size:255;"`          // 执行策略
	Concurrent     int    `json:"concurrent" gorm:"size:1;"`               // 是否并发
	Status         int    `json:"status" gorm:"size:1;"`                   // 状态
	EntryId        int    `json:"entry_id" gorm:"size:11;"`                // job启动时返回的id
	CreateBy       string `json:"createBy" gorm:"size:128;"`               //
	UpdateBy       string `json:"updateBy" gorm:"size:128;"`               //
	DataScope      string `json:"dataScope" gorm:"-"`
	BaseModel
}

func (SysJob) TableName() string {
	return "sys_job"
}

// 创建SysJob
func (e *SysJob) Create() (SysJob, error) {
	var doc SysJob
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取SysJob
func (e *SysJob) Get() (SysJob, error) {
	var doc SysJob
	table := orm.Eloquent.Table(e.TableName())

	if e.JobId != 0 {
		table = table.Where("job_id = ?", e.JobId)
	}

	if e.JobName != "" {
		table = table.Where("job_name like ?", "%"+e.JobName+"%")
	}

	if e.JobGroup != "" {
		table = table.Where("job_group = ?", e.JobGroup)
	}

	if e.CronExpression != "" {
		table = table.Where("cron_expression = ?", e.CronExpression)
	}

	if e.InvokeTarget != "" {
		table = table.Where("invoke_target = ?", e.InvokeTarget)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取SysJob带分页
func (e *SysJob) GetPage(pageSize int, pageIndex int) ([]SysJob, int, error) {
	var doc []SysJob

	table := orm.Eloquent.Select("*").Table(e.TableName())

	if e.JobId != 0 {
		table = table.Where("job_id = ?", e.JobId)
	}

	if e.JobName != "" {
		table = table.Where("job_name like ?", "%"+e.JobName+"%")
	}

	if e.JobGroup != "" {
		table = table.Where("job_group = ?", e.JobGroup)
	}

	if e.CronExpression != "" {
		table = table.Where("cron_expression = ?", e.CronExpression)
	}

	if e.InvokeTarget != "" {
		table = table.Where("invoke_target = ?", e.InvokeTarget)
	}

	if e.Status != 0 {
		table = table.Where("status = ?", e.Status)
	}

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

func (e *SysJob) GetList() ([]SysJob, error) {
	var doc []SysJob

	table := orm.Eloquent.Select("*").Table(e.TableName())

	table = table.Where("status = ?", 2)

	if err := table.Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

// 更新SysJob
func (e *SysJob) Update(id int) (update SysJob, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("job_id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysJob) RemoveAllEntryID() (update SysJob, err error) {

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Updates(map[string]interface{}{"entry_id": 0}).Error; err != nil {
		return
	}
	return
}

func (e *SysJob) RemoveEntryID(entryID int) (update SysJob, err error) {

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Where("entry_id = ?", entryID).Updates(map[string]interface{}{"entry_id": 0}).Error; err != nil {
		return
	}
	return
}

// 删除SysJob
func (e *SysJob) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("job_id = ?", id).Delete(&SysJob{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysJob) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("job_id in (?)", id).Delete(&SysJob{}).Error; err != nil {
		return
	}
	Result = true
	return
}
