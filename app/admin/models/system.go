package models

import (
	"go-admin/common/models"
)

type SysSetting struct {
	SettingsId int    `json:"settings_id" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string `json:"name" gorm:"type:varchar(256);"`
	Logo       string `json:"logo" gorm:"type:varchar(256);"`
	models.ModelTime
}

func (SysSetting) TableName() string {
	return "sys_setting"
}

func (s *SysSetting) GetId() interface{} {
	return s.SettingsId
}

//查询
//func (s *SysSetting) Get() (create SysSetting, err error) {
//	result := orm.Eloquent.Table("sys_setting").First(&create)
//	if result.Error != nil {
//		err = result.Error
//		return
//	}
//	return create, nil
//}

//修改
//func (s *SysSetting) Update() (update SysSetting, err error) {
//	if err = orm.Eloquent.Table("sys_setting").Model(&update).Where("settings_id = ?", s.SettingsId).Updates(&s).Error; err != nil {
//		return
//	}
//	return
//}

type ResponseSystemConfig struct {
	Name       string `json:"name" binding:"required"`        // 名称
	Logo       string `json:"logo" binding:"required"`        // 头像
	SettingsId int    `json:"settings_id" binding:"required"` // 头像
}
