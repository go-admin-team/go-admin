package dto

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"

	"go-admin/common/dto"
)

// SysPostSearch 列表或者搜索使用结构体
type SysPostSearch struct {
	dto.Pagination `search:"-"`
	PostId         int    `form:"postId" search:"type:exact;column:post_id;table:sys_post" comment:"id"`        // id
	PostName       string `form:"postName" search:"type:contains;column:post_name;table:sys_post" comment:"名称"` // 名称
	PostCode       string `form:"postCode" search:"type:contains;column:post_code;table:sys_post" comment:"编码"` // 编码
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_post" comment:"排序"`             // 排序
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_post" comment:"状态"`         // 状态
	Remark         string `form:"remark" search:"type:exact;column:remark;table:sys_post" comment:"备注"`         // 备注
}

func (m *SysPostSearch) GetNeedSearch() interface{} {
	return *m
}

// SysPostControl 增、改使用的结构体
type SysPostControl struct {
	PostId   int    `uri:"id"  comment:"id"`        // id
	PostName string `form:"postName"  comment:"名称"` // 名称
	PostCode string `form:"postCode" comment:"编码"`  // 编码
	Sort     int    `form:"sort" comment:"排序"`      // 排序
	Status   int    `form:"status"   comment:"状态"`  // 状态
	Remark   string `form:"remark"   comment:"备注"`  // 备注
	common.ControlBy
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysPostControl) Generate(model *models.SysPost) {
	if s.PostId != 0 {
		model.PostId = s.PostId
	}
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
func (s *SysPostControl) GetId() interface{} {
	return s.PostId
}

// SysPostById 获取单个或者删除的结构体
type SysPostById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
	common.ControlBy
}

func (s *SysPostById) Generate(model *models.SysPost) {
	if s.ControlBy.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
	if s.ControlBy.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
}

func (s *SysPostById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysPostById) GenerateM() (*models.SysPost, error) {
	return &models.SysPost{}, nil
}