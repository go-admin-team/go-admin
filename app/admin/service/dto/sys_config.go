package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

// SysConfigGetPageReq 列表或者搜索使用结构体
type SysConfigGetPageReq struct {
	dto.Pagination `search:"-"`
	ConfigName     string `form:"configName" search:"type:contains;column:config_name;table:sys_config"`
	ConfigKey      string `form:"configKey" search:"type:contains;column:config_key;table:sys_config"`
	ConfigType     string `form:"configType" search:"type:exact;column:config_type;table:sys_config"`
	IsFrontend     int    `form:"isFrontend" search:"type:exact;column:is_frontend;table:sys_config"`
	SysConfigOrder
}

type SysConfigOrder struct {
	IdOrder         string `search:"type:order;column:id;table:sys_config" form:"idOrder"`
	ConfigNameOrder string `search:"type:order;column:config_name;table:sys_config" form:"configNameOrder"`
	ConfigKeyOrder  string `search:"type:order;column:config_key;table:sys_config" form:"configKeyOrder"`
	ConfigTypeOrder string `search:"type:order;column:config_type;table:sys_config" form:"configTypeOrder"`
	CreatedAtOrder  string `search:"type:order;column:created_at;table:sys_config" form:"createdAtOrder"`
}

func (m *SysConfigGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysConfigGetToSysAppReq struct {
	IsFrontend     int    `form:"isFrontend" search:"type:exact;column:is_frontend;table:sys_config"`
}

func (m *SysConfigGetToSysAppReq) GetNeedSearch() interface{} {
	return *m
}

// SysConfigControl 增、改使用的结构体
type SysConfigControl struct {
	Id          int    `uri:"Id" comment:"编码"` // 编码
	ConfigName  string `json:"configName" comment:""`
	ConfigKey   string `uri:"configKey" json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
	ConfigType  string `json:"configType" comment:""`
	IsFrontend  int    `json:"isFrontend"`
	Remark      string `json:"remark" comment:""`
	common.ControlBy
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysConfigControl) Generate(model *models.SysConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ConfigName = s.ConfigName
	model.ConfigKey = s.ConfigKey
	model.ConfigValue = s.ConfigValue
	model.ConfigType = s.ConfigType
	model.IsFrontend = s.IsFrontend
	model.Remark = s.Remark

}

// GetId 获取数据对应的ID
func (s *SysConfigControl) GetId() interface{} {
	return s.Id
}

// GetSetSysConfigReq 增、改使用的结构体
type GetSetSysConfigReq struct {
	ConfigKey   string `json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *GetSetSysConfigReq) Generate(model *models.SysConfig) {
	model.ConfigValue = s.ConfigValue
}

type UpdateSetSysConfigReq map[string]string

// SysConfigByKeyReq 根据Key获取配置
type SysConfigByKeyReq struct {
	ConfigKey string `uri:"configKey" search:"type:contains;column:config_key;table:sys_config"`
}

func (m *SysConfigByKeyReq) GetNeedSearch() interface{} {
	return *m
}

type GetSysConfigByKEYForServiceResp struct {
	ConfigKey   string `json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
}

type SysConfigGetReq struct {
	Id int `uri:"id"`
}

func (s *SysConfigGetReq) GetId() interface{} {
	return s.Id
}

type SysConfigDeleteReq struct {
	Ids []int `json:"ids"`
	common.ControlBy
}

func (s *SysConfigDeleteReq) GetId() interface{} {
	return s.Ids
}
