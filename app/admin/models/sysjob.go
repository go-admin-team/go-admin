package models

import (
	"go-admin/common/dto"
	orm "go-admin/common/global"
	"go-admin/common/models"
	"go-admin/tools"
	"strconv"
)

type SysJob struct {
	JobId          uint   `json:"jobId" gorm:"primary_key;AUTO_INCREMENT"` // 编码
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
	BaseModel

	DataScope      string `json:"dataScope" gorm:"-"`
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

func (e *SysJob) SetCreateBy(createBy uint) {
	e.CreateBy = strconv.Itoa(int(createBy))
}

func (e *SysJob) SetUpdateBy(updateBy uint) {
	e.UpdateBy = strconv.Itoa(int(updateBy))
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
func (e *SysJob) GetPage(pageSize int, pageIndex int, v interface{}, list interface{}) (int, error) {
	table := orm.Eloquent.Table(e.TableName()).Scopes(dto.MakeCondition(v))

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	//dataPermission := new(DataPermission)
	userid, _ := tools.StringToInt(e.DataScope)

	var count int64

	if err := table.Scopes(DataScopes(e.TableName(), userid), dto.Paginate(pageSize, pageIndex)).Find(list).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (e *SysJob) GetList(list interface{}) (err error) {
	return orm.Eloquent.Table(e.TableName()).Where("status = ?", 2).Find(list).Error
}

// 更新SysJob
func (e *SysJob) Update(id interface{}) (err error) {
	return orm.Eloquent.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *SysJob) RemoveAllEntryID() (update SysJob, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("entry_id > ?", 0).Update("entry_id", 0).Error; err != nil {
		return
	}
	return
}

func (e *SysJob) RemoveEntryID(entryID int) (update SysJob, err error) {
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
