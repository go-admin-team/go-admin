package models

import (
	orm "go-admin/global"
)

type SysSetting struct {
	SettingsId int    `json:"settings_id" gorm:"primary_key;AUTO_INCREMENT"`
	Name   string `json:"name" gorm:"type:varchar(256);"`
	Logo      string `json:"logo" gorm:"type:varchar(256);"`
	BaseModel
}

func (SysSetting) TableName() string {
	return "sys_setting"
}

//查询
func (s *SysSetting) Get() (create SysSetting, err error) {
	result := orm.Eloquent.Table("sys_setting").First(&create)
	if result.Error != nil {
		err = result.Error
		return
	}
	return create,nil
}

//修改
func (s *SysSetting) Update() (update SysSetting, err error) {
	if err = orm.Eloquent.Table("sys_setting").Model(&update).Updates(&s).Error; err != nil {
		return
	}
	return
}

type ResponseSystemConfig struct {
	Name    string `json:"name" binding:"required"`    // 名称
	Logo    string `json:"logo" binding:"required"`     // 头像
}