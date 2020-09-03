package dto

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	models2 "go-admin/common/models"
	"net/http"
)

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

func (m *SysJobSearch) Bind(ctx *gin.Context) error {
	return ctx.Bind(m)
}

func (m *SysJobSearch) Generate() dto.Index {
	o := *m
	return &o
}

func (m *SysJobSearch) GetPageIndex() int {
	return m.PageIndex
}

func (m *SysJobSearch) GetPageSize() int {
	return m.PageSize
}

type SysJobControl struct {
	JobId          int    `json:"jobId"`
	JobName        string `json:"jobName" validate:"required"` // 名称
	JobGroup       string `json:"jobGroup"`                    // 任务分组
	JobType        int    `json:"jobType"`                     // 任务类型
	CronExpression string `json:"cronExpression"`              // cron表达式
	InvokeTarget   string `json:"invokeTarget"`                // 调用目标
	Args           string `json:"args"`                        // 目标参数
	MisfirePolicy  int    `json:"misfirePolicy"`               // 执行策略
	Concurrent     int    `json:"concurrent"`                  // 是否并发
	Status         int    `json:"status"`                      // 状态
	EntryId        int    `json:"entry_id"`                    // job启动时返回的id
}

func (s *SysJobControl) Bind(ctx *gin.Context) error {
	return ctx.Bind(s)
}

func (s *SysJobControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysJobControl) GenerateM() (models2.ActiveRecord, error) {
	return &models.SysJob{
		JobId:          s.JobId,
		JobName:        s.JobName,
		JobGroup:       s.JobGroup,
		JobType:        s.JobType,
		CronExpression: s.CronExpression,
		InvokeTarget:   s.InvokeTarget,
		Args:           s.Args,
		MisfirePolicy:  s.MisfirePolicy,
		Concurrent:     s.Concurrent,
		Status:         s.Status,
		EntryId:        s.EntryId,
	}, nil
}

func (s *SysJobControl) GetId() interface{} {
	return s.JobId
}

type SysJobById struct {
	Id  int   `uri:"id" validate:"required"`
	Ids []int `json:"ids"`
}

func (s *SysJobById) Bind(ctx *gin.Context) error {
	if ctx.Request.Method == http.MethodDelete {
		err := ctx.Bind(s)
		if err != nil {
			return err
		}
		if len(s.Ids) > 0 {
			return nil
		}
	}
	return ctx.BindUri(s)
}

func (s *SysJobById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysJobById) GenerateM() (models2.ActiveRecord, error) {
	return &models.SysJob{}, nil
}

func (s *SysJobById) GetId() interface{} {
	if len(s.Ids) > 0 {
		return s.Ids
	}
	return s.Id
}
