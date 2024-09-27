package models

import (
	"encoding/json"
	"go-admin/common/models"
	"time"
)

type LogNotify struct {
	models.Model

	Data    string                 `json:"data" gorm:"type:longtext;comment:解析后的数据"`
	Parm    string                 `json:"parm" gorm:"type:longtext;comment:请求参数"`
	ParmMap map[string]interface{} `json:"parm_map" gorm:"-"`
	DataMap map[string]interface{} `json:"data_map" gorm:"-"`
	Time    time.Time              `json:"time" gorm:"type:datetime;comment:订单通知时间"`
	Type    string                 `json:"type" gorm:"type:varchar(255);comment:回调名称"`
}

type LogZfbNotify struct {
	LogNotify
}

func (*LogZfbNotify) TableName() string {
	return "m_log_zfb_notify"
}

func (ln *LogZfbNotify) ParseParm() (m map[string]interface{}, e error) {
	var parmMap map[string]interface{}
	if err := json.Unmarshal([]byte(ln.Parm), &parmMap); err != nil {
		return nil, err
	}
	return parmMap, nil
}

type LogWxNotify struct {
	LogNotify
}

func (*LogWxNotify) TableName() string {
	return "m_log_wx_notify"
}

func (ln *LogWxNotify) ParseParm() (m map[string]interface{}, e error) {
	var parmMap map[string]interface{}
	if err := json.Unmarshal([]byte(ln.Parm), &parmMap); err != nil {
		return nil, err
	}
	return parmMap, nil
}

type LogWxScanNotify struct {
	LogNotify
	Result           string                 `json:"result" gorm:"type:varchar(255);comment:微信支付扫码返回结果"`
	ResultMap        map[string]interface{} `json:"result_map" gorm:"-"`
	ResultSuccessMap map[string]interface{} `json:"result_success_map" gorm:"-"`
	DataMap          map[string]interface{} `json:"data_map" gorm:"-"`
}

func (*LogWxScanNotify) TableName() string {
	return "m_log_wx_scan"
}

func (ln *LogWxScanNotify) ParseForResult() (m map[string]interface{}, e error) {
	var parmMap map[string]interface{}
	if err := json.Unmarshal([]byte(ln.Result), &parmMap); err != nil {
		return nil, err
	}
	return parmMap, nil
}

func (ln *LogWxScanNotify) ParseParm() (m map[string]interface{}, e error) {
	var parmMap map[string]interface{}
	if err := json.Unmarshal([]byte(ln.Parm), &parmMap); err != nil {
		return nil, err
	}
	return parmMap, nil
}

func (ln *LogWxScanNotify) ParseForData() (m map[string]interface{}, e error) {
	var parmMap map[string]interface{}
	if err := json.Unmarshal([]byte(ln.Data), &parmMap); err != nil {
		return nil, err
	}
	return parmMap, nil
}
