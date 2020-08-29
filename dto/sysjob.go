package dto

type SysJobSearch struct {
	Pagination     `search:"-"`
	JobId          int    `form:"jobId" search:"type:exact;column:job_id;table:sys_job"`
	JobName        string `form:"jobName" search:"type:icontains;column:job_name;table:sys_job"`
	JobGroup       string `form:"jobGroup" search:"type:exact;column:job_group;table:sys_job"`
	CronExpression string `form:"cronExpression" search:"type:exact;column:cron_expression;table:sys_job"`
	InvokeTarget   string `form:"invokeTarget" search:"type:exact;column:invoke_target;table:sys_job"`
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_job"`
}

func (m *SysJobSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysJobSearch) Validate() error {
	return nil
}

func (m *SysJobSearch) Generate() Dtor {
	o := *m
	return &o
}

func (m *SysJobSearch) GetPageIndex() int {
	return m.PageIndex
}

func (m *SysJobSearch) GetPageSize() int {
	return m.PageSize
}
