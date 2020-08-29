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
	Args           string `json:"args" gorm:"size:255;"`                   // 目标参数
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
func (e *SysJob) Create() (err error) {
	return orm.Eloquent.Table(e.TableName()).Create(e).Error
}

// 获取SysJob
func (e *SysJob) Get(id interface{}) (err error) {
	return orm.Eloquent.Table(e.TableName()).First(e, id).Error
}

// 获取SysJob带分页
func (e *SysJob) GetPage(pageSize int, pageIndex int, v interface{}) ([]SysJob, int, error) {
	var doc []SysJob

	table := orm.Eloquent.Table(e.TableName()).Scopes(tools.MakeCondition(v))

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int64

	if err := table.Scopes(tools.Paginate(pageSize, pageIndex)).Find(&doc).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return doc, int(count), nil
}

func (e *SysJob) GetList() ([]SysJob, error) {
	var doc []SysJob

	table := orm.Eloquent.Table(e.TableName())

	table = table.Where("status = ?", 2)

	if err := table.Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

// 更新SysJob
func (e *SysJob) Update(id interface{}) (rowsAffected int64, err error) {
	result := orm.Eloquent.Table(e.TableName()).Where(id).Updates(&e)
	return result.RowsAffected, result.Error
}

func (e *SysJob) RemoveAllEntryID() (update SysJob, err error) {

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Where("entry_id > ?", 0).Update("entry_id", 0).Error; err != nil {
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
	if err = orm.Eloquent.Table(e.TableName()).Where(id).Delete(&SysJob{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *SysJob) BatchDelete(id []int) error {
	return orm.Eloquent.Table(e.TableName()).Where(id).Delete(&SysJob{}).Error
}
