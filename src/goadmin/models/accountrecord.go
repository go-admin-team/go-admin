package models

import (
	"github.com/shopspring/decimal"
	orm "goadmin/database"
	_ "log"
	"time"
)

type AccountRecords struct {
	Ids []int `json:"ids"`
}

type Accountrecord struct {
	Id int64 `gorm:"column:id" json:"id"`
	//AccountType     int             `gorm:"column:account_type" json:"account_type"`
	//Type            int             `gorm:"column:type" json:"type"`
	//OmId            int64           `gorm:"column:om_id" json:"om_id"`
	//StId            int64           `gorm:"column:st_id" json:"st_id"`
	//ContractId      int64           `gorm:"column:contract_id" json:"contract_id"`
	//FeeDepartmentId int             `gorm:"column:fee_department_id" json:"fee_department_id"`
	//FeeDepartment   string          `gorm:"column:fee_department" json:"fee_department"`
	//Payer           string          `gorm:"column:payer" json:"payer"`
	//PayerPhone      string          `gorm:"column:payer_phone" json:"payer_phone"`
	//FeeType         int             `gorm:"column:fee_type" json:"fee_type"`
	//OgFeeAmt        decimal.Decimal `gorm:"column:og_fee_amt" json:"og_fee_amt"`
	//FeeAmt          decimal.Decimal `gorm:"column:fee_amt" json:"fee_amt"`
	LeftAmt      decimal.Decimal `gorm:"column:left_amt" json:"left_amt"`
	Status       string          `gorm:"column:status" json:"status"`
	Description  string          `gorm:"column:description" json:"description"`
	Certificate  string          `gorm:"column:certificate" json:"certificate"`
	FeeBeginTime string          `gorm:"column:fee_begin_time" json:"fee_begin_time"`
	FeeEndTime   string          `gorm:"column:fee_end_time" json:"fee_end_time"`
	PayTime      string          `gorm:"column:pay_time" json:"pay_time"`
	ReceiptTime  string          `gorm:"column:receipt_time" json:"receipt_time"`
	PayType      int             `gorm:"column:pay_type" json:"pay_type"`
	BankOrderNo  string          `gorm:"column:bank_order_no" json:"bank_order_no"`
	BillNo       string          `gorm:"column:bill_no" json:"bill_no"`
	Appendix     string          `gorm:"column:appendix" json:"appendix"`
	TransId      int64           `gorm:"column:trans_id" json:"trans_id"`
	Channel      int             `gorm:"column:channel" json:"channel"`
	DelFlag      string          `gorm:"column:del_flag" json:"del_flag"`
	OverdueDays  int             `gorm:"column:overdue_days" json:"overdue_days"`
	ParentId     int64           `gorm:"column:parent_id" json:"parent_id"`
	CreateTime   string          `gorm:"column:create_time" json:"create_time"`
	Creator      string          `gorm:"column:creator" json:"creator"`
	CreatorId    int             `gorm:"column:creator_id" json:"creator_id"`
	UpdateTime   string          `gorm:"column:update_time" json:"update_time"`
	Updater      string          `gorm:"column:updater" json:"updater"`
	UpdaterId    int             `gorm:"column:updater_id" json:"updater_id"`
	Remark       string          `gorm:"column:remark" json:"remark"`
}

func (data *Accountrecord) Get() (dataList Accountrecord, err error) {
	if err = orm.Eloquent.Table("tb_account_record").First(&dataList, data.Id).Error; err != nil {
		return
	}
	return
}

func (data *Accountrecord) GetList(ids []int) (dataList []Accountrecord, err error) {
	if err = orm.Eloquent.Table("tb_account_record").Where("id IN (?)", ids).Find(&dataList).Error; err != nil {
		return
	}
	return
}

func (data *Accountrecord) BatchUpdate(ids []int, maps map[string]interface{}) (err error) {
	data.UpdateTime = time.Now().String()
	if err = orm.Eloquent.Table("tb_account_record").Where("id IN (?)", ids).Updates(maps).Error; err != nil {
		return
	}
	return
}
