package models

type SysChinaAreaData struct {
	Model
	PId  string `json:"p_id" gorm:"size:11;comment:上级编码"`
	Name string `json:"name" gorm:"size:128;comment:名称"`
	ControlBy
	ModelTime
}

func (SysChinaAreaData) TableName() string {
	return "sys_china_area_data"
}