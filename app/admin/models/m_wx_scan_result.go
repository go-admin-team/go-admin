package models

import (
	"go-admin/common/models"
	"time"
)

// WxScanResult represents the structure for storing WeChat scan payment results.
type WxScanResult struct {
	models.Model
	ReturnCode    string    `gorm:"column:return_code;type:varchar(16);not null;comment:'返回状态码'" json:"return_code" mapstructure:"return_code"`
	ReturnMsg     string    `gorm:"column:return_msg;type:varchar(128);not null;comment:'返回信息'" json:"return_msg" mapstructure:"return_msg"`
	ResultCode    string    `gorm:"column:result_code;type:varchar(16);not null;comment:'业务结果'" json:"result_code" mapstructure:"result_code"`
	MchID         string    `gorm:"column:mch_id;type:varchar(32);not null;comment:'商户号'" json:"mch_id" mapstructure:"mch_id"`
	AppID         string    `gorm:"column:appid;type:varchar(32);not null;comment:'应用ID'" json:"appid" mapstructure:"appid"`
	DeviceInfo    string    `gorm:"column:device_info;type:varchar(32);not null;comment:'设备号'" json:"device_info" mapstructure:"device_info"`
	NonceStr      string    `gorm:"column:nonce_str;type:varchar(64);not null;comment:'随机字符串'" json:"nonce_str" mapstructure:"nonce_str"`
	Sign          string    `gorm:"column:sign;type:varchar(64);not null;comment:'签名'" json:"sign" mapstructure:"sign"`
	OpenID        string    `gorm:"column:openid;type:varchar(128);not null;comment:'用户标识'" json:"openid" mapstructure:"openid"`
	IsSubscribe   string    `gorm:"column:is_subscribe;type:char(1);not null;comment:'是否关注公众号'" json:"is_subscribe" mapstructure:"is_subscribe"`
	TradeType     string    `gorm:"column:trade_type;type:varchar(32);not null;comment:'交易类型'" json:"trade_type" mapstructure:"trade_type"`
	BankType      string    `gorm:"column:bank_type;type:varchar(32);not null;comment:'付款银行'" json:"bank_type" mapstructure:"bank_type"`
	FeeType       string    `gorm:"column:fee_type;type:varchar(16);not null;comment:'货币种类'" json:"fee_type" mapstructure:"fee_type"`
	TotalFee      int       `gorm:"column:total_fee;type:int;not null;comment:'总金额'" json:"total_fee" mapstructure:"total_fee"`
	CashFeeType   string    `gorm:"column:cash_fee_type;type:varchar(16);not null;comment:'现金支付货币类型'" json:"cash_fee_type" mapstructure:"cash_fee_type"`
	CashFee       int       `gorm:"column:cash_fee;type:int;not null;comment:'现金支付金额'" json:"cash_fee" mapstructure:"cash_fee"`
	TransactionID string    `gorm:"column:transaction_id;type:varchar(64);not null;comment:'微信支付订单号'" json:"transaction_id" mapstructure:"transaction_id"`
	OutTradeNo    string    `gorm:"column:out_trade_no;type:varchar(64);not null;comment:'商户订单号'" json:"out_trade_no" mapstructure:"out_trade_no"`
	Attach        string    `gorm:"column:attach;type:varchar(256);not null;comment:'附加数据'" json:"attach" mapstructure:"attach"`
	TimeEnd       time.Time `gorm:"column:time_end;type:datetime;not null;comment:'交易结束时间'" json:"time_end" mapstructure:"time_end"`
}

// TableName sets the insert table name for this struct type
func (WxScanResult) TableName() string {
	return "m_wx_scan_result"
}

// BeforeCreate hook to parse time_end field if needed.
// func (w *WxScanResult) BeforeCreate(tx *gorm.DB) (err error) {
// 	if w.TimeEnd.IsZero() && w.TimeEnd.String() != "" {
// 		t, err := time.Parse("2006-01-02 15:04:05", w.TimeEnd.String())
// 		if err != nil {
// 			return err
// 		}
// 		w.TimeEnd = t
// 	}
// 	return
// }
