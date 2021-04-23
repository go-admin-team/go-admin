package models

import (
	"time"
)

type SysLoginLog struct {
	Model
	Username      string    `json:"username" gorm:"type:varchar(128);comment:用户名"`
	Status        string    `json:"status" gorm:"type:varchar(4);comment:状态"`
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(255);comment:ip地址"`
	LoginLocation string    `json:"loginLocation" gorm:"type:varchar(255);comment:归属地"`
	Browser       string    `json:"browser" gorm:"type:varchar(255);comment:浏览器"`
	Os            string    `json:"os" gorm:"type:varchar(255);comment:系统"`
	Platform      string    `json:"platform" gorm:"type:varchar(255);comment:固件"`
	LoginTime     time.Time `json:"loginTime" gorm:"type:timestamp;comment:登录时间"`
	Remark        string    `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Msg           string    `json:"msg" gorm:"type:varchar(255);comment:信息"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	ControlBy
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
