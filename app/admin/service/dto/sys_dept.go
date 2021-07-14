package dto

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"
)

// SysDeptGetPageReq 列表或者搜索使用结构体
type SysDeptGetPageReq struct {
	DeptId         int    `form:"deptId" search:"type:exact;column:dept_id;table:sys_dept" comment:"id"`       //id
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"上级部门"` //上级部门
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`     //路径
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"部门名称"` //部门名称
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"排序"`            //排序
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"负责人"`       //负责人
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"手机"`          //手机
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"邮箱"`          //邮箱
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"状态"`        //状态
}

func (m *SysDeptGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysDeptInsertReq struct {
	DeptId   int    `uri:"id" comment:"编码"`                                         // 编码
	ParentId int    `json:"parentId" comment:"上级部门" vd:"?"`                         //上级部门
	DeptPath string `json:"deptPath" comment:""`                                    //路径
	DeptName string `json:"deptName" comment:"部门名称" vd:"len($)>0"`                  //部门名称
	Sort     int    `json:"sort" comment:"排序" vd:"?"`                               //排序
	Leader   string `json:"leader" comment:"负责人" vd:"@:len($)>0; msg:'leader不能为空'"` //负责人
	Phone    string `json:"phone" comment:"手机" vd:"?"`                              //手机
	Email    string `json:"email" comment:"邮箱" vd:"?"`                              //邮箱
	Status   int    `json:"status" comment:"状态" vd:"$>0"`                           //状态
	common.ControlBy
}

func (s *SysDeptInsertReq) Generate(model *models.SysDept) {
	if s.DeptId != 0 {
		model.DeptId = s.DeptId
	}
	model.DeptName = s.DeptName
	model.ParentId = s.ParentId
	model.DeptPath = s.DeptPath
	model.Sort = s.Sort
	model.Leader = s.Leader
	model.Phone = s.Phone
	model.Email = s.Email
	model.Status = s.Status
}

// GetId 获取数据对应的ID
func (s *SysDeptInsertReq) GetId() interface{} {
	return s.DeptId
}

type SysDeptUpdateReq struct {
	DeptId   int    `uri:"id" comment:"编码"`                                         // 编码
	ParentId int    `json:"parentId" comment:"上级部门" vd:"?"`                         //上级部门
	DeptPath string `json:"deptPath" comment:""`                                    //路径
	DeptName string `json:"deptName" comment:"部门名称" vd:"len($)>0"`                  //部门名称
	Sort     int    `json:"sort" comment:"排序" vd:"?"`                               //排序
	Leader   string `json:"leader" comment:"负责人" vd:"@:len($)>0; msg:'leader不能为空'"` //负责人
	Phone    string `json:"phone" comment:"手机" vd:"?"`                              //手机
	Email    string `json:"email" comment:"邮箱" vd:"?"`                              //邮箱
	Status   int    `json:"status" comment:"状态" vd:"$>0"`                           //状态
	common.ControlBy
}

// Generate 结构体数据转化 从 SysDeptControl 至 SysDept 对应的模型
func (s *SysDeptUpdateReq) Generate(model *models.SysDept) {
	if s.DeptId != 0 {
		model.DeptId = s.DeptId
	}
	model.DeptName = s.DeptName
	model.ParentId = s.ParentId
	model.DeptPath = s.DeptPath
	model.Sort = s.Sort
	model.Leader = s.Leader
	model.Phone = s.Phone
	model.Email = s.Email
	model.Status = s.Status
}

// GetId 获取数据对应的ID
func (s *SysDeptUpdateReq) GetId() interface{} {
	return s.DeptId
}

type SysDeptGetReq struct {
	Id int `uri:"id"`
}

func (s *SysDeptGetReq) GetId() interface{} {
	return s.Id
}

type SysDeptDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysDeptDeleteReq) GetId() interface{} {
	return s.Ids
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}
