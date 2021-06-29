package dto

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"

	"go-admin/common/dto"
)

// SysPostPageReq 列表或者搜索使用结构体
type SysPostPageReq struct {
	dto.Pagination `search:"-"`
	PostId         int    `form:"postId" search:"type:exact;column:post_id;table:sys_post" comment:"id"`        // id
	PostName       string `form:"postName" search:"type:contains;column:post_name;table:sys_post" comment:"名称"` // 名称
	PostCode       string `form:"postCode" search:"type:contains;column:post_code;table:sys_post" comment:"编码"` // 编码
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_post" comment:"排序"`             // 排序
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_post" comment:"状态"`         // 状态
	Remark         string `form:"remark" search:"type:exact;column:remark;table:sys_post" comment:"备注"`         // 备注
}

func (m *SysPostPageReq) GetNeedSearch() interface{} {
	return *m
}

// SysPostInsertReq 增使用的结构体
type SysPostInsertReq struct {
	PostId   int    `uri:"id"  comment:"id"`
	PostName string `form:"postName"  comment:"名称"`
	PostCode string `form:"postCode" comment:"编码"`
	Sort     int    `form:"sort" comment:"排序"`
	Status   int    `form:"status"   comment:"状态"`
	Remark   string `form:"remark"   comment:"备注"`
	common.ControlBy
}

func (s *SysPostInsertReq) Generate(model *models.SysPost) {
	model.PostName = s.PostName
	model.PostCode = s.PostCode
	model.Sort = s.Sort
	model.Status = s.Status
	model.Remark = s.Remark
	if s.ControlBy.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.ControlBy.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
}

// GetId 获取数据对应的ID
func (s *SysPostInsertReq) GetId() interface{} {
	return s.PostId
}

// SysPostUpdateReq 改使用的结构体
type SysPostUpdateReq struct {
	PostId   int    `uri:"id"  comment:"id"`
	PostName string `form:"postName"  comment:"名称"`
	PostCode string `form:"postCode" comment:"编码"`
	Sort     int    `form:"sort" comment:"排序"`
	Status   int    `form:"status"   comment:"状态"`
	Remark   string `form:"remark"   comment:"备注"`
	common.ControlBy
}

func (s *SysPostUpdateReq) Generate(model *models.SysPost) {
	model.PostId = s.PostId
	model.PostName = s.PostName
	model.PostCode = s.PostCode
	model.Sort = s.Sort
	model.Status = s.Status
	model.Remark = s.Remark
	if s.ControlBy.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.ControlBy.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
}

func (s *SysPostUpdateReq) GetId() interface{} {
	return s.PostId
}

// SysPostGetReq 获取单个的结构体
type SysPostGetReq struct {
	Id int `uri:"id"`
}

func (s *SysPostGetReq) GetId() interface{} {
	return s.Id
}

// SysPostDeleteReq 删除的结构体
type SysPostDeleteReq struct {
	Ids []int `json:"ids"`
	common.ControlBy
}

func (s *SysPostDeleteReq) Generate(model *models.SysPost) {
	if s.ControlBy.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.ControlBy.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
}

func (s *SysPostDeleteReq) GetId() interface{} {
	return s.Ids
}
