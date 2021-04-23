package models

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
	CreateBy       string `json:"createBy" gorm:"size:128;"`             //
	UpdateBy       string `json:"updateBy" gorm:"size:128;"`             //
	BaseModel

	DataScope string `json:"dataScope" gorm:"-"`
}

func (SysJob) TableName() string {
	return "sys_job"
}
