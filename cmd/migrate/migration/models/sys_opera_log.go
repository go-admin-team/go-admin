package models

import (
	"time"
)

type SysOperaLog struct {
	Model
	Title         string    `json:"title" gorm:"type:varchar(255);comment:操作模块"`
	BusinessType  string    `json:"businessType" gorm:"type:varchar(128);comment:操作类型"`
	BusinessTypes string    `json:"businessTypes" gorm:"type:varchar(128);comment:BusinessTypes"`
	Method        string    `json:"method" gorm:"type:varchar(128);comment:函数"`
	RequestMethod string    `json:"requestMethod" gorm:"type:varchar(128);comment:请求方式"`
	OperatorType  string    `json:"operatorType" gorm:"type:varchar(128);comment:操作类型"`
	OperName      string    `json:"operName" gorm:"type:varchar(128);comment:操作者"`
	DeptName      string    `json:"deptName" gorm:"type:varchar(128);comment:部门名称"`
	OperUrl       string    `json:"operUrl" gorm:"type:varchar(255);comment:访问地址"`
	OperIp        string    `json:"operIp" gorm:"type:varchar(128);comment:客户端ip"`
	OperLocation  string    `json:"operLocation" gorm:"type:varchar(128);comment:访问位置"`
	OperParam     string    `json:"operParam" gorm:"type:varchar(255);comment:请求参数"`
	Status        string    `json:"status" gorm:"type:varchar(4);comment:操作状态"`
	OperTime      time.Time `json:"operTime" gorm:"type:timestamp;comment:操作时间"`
	JsonResult    string    `json:"jsonResult" gorm:"type:varchar(255);comment:返回数据"`
	Remark        string    `json:"remark" gorm:"type:varchar(255);comment:备注"`
	LatencyTime   string    `json:"latencyTime" gorm:"type:varchar(128);comment:耗时"`
	UserAgent     string    `json:"userAgent" gorm:"type:varchar(255);comment:ua"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	ControlBy
}

func (SysOperaLog) TableName() string {
	return "sys_opera_log"
}
