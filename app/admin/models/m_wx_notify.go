package models

import (
	"go-admin/common/models"
	"time"
)

// WxNotify 结构体用于接收和处理微信支付回调的数据
type WxNotify struct {
	models.Model
	Mchid          string    `json:"mchid" gorm:"column:mchid;comment:商户号;type:varchar(20)" mapstructure:"mchid"`
	Appid          string    `json:"appid" gorm:"column:appid;comment:应用ID;type:varchar(30)" mapstructure:"appid"`
	OutTradeNo     string    `json:"out_trade_no" gorm:"column:out_trade_no;comment:商户订单号;type:varchar(64)" mapstructure:"out_trade_no"`
	TransactionID  string    `json:"transaction_id" gorm:"column:transaction_id;comment:微信支付订单号;type:varchar(64)" mapstructure:"transaction_id"`
	TradeType      string    `json:"trade_type" gorm:"column:trade_type;comment:交易类型;type:varchar(10)" mapstructure:"trade_type"`
	TradeState     string    `json:"trade_state" gorm:"column:trade_state;comment:交易状态;type:varchar(10)" mapstructure:"trade_state"`
	TradeStateDesc string    `json:"trade_state_desc" gorm:"column:trade_state_desc;comment:交易状态描述;type:varchar(255)" mapstructure:"trade_state_desc"`
	BankType       string    `json:"bank_type" gorm:"column:bank_type;comment:银行类型;type:varchar(20)" mapstructure:"bank_type"`
	Attach         string    `json:"attach" gorm:"column:attach;comment:自定义数据;type:varchar(255)" mapstructure:"attach"`
	SuccessTime    time.Time `json:"success_time" gorm:"column:success_time;comment:成功时间;type:datetime" mapstructure:"success_time"`
	Openid         string    `json:"openid" gorm:"column:openid;comment:用户标识;type:varchar(64)" mapstructure:"openid"`
	Total          int       `json:"total" gorm:"column:total;comment:订单金额（单位：分）;type:int(11)" mapstructure:"total"`
	PayerTotal     int       `json:"payer_total" gorm:"column:payer_total;comment:付款金额（单位：分）;type:int(11)" mapstructure:"payer_total"`
	Currency       string    `json:"currency" gorm:"column:currency;comment:货币种类;type:varchar(10)" mapstructure:"currency"`
	PayerCurrency  string    `json:"payer_currency" gorm:"column:payer_currency;comment:付款货币种类;type:varchar(10)" mapstructure:"payer_currency"`
	TypeName       string    `json:"type_name" gorm:"column:type_name;type:varchar(50);comment:业务类型" mapstructure:"type_name"`
}

func (*WxNotify) TableName() string {
	return "m_wx_notify"
}
