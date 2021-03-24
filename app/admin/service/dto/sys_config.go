package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

// SysConfigSearch 列表或者搜索使用结构体
type SysConfigSearch struct {
	dto.Pagination `search:"-"`
	ConfigName     string `form:"configName" search:"type:contains;column:config_name;table:sys_config" comment:""`
	ConfigKey      string `form:"configKey" search:"type:contains;column:config_key;table:sys_config" comment:""`
	ConfigType     string `form:"configType" search:"type:exact;column:config_type;table:sys_config" comment:""`
	IsFrontend     int    `json:"isFrontend" search:"type:exact;column:is_frontend;table:sys_config" comment:""`
}

func (m *SysConfigSearch) GetNeedSearch() interface{} {
	return *m
}

// Bind 映射上下文中的结构体数据
func (m *SysConfigSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

// SysConfigControl 增、改使用的结构体
type SysConfigControl struct {
	ID          int    `uri:"Id" comment:"编码"` // 编码
	ConfigName  string `json:"configName" comment:""`
	ConfigKey   string `uri:"configKey" json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
	ConfigType  string `json:"configType" comment:""`
	IsFrontend  int    `json:"isFrontend"`
	Remark      string `json:"remark" comment:""`
}

// Bind 映射上下文中的结构体数据
func (s *SysConfigControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBindBodyWith(s, binding.JSON)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	var jsonStr []byte
	jsonStr, err = json.Marshal(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	ctx.Set("body", string(jsonStr))
	return err
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysConfigControl) Generate() (*system.SysConfig, error) {
	return &system.SysConfig{
		Model:       common.Model{Id: s.ID},
		ConfigName:  s.ConfigName,
		ConfigKey:   s.ConfigKey,
		ConfigValue: s.ConfigValue,
		ConfigType:  s.ConfigType,
		IsFrontend:  s.IsFrontend,
		Remark:      s.Remark,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysConfigControl) GetId() interface{} {
	return s.ID
}

// SysConfigById 获取单个或者删除的结构体
type SysConfigById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysConfigById) Generate() *SysConfigById {
	cp := *s
	return &cp
}

func (s *SysConfigById) GetId() interface{} {
	return s.Id
}

func (s *SysConfigById) Bind(ctx *gin.Context) error {
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

func (s *SysConfigById) GenerateM() (*system.SysConfig, error) {
	return &system.SysConfig{}, nil
}
