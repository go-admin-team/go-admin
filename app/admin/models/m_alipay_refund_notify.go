package models

import (
	"go-admin/common/models"
	"time"
)

// AlipayRefundNotify 支付宝应退款通知
type AlipayRefundNotify struct {
	models.Model
	Code          string    `json:"code" gorm:"column:code;type:varchar(20);not null;comment:'业务状态码，10000表示成功';" mapstructure:"code"`
	Msg           string    `json:"msg" gorm:"type:varchar(255);not null;comment:'描述性消息，进一步确认请求状态';column:msg" mapstructure:"msg"`
	BuyerLogonID  string    `json:"buyer_logon_id" gorm:"type:varchar(20);not null;comment:'买家登录ID，部分信息可能被屏蔽';column:buyer_logon_id" mapstructure:"buyer_logon_id"`
	FundChange    string    `json:"fund_change" gorm:"type:char(1);not null;comment:'是否涉及资金变化，Y表示有资金变化';column:fund_change" mapstructure:"fund_change"`
	GmtRefundPay  time.Time `json:"gmt_refund_pay" gorm:"type:datetime;not null;comment:'退款发生的具体时间';column:gmt_refund_pay" mapstructure:"gmt_refund_pay"`
	OutTradeNo    string    `json:"out_trade_no" gorm:"type:varchar(50);not null;comment:'商户系统内部的订单号';column:out_trade_no" mapstructure:"out_trade_no"`
	RefundFee     float64   `json:"refund_fee" gorm:"type:decimal(10,2);not null;comment:'退款的总金额';column:refund_fee" mapstructure:"refund_fee"`
	SendBackFee   float64   `json:"send_back_fee" gorm:"type:decimal(10,2);not null;comment:'退回给买家的金额，通常与退款金额相同';column:send_back_fee" mapstructure:"send_back_fee"`
	TradeNo       string    `json:"trade_no" gorm:"type:varchar(50);not null;comment:'支付宝交易号，用于唯一识别这笔退款交易';column:trade_no" mapstructure:"trade_no"`
	BuyerOpenID   string    `json:"buyer_open_id" gorm:"type:varchar(100);not null;comment:'买家的支付宝Open ID';column:buyer_open_id" mapstructure:"buyer_open_id"`
	RefundDetails string    `json:"refund_details" gorm:"type:longtext;not null;comment:'退款详情列表，包含退款金额和资金渠道';column:refund_details" mapstructure:"refund_details"`
	TypeName      string    `json:"type_name" gorm:"column:type_name;type:varchar(50);comment:业务类型" mapstructure:"type_name"`
}

// TableName 设置表名
func (*AlipayRefundNotify) TableName() string {
	return "m_alipay_refund_notify"
}
