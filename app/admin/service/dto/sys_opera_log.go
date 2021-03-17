package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysOperaLogSearch struct {
	dto.Pagination `search:"-"`

	Title         string `form:"title" search:"type:contains;column:title;table:sys_opera_log" comment:"操作模块"`
	Method        string `form:"method" search:"type:contains;column:method;table:sys_opera_log" comment:"函数"`
	RequestMethod string `form:"requestMethod" search:"type:contains;column:request_method;table:sys_opera_log" comment:"请求方式"`
	OperUrl       string `form:"operUrl" search:"type:contains;column:oper_url;table:sys_opera_log" comment:"访问地址"`
	OperIp        string `form:"operIp" search:"type:exact;column:oper_ip;table:sys_opera_log" comment:"客户端ip"`
}

func (m *SysOperaLogSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysOperaLogSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

type SysOperaLogControl struct {
	ID            int       `uri:"Id" comment:"编码"` // 编码
	Title         string    `json:"title" comment:"操作模块"`
	BusinessType  string    `json:"businessType" comment:"操作类型"`
	BusinessTypes string    `json:"businessTypes" comment:""`
	Method        string    `json:"method" comment:"函数"`
	RequestMethod string    `json:"requestMethod" comment:"请求方式"`
	OperatorType  string    `json:"operatorType" comment:"操作类型"`
	OperName      string    `json:"operName" comment:"操作者"`
	DeptName      string    `json:"deptName" comment:"部门名称"`
	OperUrl       string    `json:"operUrl" comment:"访问地址"`
	OperIp        string    `json:"operIp" comment:"客户端ip"`
	OperLocation  string    `json:"operLocation" comment:"访问位置"`
	OperParam     string    `json:"operParam" comment:"请求参数"`
	Status        string    `json:"status" comment:"操作状态"`
	OperTime      time.Time `json:"operTime" comment:"操作时间"`
	JsonResult    string    `json:"jsonResult" comment:"返回数据"`
	Remark        string    `json:"remark" comment:"备注"`
	LatencyTime   string    `json:"latencyTime" comment:"耗时"`
	UserAgent     string    `json:"userAgent" comment:"ua"`
}

func (s *SysOperaLogControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysOperaLogControl) Generate() (*system.SysOperaLog, error) {
	return &system.SysOperaLog{
		Model:         common.Model{Id: s.ID},
		Title:         s.Title,
		BusinessType:  s.BusinessType,
		BusinessTypes: s.BusinessTypes,
		Method:        s.Method,
		RequestMethod: s.RequestMethod,
		OperatorType:  s.OperatorType,
		OperName:      s.OperName,
		DeptName:      s.DeptName,
		OperUrl:       s.OperUrl,
		OperIp:        s.OperIp,
		OperLocation:  s.OperLocation,
		OperParam:     s.OperParam,
		Status:        s.Status,
		OperTime:      s.OperTime,
		JsonResult:    s.JsonResult,
		Remark:        s.Remark,
		LatencyTime:   s.LatencyTime,
		UserAgent:     s.UserAgent,
	}, nil
}

func (s *SysOperaLogControl) GetId() interface{} {
	return s.ID
}

type SysOperaLogById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysOperaLogById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysOperaLogById) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(&s.Ids)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysOperaLogById) SetUpdateBy(id int) {

}
