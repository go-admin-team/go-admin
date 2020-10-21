package dto

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/admin/models/system"
	"gorm.io/gorm"

	"go-admin/common/dto"
	"go-admin/common/log"
	common "go-admin/common/models"
	"go-admin/tools"
)

type SysConfigSearch struct {
	dto.Pagination `search:"-"`
	ConfigName     string `form:"configName" search:"type:exact;column:config_name;table:sys_config" comment:""`
	ConfigKey      string `form:"configKey" search:"type:exact;column:config_key;table:sys_config" comment:""`
	ConfigType     string `form:"configType" search:"type:exact;column:config_type;table:sys_config" comment:""`
}

func (m *SysConfigSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysConfigSearch) Bind(ctx *gin.Context) error {
	msgID := tools.GenerateMsgIDFromContext(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("MsgID[%s] ShouldBind error: %s", msgID, err.Error())
	}
	return err
}

func (m *SysConfigSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysConfigControl struct {
	ID          uint   `uri:"ID" comment:"编码"` // 编码
	ConfigName  string `json:"configName" comment:""`
	ConfigKey   string `json:"configKey" comment:""`
	ConfigValue string `json:"configValue" comment:""`
	ConfigType  string `json:"configType" comment:""`
	Remark      string `json:"remark" comment:""`
}

func (s *SysConfigControl) Bind(ctx *gin.Context) error {
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

func (s *SysConfigControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysConfigControl) GenerateM() (common.ActiveRecord, error) {
	return &system.SysConfig{
		Model:       gorm.Model{ID: s.ID},
		ConfigName:  s.ConfigName,
		ConfigKey:   s.ConfigKey,
		ConfigValue: s.ConfigValue,
		ConfigType:  s.ConfigType,
		Remark:      s.Remark,
	}, nil
}

func (s *SysConfigControl) GetId() interface{} {
	return s.ID
}

type SysConfigById struct {
	dto.ObjectById
}

func (s *SysConfigById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysConfigById) GenerateM() (common.ActiveRecord, error) {
	return &system.SysConfig{}, nil
}
