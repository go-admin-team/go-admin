package models

import (
	orm "go-admin/database"
	"go-admin/utils"
)

//sys_operlog
type SysOperLog struct {

	//日志编码
	OperId int64 `json:"operId" gorm:"column:operId;primary_key"`
	//操作模块
	Title string `json:"title" gorm:"column:title;"`
	//操作类型
	BusinessType  string `json:"businessType" gorm:"column:businessType;"`
	BusinessTypes string `json:"businessTypes" gorm:"column:businessTypes;"`
	//函数
	Method string `json:"method" gorm:"column:method;"`
	//请求方式
	RequestMethod string `json:"requestMethod" gorm:"column:requestMethod;"`
	//操作类型
	OperatorType string `json:"operatorType" gorm:"column:operatorType;"`
	//操作者
	OperName string `json:"operName" gorm:"column:operName;"`
	//部门名称
	DeptName string `json:"deptName" gorm:"column:deptName;"`
	//访问地址
	OperUrl string `json:"operUrl" gorm:"column:operUrl;"`
	//客户端ip
	OperIp string `json:"operIp" gorm:"column:operIp;"`
	//访问位置
	OperLocation string `json:"operLocation" gorm:"column:operLocation;"`
	//请求参数
	OperParam string `json:"operParam" gorm:"column:operParam;"`
	//操作状态
	Status string `json:"status" gorm:"column:status;"`
	//操作时间
	OperTime string `json:"operTime" gorm:"column:operTime;"`
	//返回数据
	JsonResult string `json:"jsonResult" gorm:"column:jsonResult;"`
	//创建人
	CreateBy string `json:"createBy" gorm:"column:create_by;"`
	//创建时间
	CreateTime string `json:"createTime" gorm:"column:create_time;"`
	//更新者
	UpdateBy string `json:"updateBy" gorm:"column:update_by;"`
	//更新时间
	UpdateTime string `json:"updateTime" gorm:"column:update_time;"`
	//数据
	DataScope string `json:"dataScope" gorm:"column:data_scope;"`
	//参数
	Params string `json:"params" gorm:"column:params;"`
	//备注
	Remark string `json:"remark" gorm:"column:remark;"`
	//是否删除
	IsDel string `json:"isDel" gorm:"column:is_del;"`
	//耗时
	LatencyTime string `json:"latencyime" gorm:"column:latency_time;"`
	//us
	UserAgent string `json:"userAgent" gorm:"column:user_agent;"`
}

func (e *SysOperLog) Get() (SysOperLog, error) {
	var doc SysOperLog

	table := orm.Eloquent.Table("sys_operlog")
	if e.OperIp != "" {
		table = table.Where("operIp = ?", e.OperIp)
	}
	if e.OperId != 0 {
		table = table.Where("operId = ?", e.OperId)
	}

	if err := table.Where("is_del = 0").First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysOperLog) GetPage(pageSize int, pageIndex int) ([]SysOperLog, int32, error) {
	var doc []SysOperLog

	table := orm.Eloquent.Select("*").Table("sys_operlog")
	if e.OperIp != "" {
		table = table.Where("operIp = ?", e.OperIp)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.OperName != "" {
		table = table.Where("operName = ?", e.OperName)
	}
	if e.BusinessType != "" {
		table = table.Where("businessType = ?", e.BusinessType)
	}

	var count int32

	if err := table.Where("is_del = 0").Order("operId desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("is_del = 0").Count(&count)
	return doc, count, nil
}

func (e *SysOperLog) Create() (SysOperLog, error) {
	var doc SysOperLog
	e.CreateTime = utils.GetCurrntTime()
	e.UpdateTime = utils.GetCurrntTime()
	e.CreateBy = "0"
	e.UpdateBy = "0"
	result := orm.Eloquent.Table("sys_operlog").Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

func (e *SysOperLog) Update(id int64) (update SysOperLog, err error) {
	e.UpdateTime = utils.GetCurrntTime()

	if err = orm.Eloquent.Table("sys_operlog").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table("sys_operlog").Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysOperLog) BatchDelete(id []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table("sys_operlog").Where("is_del=0 and operId in (?)", id).Update(map[string]interface{}{"is_del": "1", "update_time": utils.GetCurrntTime(), "update_by": e.UpdateBy}).Error; err != nil {
		return
	}
	Result = true
	return
}
