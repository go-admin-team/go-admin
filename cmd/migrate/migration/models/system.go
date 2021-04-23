package models

type SysSetting struct {
	SettingsId int    `json:"settings_id" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string `json:"name" gorm:"type:varchar(256);"`
	Logo       string `json:"logo" gorm:"type:varchar(256);"`
	ModelTime
}

func (SysSetting) TableName() string {
	return "sys_setting"
}
