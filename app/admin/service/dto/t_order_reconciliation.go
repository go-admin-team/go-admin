package dto

import (
	"time"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type OrderReconciliationGetPageReq struct {
	dto.Pagination `search:"-"`
	OrderReconciliationOrder
}

type OrderReconciliationOrder struct {
	Id                            string `form:"idOrder"  search:"type:order;column:id;table:t_order_reconciliation"`
	ProductName                   string `form:"productNameOrder"  search:"type:order;column:product_name;table:t_order_reconciliation"`
	ProductCategoryId             string `form:"productCategoryIdOrder"  search:"type:order;column:product_category_id;table:t_order_reconciliation"`
	ProductCategoryName           string `form:"productCategoryNameOrder"  search:"type:order;column:product_category_name;table:t_order_reconciliation"`
	PaymentTypeId                 string `form:"paymentTypeIdOrder"  search:"type:order;column:payment_type_id;table:t_order_reconciliation"`
	PaymentTypeName               string `form:"paymentTypeNameOrder"  search:"type:order;column:payment_type_name;table:t_order_reconciliation"`
	TransactionSceneId            string `form:"transactionSceneIdOrder"  search:"type:order;column:transaction_scene_id;table:t_order_reconciliation"`
	TransactionSceneName          string `form:"transactionSceneNameOrder"  search:"type:order;column:transaction_scene_name;table:t_order_reconciliation"`
	OrderNumber                   string `form:"orderNumberOrder"  search:"type:order;column:order_number;table:t_order_reconciliation"`
	TransactionSerialNumber       string `form:"transactionSerialNumberOrder"  search:"type:order;column:transaction_serial_number;table:t_order_reconciliation"`
	CentralSettlementNumber       string `form:"centralSettlementNumberOrder"  search:"type:order;column:central_settlement_number;table:t_order_reconciliation"`
	BankOrderNumber               string `form:"bankOrderNumberOrder"  search:"type:order;column:bank_order_number;table:t_order_reconciliation"`
	UserId                        string `form:"userIdOrder"  search:"type:order;column:user_id;table:t_order_reconciliation"`
	UserName                      string `form:"userNameOrder"  search:"type:order;column:user_name;table:t_order_reconciliation"`
	Openid                        string `form:"openidOrder"  search:"type:order;column:openid;table:t_order_reconciliation"`
	UserAccount                   string `form:"userAccountOrder"  search:"type:order;column:user_account;table:t_order_reconciliation"`
	SystemOrderAmount             string `form:"systemOrderAmountOrder"  search:"type:order;column:system_order_amount;table:t_order_reconciliation"`
	MerchantOrderAmount           string `form:"merchantOrderAmountOrder"  search:"type:order;column:merchant_order_amount;table:t_order_reconciliation"`
	CentralAmount                 string `form:"centralAmountOrder"  search:"type:order;column:central_amount;table:t_order_reconciliation"`
	SystemRefundAmount            string `form:"systemRefundAmountOrder"  search:"type:order;column:system_refund_amount;table:t_order_reconciliation"`
	MerchantRefundAmount          string `form:"merchantRefundAmountOrder"  search:"type:order;column:merchant_refund_amount;table:t_order_reconciliation"`
	MerchantActualRefundAmount    string `form:"merchantActualRefundAmountOrder"  search:"type:order;column:merchant_actual_refund_amount;table:t_order_reconciliation"`
	MerchantTransactionStatusId   string `form:"merchantTransactionStatusIdOrder"  search:"type:order;column:merchant_transaction_status_id;table:t_order_reconciliation"`
	MerchantTransactionStatusName string `form:"merchantTransactionStatusNameOrder"  search:"type:order;column:merchant_transaction_status_name;table:t_order_reconciliation"`
	SystemTransactionStatusId     string `form:"systemTransactionStatusIdOrder"  search:"type:order;column:system_transaction_status_id;table:t_order_reconciliation"`
	SystemTransactionStatusName   string `form:"systemTransactionStatusNameOrder"  search:"type:order;column:system_transaction_status_name;table:t_order_reconciliation"`
	CentralTransactionStatusId    string `form:"centralTransactionStatusIdOrder"  search:"type:order;column:central_transaction_status_id;table:t_order_reconciliation"`
	CentralTransactionStatusName  string `form:"centralTransactionStatusNameOrder"  search:"type:order;column:central_transaction_status_name;table:t_order_reconciliation"`
	CentralSettlementMethodId     string `form:"centralSettlementMethodIdOrder"  search:"type:order;column:central_settlement_method_id;table:t_order_reconciliation"`
	CentralSettlementMethodName   string `form:"centralSettlementMethodNameOrder"  search:"type:order;column:central_settlement_method_name;table:t_order_reconciliation"`
	OrderCreationTime             string `form:"orderCreationTimeOrder"  search:"type:order;column:order_creation_time;table:t_order_reconciliation"`
	TransactionTime               string `form:"transactionTimeOrder"  search:"type:order;column:transaction_time;table:t_order_reconciliation"`
	OperatorName                  string `form:"operatorNameOrder"  search:"type:order;column:operator_name;table:t_order_reconciliation"`
}

func (m *OrderReconciliationGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderReconciliationInsertReq struct {
	Id                            int       `json:"-" comment:"主键ID"` // 主键ID
	ProductName                   string    `json:"productName" comment:"商品名称"`
	ProductCategoryId             string    `json:"productCategoryId" comment:"商品分类ID"`
	ProductCategoryName           string    `json:"productCategoryName" comment:"商品分类名称"`
	PaymentTypeId                 string    `json:"paymentTypeId" comment:"支付类型ID"`
	PaymentTypeName               string    `json:"paymentTypeName" comment:"支付类型名称"`
	TransactionSceneId            string    `json:"transactionSceneId" comment:"交易场景ID"`
	TransactionSceneName          string    `json:"transactionSceneName" comment:"交易场景名称"`
	OrderNumber                   string    `json:"orderNumber" comment:"订单号"`
	TransactionSerialNumber       string    `json:"transactionSerialNumber" comment:"交易流水号"`
	CentralSettlementNumber       string    `json:"centralSettlementNumber" comment:"中联结算号"`
	BankOrderNumber               string    `json:"bankOrderNumber" comment:"银行订单号"`
	UserId                        string    `json:"userId" comment:"用户ID"`
	UserName                      string    `json:"userName" comment:"用户名称"`
	Openid                        string    `json:"openid" comment:"OpenID"`
	UserAccount                   string    `json:"userAccount" comment:"用户账户"`
	SystemOrderAmount             string    `json:"systemOrderAmount" comment:"系统订单金额"`
	MerchantOrderAmount           string    `json:"merchantOrderAmount" comment:"商家订单金额"`
	CentralAmount                 string    `json:"centralAmount" comment:"中联金额"`
	SystemRefundAmount            string    `json:"systemRefundAmount" comment:"系统退款金额"`
	MerchantRefundAmount          string    `json:"merchantRefundAmount" comment:"商家退款金额"`
	MerchantActualRefundAmount    string    `json:"merchantActualRefundAmount" comment:"商家实退金额"`
	MerchantTransactionStatusId   string    `json:"merchantTransactionStatusId" comment:"商家交易状态ID"`
	MerchantTransactionStatusName string    `json:"merchantTransactionStatusName" comment:"商家交易状态名称"`
	SystemTransactionStatusId     string    `json:"systemTransactionStatusId" comment:"系统交易状态ID"`
	SystemTransactionStatusName   string    `json:"systemTransactionStatusName" comment:"系统交易状态名称"`
	CentralTransactionStatusId    string    `json:"centralTransactionStatusId" comment:"中联交易状态ID"`
	CentralTransactionStatusName  string    `json:"centralTransactionStatusName" comment:"中联交易状态名称"`
	CentralSettlementMethodId     string    `json:"centralSettlementMethodId" comment:"中联结算方式ID"`
	CentralSettlementMethodName   string    `json:"centralSettlementMethodName" comment:"中联结算方式名称"`
	OrderCreationTime             time.Time `json:"orderCreationTime" comment:"订单创建时间"`
	TransactionTime               time.Time `json:"transactionTime" comment:"交易时间"`
	OperatorName                  string    `json:"operatorName" comment:"操作员姓名"`
	common.ControlBy
}

func (s *OrderReconciliationInsertReq) Generate(model *models.OrderReconciliation) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ProductName = s.ProductName
	model.ProductCategoryId = s.ProductCategoryId
	model.ProductCategoryName = s.ProductCategoryName
	model.PaymentTypeId = s.PaymentTypeId
	model.PaymentTypeName = s.PaymentTypeName
	model.TransactionSceneId = s.TransactionSceneId
	model.TransactionSceneName = s.TransactionSceneName
	model.OrderNumber = s.OrderNumber
	model.TransactionSerialNumber = s.TransactionSerialNumber
	model.CentralSettlementNumber = s.CentralSettlementNumber
	model.BankOrderNumber = s.BankOrderNumber
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.Openid = s.Openid
	model.UserAccount = s.UserAccount
	model.SystemOrderAmount = s.SystemOrderAmount
	model.MerchantOrderAmount = s.MerchantOrderAmount
	model.CentralAmount = s.CentralAmount
	model.SystemRefundAmount = s.SystemRefundAmount
	model.MerchantRefundAmount = s.MerchantRefundAmount
	model.MerchantActualRefundAmount = s.MerchantActualRefundAmount
	model.MerchantTransactionStatusId = s.MerchantTransactionStatusId
	model.MerchantTransactionStatusName = s.MerchantTransactionStatusName
	model.SystemTransactionStatusId = s.SystemTransactionStatusId
	model.SystemTransactionStatusName = s.SystemTransactionStatusName
	model.CentralTransactionStatusId = s.CentralTransactionStatusId
	model.CentralTransactionStatusName = s.CentralTransactionStatusName
	model.CentralSettlementMethodId = s.CentralSettlementMethodId
	model.CentralSettlementMethodName = s.CentralSettlementMethodName
	model.OrderCreationTime = s.OrderCreationTime
	model.TransactionTime = s.TransactionTime
	model.OperatorName = s.OperatorName
}

func (s *OrderReconciliationInsertReq) GetId() interface{} {
	return s.Id
}

type OrderReconciliationUpdateReq struct {
	Id                            int       `uri:"id" comment:"主键ID"` // 主键ID
	ProductName                   string    `json:"productName" comment:"商品名称"`
	ProductCategoryId             string    `json:"productCategoryId" comment:"商品分类ID"`
	ProductCategoryName           string    `json:"productCategoryName" comment:"商品分类名称"`
	PaymentTypeId                 string    `json:"paymentTypeId" comment:"支付类型ID"`
	PaymentTypeName               string    `json:"paymentTypeName" comment:"支付类型名称"`
	TransactionSceneId            string    `json:"transactionSceneId" comment:"交易场景ID"`
	TransactionSceneName          string    `json:"transactionSceneName" comment:"交易场景名称"`
	OrderNumber                   string    `json:"orderNumber" comment:"订单号"`
	TransactionSerialNumber       string    `json:"transactionSerialNumber" comment:"交易流水号"`
	CentralSettlementNumber       string    `json:"centralSettlementNumber" comment:"中联结算号"`
	BankOrderNumber               string    `json:"bankOrderNumber" comment:"银行订单号"`
	UserId                        string    `json:"userId" comment:"用户ID"`
	UserName                      string    `json:"userName" comment:"用户名称"`
	Openid                        string    `json:"openid" comment:"OpenID"`
	UserAccount                   string    `json:"userAccount" comment:"用户账户"`
	SystemOrderAmount             string    `json:"systemOrderAmount" comment:"系统订单金额"`
	MerchantOrderAmount           string    `json:"merchantOrderAmount" comment:"商家订单金额"`
	CentralAmount                 string    `json:"centralAmount" comment:"中联金额"`
	SystemRefundAmount            string    `json:"systemRefundAmount" comment:"系统退款金额"`
	MerchantRefundAmount          string    `json:"merchantRefundAmount" comment:"商家退款金额"`
	MerchantActualRefundAmount    string    `json:"merchantActualRefundAmount" comment:"商家实退金额"`
	MerchantTransactionStatusId   string    `json:"merchantTransactionStatusId" comment:"商家交易状态ID"`
	MerchantTransactionStatusName string    `json:"merchantTransactionStatusName" comment:"商家交易状态名称"`
	SystemTransactionStatusId     string    `json:"systemTransactionStatusId" comment:"系统交易状态ID"`
	SystemTransactionStatusName   string    `json:"systemTransactionStatusName" comment:"系统交易状态名称"`
	CentralTransactionStatusId    string    `json:"centralTransactionStatusId" comment:"中联交易状态ID"`
	CentralTransactionStatusName  string    `json:"centralTransactionStatusName" comment:"中联交易状态名称"`
	CentralSettlementMethodId     string    `json:"centralSettlementMethodId" comment:"中联结算方式ID"`
	CentralSettlementMethodName   string    `json:"centralSettlementMethodName" comment:"中联结算方式名称"`
	OrderCreationTime             time.Time `json:"orderCreationTime" comment:"订单创建时间"`
	TransactionTime               time.Time `json:"transactionTime" comment:"交易时间"`
	OperatorName                  string    `json:"operatorName" comment:"操作员姓名"`
	common.ControlBy
}

func (s *OrderReconciliationUpdateReq) Generate(model *models.OrderReconciliation) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ProductName = s.ProductName
	model.ProductCategoryId = s.ProductCategoryId
	model.ProductCategoryName = s.ProductCategoryName
	model.PaymentTypeId = s.PaymentTypeId
	model.PaymentTypeName = s.PaymentTypeName
	model.TransactionSceneId = s.TransactionSceneId
	model.TransactionSceneName = s.TransactionSceneName
	model.OrderNumber = s.OrderNumber
	model.TransactionSerialNumber = s.TransactionSerialNumber
	model.CentralSettlementNumber = s.CentralSettlementNumber
	model.BankOrderNumber = s.BankOrderNumber
	model.UserId = s.UserId
	model.UserName = s.UserName
	model.Openid = s.Openid
	model.UserAccount = s.UserAccount
	model.SystemOrderAmount = s.SystemOrderAmount
	model.MerchantOrderAmount = s.MerchantOrderAmount
	model.CentralAmount = s.CentralAmount
	model.SystemRefundAmount = s.SystemRefundAmount
	model.MerchantRefundAmount = s.MerchantRefundAmount
	model.MerchantActualRefundAmount = s.MerchantActualRefundAmount
	model.MerchantTransactionStatusId = s.MerchantTransactionStatusId
	model.MerchantTransactionStatusName = s.MerchantTransactionStatusName
	model.SystemTransactionStatusId = s.SystemTransactionStatusId
	model.SystemTransactionStatusName = s.SystemTransactionStatusName
	model.CentralTransactionStatusId = s.CentralTransactionStatusId
	model.CentralTransactionStatusName = s.CentralTransactionStatusName
	model.CentralSettlementMethodId = s.CentralSettlementMethodId
	model.CentralSettlementMethodName = s.CentralSettlementMethodName
	model.OrderCreationTime = s.OrderCreationTime
	model.TransactionTime = s.TransactionTime
	model.OperatorName = s.OperatorName
}

func (s *OrderReconciliationUpdateReq) GetId() interface{} {
	return s.Id
}

// OrderReconciliationGetReq 功能获取请求参数
type OrderReconciliationGetReq struct {
	Id int `uri:"id"`
}

func (s *OrderReconciliationGetReq) GetId() interface{} {
	return s.Id
}

// OrderReconciliationDeleteReq 功能删除请求参数
type OrderReconciliationDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OrderReconciliationDeleteReq) GetId() interface{} {
	return s.Ids
}
