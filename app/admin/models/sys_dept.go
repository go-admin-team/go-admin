package models

import "go-admin/common/models"

type SysDept struct {
	DeptId   int    `json:"deptId" gorm:"primaryKey;autoIncrement;"` //部门编码
	ParentId int    `json:"parentId" gorm:""`                        //上级部门
	DeptPath string `json:"deptPath" gorm:"size:255;"`               //
	DeptName string `json:"deptName"  gorm:"size:128;"`              //部门名称
	Sort     int    `json:"sort" gorm:"size:4;"`                            //排序
	Leader   string `json:"leader" gorm:"size:128;"`                 //负责人
	Phone    string `json:"phone" gorm:"size:11;"`                   //手机
	Email    string `json:"email" gorm:"size:64;"`                   //邮箱
	Status   int    `json:"status" gorm:"size:4;"`                   //状态
	models.ControlBy
	models.ModelTime
	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params" gorm:"-"`
	Children  []SysDept `json:"children" gorm:"-"`
}

func (SysDept) TableName() string {
	return "sys_dept"
}

func (e *SysDept) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDept) GetId() interface{} {
	return e.DeptId
}
