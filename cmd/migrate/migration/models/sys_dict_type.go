package models

type DictType struct {
	DictId   int    `gorm:"primaryKey;autoIncrement;" json:"dictId"`
	DictName string `gorm:"size:128;" json:"dictName"` //字典名称
	DictType string `gorm:"size:128;" json:"dictType"` //字典类型
	Status   int    `gorm:"size:4;" json:"status"`     //状态
	Remark   string `gorm:"size:255;" json:"remark"`   //备注
	ControlBy
	ModelTime
}

func (DictType) TableName() string {
	return "sys_dict_type"
}
