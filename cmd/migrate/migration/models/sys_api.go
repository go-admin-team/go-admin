package models

type SysApi struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	Handle string `json:"handle" gorm:"size:128;comment:handle"`
	Title  string `json:"title" gorm:"size:128;comment:标题"`
	Path   string `json:"path" gorm:"size:128;comment:地址"`
	Type   string `json:"type" gorm:"size:16;comment:接口类型"`
	Action string `json:"action" gorm:"size:16;comment:请求类型"`
	ModelTime
	ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}