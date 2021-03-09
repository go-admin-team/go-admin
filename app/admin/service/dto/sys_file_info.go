package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysFileInfoSearch struct {
	dto.Pagination `search:"-"`
	ID             int    `form:"Id" search:"type:exact;column:id;table:sys_file_info" comment:"标识"`
	Type           string `form:"type" search:"type:exact;column:type;table:sys_file_info" comment:"类型"`
	Name           string `form:"name" search:"type:exact;column:name;table:sys_file_info" comment:"名称"`
	PId            string `form:"pId" search:"type:exact;column:p_id;table:sys_file_info" comment:"目录"`
	Source         string `form:"source" search:"type:exact;column:source;table:sys_file_info" comment:"来源"`
}

func (m *SysFileInfoSearch) GetNeedSearch() interface{} {
	return *m
}

type SysFileInfoControl struct {
	ID       int    `uri:"id" comment:"标识"` // 标识
	Type     string `json:"type" comment:"类型"`
	Name     string `json:"name" comment:"名称"`
	Size     string `json:"size" comment:"大小"`
	PId      int    `json:"pId" comment:"目录"`
	Source   string `json:"source" comment:"来源"`
	Url      string `json:"url" comment:"地址"`
	FullUrl  string `json:"fullUrl" comment:"全地址"`
	CreateBy int    `json:"-"`
	UpdateBy int    `json:"-"`
}

func (s *SysFileInfoControl) Generate() (*models.SysFileInfo, error) {
	return &models.SysFileInfo{
		Model:     common.Model{Id: s.ID},
		Type:      s.Type,
		Name:      s.Name,
		Size:      s.Size,
		PId:       s.PId,
		Source:    s.Source,
		Url:       s.Url,
		FullUrl:   s.FullUrl,
		ControlBy: common.ControlBy{CreateBy: s.CreateBy, UpdateBy: s.UpdateBy},
	}, nil
}

type SysFileInfoById struct {
	dto.ObjectById
	UpdateBy int `json:"-"`
}
