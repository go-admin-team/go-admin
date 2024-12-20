package models

import (
	"go-admin/common/models"
	"time"
)

type WxScan struct {
	models.Model
	Time           time.Time `json:"time" gorm:"type:datetime;comment:请求时间"`
	Appid          string    `json:"appid" gorm:"column:appid;comment:微信分配的公众账号ID;type:varchar(32);not null" mapstructure:"appid"`
	MchID          string    `json:"mch_id" gorm:"column:mch_id;comment:微信支付分配的商户号;type:varchar(32);not null" mapstructure:"mch_id"`
	DeviceInfo     string    `json:"device_info" gorm:"column:device_info;comment:终端设备号;type:varchar(16)" mapstructure:"device_info"`
	NonceStr       string    `json:"nonce_str" gorm:"column:nonce_str;comment:随机字符串;type:varchar(64);not null" mapstructure:"nonce_str"`
	Body           string    `json:"body" gorm:"column:body;comment:商品描述;type:varchar(128);not null" mapstructure:"body"`
	Attach         string    `json:"attach" gorm:"column:attach;comment:附加数据;type:varchar(128)" mapstructure:"attach"`
	OutTradeNo     string    `json:"out_trade_no" gorm:"column:out_trade_no;comment:商户订单号;type:varchar(64);not null;unique_index:uniq_out_trade_no" mapstructure:"out_trade_no"`
	TotalFee       int       `json:"total_fee" gorm:"column:total_fee;comment:订单金额，单位为分;type:int;not null" mapstructure:"total_fee"`
	SpbillCreateIP string    `json:"spbill_create_ip" gorm:"column:spbill_create_ip;comment:APP和网页支付提交用户端ip;type:varchar(15);not null" mapstructure:"spbill_create_ip"`
	AuthCode       string    `json:"auth_code" gorm:"column:auth_code;comment:扫码支付授权码;type:varchar(64);not null" mapstructure:"auth_code"`
	Sign           string    `json:"sign" gorm:"column:sign;comment:签名;type:varchar(64);not null" mapstructure:"sign"`
}

func (*WxScan) TableName() string {
	return "m_wx_scan"
}
