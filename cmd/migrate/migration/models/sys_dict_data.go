package models

type DictData struct {
	DictCode  int    `gorm:"primaryKey;autoIncrement;" json:"dictCode" example:"1"` //字典编码
	DictSort  int    `gorm:"" json:"dictSort"`                                      //显示顺序
	DictLabel string `gorm:"size:128;" json:"dictLabel"`                            //数据标签
	DictValue string `gorm:"size:255;" json:"dictValue"`                            //数据键值
	DictType  string `gorm:"size:64;" json:"dictType"`                              //字典类型
	CssClass  string `gorm:"size:128;" json:"cssClass"`                             //
	ListClass string `gorm:"size:128;" json:"listClass"`                            //
	IsDefault string `gorm:"size:8;" json:"isDefault"`                              //
	Status    int    `gorm:"size:4;" json:"status"`                                 //状态
	Default   string `gorm:"size:8;" json:"default"`                                //
	Remark    string `gorm:"size:255;" json:"remark"`                               //备注
	ControlBy
	ModelTime
}

func (DictData) TableName() string {
	return "sys_dict_data"
}
