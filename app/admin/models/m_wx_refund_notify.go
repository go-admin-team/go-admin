package models

import (
	"go-admin/common/models"
	"time"
)

// WxRefundNotify 结构体用于接收和处理微信支付退款回调的数据
type WxRefundNotify struct {
	models.Model
	Mchid               string    `json:"mchid" gorm:"column:mchid;comment:商户号;type:varchar(20)" mapstructure:"mchid"`
	OutTradeNo          string    `json:"out_trade_no" gorm:"column:out_trade_no;comment:商户订单号;type:varchar(64)" mapstructure:"out_trade_no"`
	TransactionID       string    `json:"transaction_id" gorm:"column:transaction_id;comment:微信支付订单号;type:varchar(64)" mapstructure:"transaction_id"`
	OutRefundNo         string    `json:"out_refund_no" gorm:"column:out_refund_no;comment:商户退款单号;type:varchar(64)" mapstructure:"out_refund_no"`
	RefundID            string    `json:"refund_id" gorm:"column:refund_id;comment:微信退款单号;type:varchar(64)" mapstructure:"refund_id"`
	RefundStatus        string    `json:"refund_status" gorm:"column:refund_status;comment:退款状态;type:varchar(10)" mapstructure:"refund_status"`
	SuccessTime         time.Time `json:"success_time" gorm:"column:success_time;comment:退款成功时间;type:datetime" mapstructure:"success_time"`
	Total               int       `json:"total" gorm:"column:total;comment:订单金额（单位：分）;type:int(11)" mapstructure:"total"`
	Refund              int       `json:"refund" gorm:"column:refund;comment:退款金额（单位：分）;type:int(11)" mapstructure:"refund"`
	PayerTotal          int       `json:"payer_total" gorm:"column:payer_total;comment:付款金额（单位：分）;type:int(11)" mapstructure:"payer_total"`
	PayerRefund         int       `json:"payer_refund" gorm:"column:payer_refund;comment:付款方退款金额（单位：分）;type:int(11)" mapstructure:"payer_refund"`
	UserReceivedAccount string    `json:"user_received_account" gorm:"column:user_received_account;comment:收款账户;type:varchar(255)" mapstructure:"user_received_account"`
	TypeName            string    `json:"type_name" gorm:"column:type_name;type:varchar(50);comment:业务类型" mapstructure:"type_name"`
}

func (*WxRefundNotify) TableName() string {
	return "m_wx_refund_notify"
}
