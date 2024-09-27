package models

import (

	"go-admin/common/models"

)

type TApiZl struct {
    models.Model
    
    Code string `json:"code" gorm:"type:varchar(50);comment:业务编码"` 
    Handle string `json:"handle" gorm:"type:varchar(128);comment:handle"` 
    Title string `json:"title" gorm:"type:varchar(128);comment:标题"` 
    Path string `json:"path" gorm:"type:varchar(128);comment:地址"` 
    Type string `json:"type" gorm:"type:varchar(16);comment:接口类型"` 
    Action string `json:"action" gorm:"type:varchar(16);comment:请求类型"` 
    Req string `json:"req" gorm:"type:longtext;comment:请求入参"` 
    Res string `json:"res" gorm:"type:longtext;comment:响应参数"` 
    ResError string `json:"resError" gorm:"type:longtext;comment:错误返回"` 
    models.ModelTime
    models.ControlBy
}

func (TApiZl) TableName() string {
    return "t_api_zl"
}

func (e *TApiZl) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *TApiZl) GetId() interface{} {
	return e.Id
}