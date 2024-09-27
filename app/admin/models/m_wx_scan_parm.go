package models

import (
	"go-admin/common/models"
	"time"
)

// WxScanParm 代表微信扫码支付-第三方传来的参数模型
type WxScanParm struct {
	models.Model
	AuthCode string    `json:"auth_code" gorm:"column:auth_code;type:varchar(255);not null;comment:'付款码'" mapstructure:"auth_code"`
	Body     string    `json:"body" gorm:"column:body;type:varchar(255);not null;comment:'商品描述'" mapstructure:"body"`
	OrderID  string    `json:"order_id" gorm:"column:order_id;type:varchar(255);not null;comment:'订单号'" mapstructure:"order_id"`
	PayType  int8      `json:"pay_type" gorm:"column:pay_type;type:tinyint;not null;comment:'扫码类型: 1-微信, 2-支付宝'" mapstructure:"pay_type"`
	Total    float64   `json:"total" gorm:"column:total;type:decimal(10,2);not null;comment:'总金额'" mapstructure:"total"`
	UserID   string    `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:'设备信息'" mapstructure:"user_id"`
	Time     time.Time `json:"time" gorm:"type:datetime;comment:通知时间"`
}

// TableName 返回表名
func (WxScanParm) TableName() string {
	return "m_wx_scan_parm"
}
