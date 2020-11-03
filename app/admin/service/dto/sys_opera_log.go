package dto

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/tools"

	"time"
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
	msgID := tools.GenerateMsgIDFromContext(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
	}
	return err
}

func (m *SysOperaLogSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysOperaLogControl struct {
	ID            uint      `uri:"ID" comment:"编码"` // 编码
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
	msgID := tools.GenerateMsgIDFromContext(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBindUri error: %s", msgID, err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %#v", msgID, err.Error())
	}
	return err
}

func (s *SysOperaLogControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysOperaLogControl) GenerateM() (common.ActiveRecord, error) {
	return &system.SysOperaLog{
		Model:         gorm.Model{ID: s.ID},
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
	dto.ObjectById
}

func (s *SysOperaLogById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysOperaLogById) GenerateM() (common.ActiveRecord, error) {
	return &system.SysOperaLog{}, nil
}
