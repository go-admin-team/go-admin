package models

import (
	"time"

	"go-admin/common/models"
)

type OrderReconciliation struct {
	models.Model

	ProductName                   string    `json:"productName" gorm:"type:varchar(255);comment:商品名称"`
	ProductCategoryId             string    `json:"productCategoryId" gorm:"type:int(11);comment:商品分类ID"`
	ProductCategoryName           string    `json:"productCategoryName" gorm:"type:varchar(255);comment:商品分类名称"`
	PaymentTypeId                 string    `json:"paymentTypeId" gorm:"type:int(11);comment:支付类型ID"`
	PaymentTypeName               string    `json:"paymentTypeName" gorm:"type:varchar(255);comment:支付类型名称"`
	TransactionSceneId            string    `json:"transactionSceneId" gorm:"type:int(11);comment:交易场景ID"`
	TransactionSceneName          string    `json:"transactionSceneName" gorm:"type:varchar(255);comment:交易场景名称"`
	OrderNumber                   string    `json:"orderNumber" gorm:"type:varchar(100);comment:订单号"`
	TransactionSerialNumber       string    `json:"transactionSerialNumber" gorm:"type:varchar(100);comment:交易流水号"`
	CentralSettlementNumber       string    `json:"centralSettlementNumber" gorm:"type:varchar(100);comment:中联结算号"`
	BankOrderNumber               string    `json:"bankOrderNumber" gorm:"type:varchar(100);comment:银行订单号"`
	UserId                        string    `json:"userId" gorm:"type:int(11);comment:用户ID"`
	UserName                      string    `json:"userName" gorm:"type:varchar(255);comment:用户名称"`
	Openid                        string    `json:"openid" gorm:"type:varchar(100);comment:OpenID"`
	UserAccount                   string    `json:"userAccount" gorm:"type:varchar(100);comment:用户账户"`
	SystemOrderAmount             string    `json:"systemOrderAmount" gorm:"type:decimal(10,2);comment:系统订单金额"`
	MerchantOrderAmount           string    `json:"merchantOrderAmount" gorm:"type:decimal(10,2);comment:商家订单金额"`
	CentralAmount                 string    `json:"centralAmount" gorm:"type:decimal(10,2);comment:中联金额"`
	SystemRefundAmount            string    `json:"systemRefundAmount" gorm:"type:decimal(10,2);comment:系统退款金额"`
	MerchantRefundAmount          string    `json:"merchantRefundAmount" gorm:"type:decimal(10,2);comment:商家退款金额"`
	MerchantActualRefundAmount    string    `json:"merchantActualRefundAmount" gorm:"type:decimal(10,2);comment:商家实退金额"`
	MerchantTransactionStatusId   string    `json:"merchantTransactionStatusId" gorm:"type:int(11);comment:商家交易状态ID"`
	MerchantTransactionStatusName string    `json:"merchantTransactionStatusName" gorm:"type:varchar(255);comment:商家交易状态名称"`
	SystemTransactionStatusId     string    `json:"systemTransactionStatusId" gorm:"type:int(11);comment:系统交易状态ID"`
	SystemTransactionStatusName   string    `json:"systemTransactionStatusName" gorm:"type:varchar(255);comment:系统交易状态名称"`
	CentralTransactionStatusId    string    `json:"centralTransactionStatusId" gorm:"type:int(11);comment:中联交易状态ID"`
	CentralTransactionStatusName  string    `json:"centralTransactionStatusName" gorm:"type:varchar(255);comment:中联交易状态名称"`
	CentralSettlementMethodId     string    `json:"centralSettlementMethodId" gorm:"type:int(11);comment:中联结算方式ID"`
	CentralSettlementMethodName   string    `json:"centralSettlementMethodName" gorm:"type:varchar(255);comment:中联结算方式名称"`
	OrderCreationTime             time.Time `json:"orderCreationTime" gorm:"type:datetime;comment:订单创建时间"`
	TransactionTime               time.Time `json:"transactionTime" gorm:"type:datetime;comment:交易时间"`
	OperatorName                  string    `json:"operatorName" gorm:"type:varchar(255);comment:操作员姓名"`
	models.ModelTime
	models.ControlBy
}

func (OrderReconciliation) TableName() string {
	return "t_order_reconciliation"
}

func (e *OrderReconciliation) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *OrderReconciliation) GetId() interface{} {
	return e.Id
}
