package models

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/storage"

	"go-admin/common/models"
)

type SysOperaLog struct {
	models.Model
	Title         string    `json:"title" gorm:"size:255;comment:操作模块"`
	BusinessType  string    `json:"businessType" gorm:"size:128;comment:操作类型"`
	BusinessTypes string    `json:"businessTypes" gorm:"size:128;comment:BusinessTypes"`
	Method        string    `json:"method" gorm:"size:128;comment:函数"`
	RequestMethod string    `json:"requestMethod" gorm:"size:128;comment:请求方式"`
	OperatorType  string    `json:"operatorType" gorm:"size:128;comment:操作类型"`
	OperName      string    `json:"operName" gorm:"size:128;comment:操作者"`
	DeptName      string    `json:"deptName" gorm:"size:128;comment:部门名称"`
	OperUrl       string    `json:"operUrl" gorm:"size:255;comment:访问地址"`
	OperIp        string    `json:"operIp" gorm:"size:128;comment:客户端ip"`
	OperLocation  string    `json:"operLocation" gorm:"size:128;comment:访问位置"`
	OperParam     string    `json:"operParam" gorm:"size:255;comment:请求参数"`
	Status        string    `json:"status" gorm:"size:4;comment:操作状态"`
	OperTime      time.Time `json:"operTime" gorm:"comment:操作时间"`
	JsonResult    string    `json:"jsonResult" gorm:"size:255;comment:返回数据"`
	Remark        string    `json:"remark" gorm:"size:255;comment:备注"`
	LatencyTime   string    `json:"latencyTime" gorm:"size:128;comment:耗时"`
	UserAgent     string    `json:"userAgent" gorm:"size:255;comment:ua"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	models.ControlBy
}

func (SysOperaLog) TableName() string {
	return "sys_opera_log"
}

func (e *SysOperaLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysOperaLog) GetId() interface{} {
	return e.Id
}

// SaveOperaLog 从队列中获取操作日志
func SaveOperaLog(message storage.Messager) (err error) {
	//准备db
	db := sdk.Runtime.GetDbByKey(message.GetPrefix())
	if db == nil {
		err = errors.New("db not exist")
		log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
		// Log writing to the database ignores error
		return nil
	}
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	var l SysOperaLog
	err = json.Unmarshal(rb, &l)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	// 超出100个字符返回值截断
	if len(l.JsonResult) > 100 {
		l.JsonResult = l.JsonResult[:100]
	}
	err = db.Create(&l).Error
	if err != nil {
		log.Errorf("db create error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	return nil
}
