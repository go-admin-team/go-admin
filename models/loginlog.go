package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

type LoginLog struct {
	//主键
	InfoId int64 `json:"infoId" gorm:"column:infoId;primary_key"`
	//用户名
	UserName string `json:"userName" gorm:"column:userName;"`
	//状态
	Status string `json:"status" gorm:"column:status;"`
	//ip地址
	Ipaddr string `json:"ipaddr" gorm:"column:ipaddr;"`
	//归属地
	LoginLocation string `json:"loginLocation" gorm:"column:loginLocation;"`
	//浏览器
	Browser string `json:"browser" gorm:"column:browser;"`
	//系统
	Os string `json:"os" gorm:"column:os;"`
	// 固件
	Platform string `json:"platform" gorm:"column:platform;"`
	//登录时间
	LoginTime string `json:"loginTime" gorm:"column:loginTime;"`
	//创建人
	CreateBy string `json:"createBy" gorm:"column:create_by;"`
	//创建时间
	CreateTime string `json:"createTime" gorm:"column:create_time;"`
	//更新者
	UpdateBy string `json:"updateBy" gorm:"column:update_by;"`
	//更新时间
	UpdateTime string `json:"updateTime" gorm:"column:update_time;"`
	//数据
	DataScope string `json:"dataScope" gorm:"column:dataScope;"`
	//参数
	Params string `json:"params" gorm:"column:params;"`
	//备注
	Remark string `json:"remark" gorm:"column:remark;"`
	//是否删除
	IsDel string `json:"isDel" gorm:"column:is_del;"`
	Msg   string `json:"msg" gorm:"column:msg;"`
}

func (e *LoginLog) Get() (LoginLog, error) {
	var doc LoginLog

	table := orm.Eloquent.Table("sys_loginlog")
	if e.Ipaddr != "" {
		table = table.Where("ipaddr = ?", e.Ipaddr)
	}
	if e.InfoId != 0 {
		table = table.Where("infoId = ?", e.InfoId)
	}

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *LoginLog) GetPage(pageSize int, pageIndex int) ([]LoginLog, int32, error) {
	var doc []LoginLog

	table := orm.Eloquent.Select("*").Table("sys_loginlog")
	if e.Ipaddr != "" {
		table = table.Where("ipaddr = ?", e.Ipaddr)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.UserName != "" {
		table = table.Where("userName = ?", e.UserName)
	}

	var count int32

	if err := table.Where("is_del = 0").Order("infoId desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (e *LoginLog) Create() (LoginLog, error) {
	var doc LoginLog
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	e.CreateBy = "0"
	e.UpdateBy = "0"
	result := orm.Eloquent.Table("sys_loginlog").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *LoginLog) Update(id int64) (update LoginLog, err error) {
	e.UpdateTime = utils.GetCurrntTime()

	if err = orm.Eloquent.Table("sys_loginlog").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_loginlog").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *LoginLog) BatchDelete(id []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table("sys_loginlog").Where("is_del=0 and infoId in (?)", id).Update(map[string]interface{}{"is_del": "1", "update_time": utils.GetCurrntTime(), "update_by": e.UpdateBy}).Error; err != nil {
		return
	}
	Result = true
	return
}
