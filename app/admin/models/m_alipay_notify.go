package models

import (
	"go-admin/common/models"
	"time"
)

// AlipayNotify 支付宝支付通知
type AlipayNotify struct {
	models.Model

	// GmtCreate 代表交易创建的时间
	GmtCreate time.Time `json:"gmt_create" gorm:"column:gmt_create;comment:交易创建的时间" mapstructure:"gmt_create"`

	//将以下字段都加入tag为map

	// Charset 指定字符集编码
	Charset string `json:"charset" gorm:"column:charset;type:varchar(10);comment:字符集编码" mapstructure:"charset"`

	// SellerEmail 是卖家的电子邮件地址
	SellerEmail string `json:"seller_email" gorm:"column:seller_email;type:varchar(100);comment:卖家的电子邮件地址" mapstructure:"seller_email"`

	// Subject 描述交易的主题或商品描述
	Subject string `json:"subject" gorm:"column:subject;type:varchar(255);comment:交易的主题或商品描述" mapstructure:"subject"`

	// Sign 用于验证消息真实性的签名
	Sign string `json:"sign" gorm:"column:sign;type:text;comment:用于验证消息真实性的签名" mapstructure:"sign"`

	// BuyerOpenID 是买家的支付宝Open ID
	BuyerOpenID string `json:"buyer_open_id" gorm:"column:buyer_open_id;type:varchar(100);comment:买家的支付宝Open ID" mapstructure:"buyer_open_id"`

	// InvoiceAmount 是发票金额
	InvoiceAmount float64 `json:"invoice_amount" gorm:"column:invoice_amount;type:decimal(10,2);comment:发票金额" mapstructure:"invoice_amount"`

	// NotifyID 是支付宝分配给此次通知的唯一ID
	NotifyID string `json:"notify_id" gorm:"column:notify_id;type:varchar(100);comment:支付宝分配给此次通知的唯一ID" mapstructure:"notify_id"`

	// FundBillList 包含资金账单的列表，列出支付方式和金额
	FundBillList string `json:"fund_bill_list" gorm:"column:fund_bill_list;type:json;comment:包含资金账单的列表，列出支付方式和金额" mapstructure:"fund_bill_list"`

	// NotifyType 表示通知的类型
	NotifyType string `json:"notify_type" gorm:"column:notify_type;type:varchar(50);comment:表示通知的类型" mapstructure:"notify_type"`

	// TradeStatus 显示当前交易的状态
	TradeStatus string `json:"trade_status" gorm:"column:trade_status;type:varchar(50);comment:显示当前交易的状态" mapstructure:"trade_status"`

	// ReceiptAmount 是实际收款金额
	ReceiptAmount float64 `json:"receipt_amount" gorm:"column:receipt_amount;type:decimal(10,2);comment:实际收款金额" mapstructure:"receipt_amount"`

	// BuyerPayAmount 是买家实际支付的金额
	BuyerPayAmount float64 `json:"buyer_pay_amount" gorm:"column:buyer_pay_amount;type:decimal(10,2);comment:买家实际支付的金额" mapstructure:"buyer_pay_amount"`

	// AppID 是商户的应用ID
	AppID string `json:"app_id" gorm:"column:app_id;type:varchar(50);comment:商户的应用ID" mapstructure:"app_id`

	// SignType 是签名算法类型
	SignType string `json:"sign_type" gorm:"column:sign_type;type:varchar(10);comment:签名算法类型" mapstructure:"sign_type"`

	// SellerID 是卖家的支付宝用户ID
	SellerID string `json:"seller_id" gorm:"column:seller_id;type:varchar(50);comment:卖家的支付宝用户ID" mapstructure:"seller_id"`

	// GmtPayment 是实际支付发生的时间
	GmtPayment time.Time `json:"gmt_payment" gorm:"column:gmt_payment;comment:实际支付发生的时间" mapstructure:"gmt_payment"`

	// NotifyTime 是支付宝发送此通知的时间
	NotifyTime time.Time `json:"notify_time" gorm:"column:notify_time;comment:支付宝发送此通知的时间" mapstructure:"notify_time"`

	// MerchantAppID 是商户应用ID，与 AppID 相同
	MerchantAppID string `json:"merchant_app_id" gorm:"column:merchant_app_id;type:varchar(50);comment:商户应用ID，与 AppID 相同" mapstructure:"merchant_app_id"`

	// Version 是接口版本号
	Version string `json:"version" gorm:"column:version;type:varchar(10);comment:接口版本号" mapstructure:"version"`

	// OutTradeNo 是商户系统内部的订单号
	OutTradeNo string `json:"out_trade_no" gorm:"column:out_trade_no;type:varchar(100);comment:商户系统内部的订单号" mapstructure:"out_trade_no"`

	// TotalAmount 是订单总金额
	TotalAmount float64 `json:"total_amount" gorm:"column:total_amount;type:decimal(10,2);comment:订单总金额" mapstructure:"total_amount"`

	// TradeNo 是支付宝交易号，用于唯一识别一笔交易
	TradeNo string `json:"trade_no" gorm:"column:trade_no;type:varchar(100);comment:支付宝交易号，用于唯一识别一笔交易" mapstructure:"trade_no"`

	// AuthAppID 是授权应用ID，通常与 AppID 相同
	AuthAppID string `json:"auth_app_id" gorm:"column:auth_app_id;type:varchar(50);comment:授权应用ID，通常与 AppID 相同" mapstructure:"auth_app_id"`

	// BuyerLogonID 是买家的登录ID，部分信息可能被屏蔽以保护隐私
	BuyerLogonID string `json:"buyer_logon_id" gorm:"column:buyer_logon_id;type:varchar(50);comment:买家的登录ID，部分信息可能被屏蔽以保护隐私" mapstructure:"buyer_logon_id"`

	// TypeName 业务类型
	TypeName string `json:"type_name" gorm:"column:type_name;type:varchar(50);comment:业务类型" mapstructure:"type_name"`

	// PointAmount 是积分抵扣的金额
	PointAmount float64 `json:"point_amount" gorm:"column:point_amount;type:decimal(10,2);comment:积分抵扣的金额" mapstructure:"point_amount"`
}

func (*AlipayNotify) TableName() string {
	return "m_alipay_notify"
}
