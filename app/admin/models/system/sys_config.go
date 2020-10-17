package system

import (
	"gorm.io/gorm"

	"go-admin/common/models"
)

type SysConfig struct {
	gorm.Model
	models.ControlBy
	ConfigName  string `json:"configName" gorm:"type:varchar(128);comment:ConfigName"`   //
	ConfigKey   string `json:"configKey" gorm:"type:varchar(128);comment:ConfigKey"`     //
	ConfigValue string `json:"configValue" gorm:"type:varchar(255);comment:ConfigValue"` //
	ConfigType  string `json:"configType" gorm:"type:varchar(64);comment:ConfigType"`    //
	Remark      string `json:"remark" gorm:"type:varchar(128);comment:Remark"`           //
}

func (SysConfig) TableName() string {
	return "sys_config"
}

func (e *SysConfig) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysConfig) GetId() interface{} {
	return e.ID
}
