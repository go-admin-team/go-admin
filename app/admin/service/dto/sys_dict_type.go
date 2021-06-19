package dto

import (
	"go-admin/app/admin/models"

	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysDictTypeSearch struct {
	dto.Pagination `search:"-"`
	DictId         []int  `form:"dictId" search:"type:in;column:dict_id;table:sys_dict_type"`
	DictName       string `form:"dictName" search:"type:icontains;column:dict_name;table:sys_dict_type"`
	DictType       string `form:"dictType" search:"type:icontains;column:dict_type;table:sys_dict_type"`
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_dict_type"`
}

type SysDictTypeOrder struct {
	DictIdOrder string `search:"type:order;column:dict_id;table:sys_dict_type" form:"dictIdOrder"`
}

func (m *SysDictTypeSearch) GetNeedSearch() interface{} {
	return *m
}

type SysDictTypeControl struct {
	Id       int    `uri:"id"`
	DictName string `json:"dictName"`
	DictType string `json:"dictType"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
	common.ControlBy
}

func (s *SysDictTypeControl) Generate(model *models.SysDictType) {
	if s.Id != 0 {
		model.ID = s.Id
	}
	model.DictName = s.DictName
	model.DictType = s.DictType
	model.Status = s.Status
	model.Remark = s.Remark

}

func (s *SysDictTypeControl) GetId() interface{} {
	return s.Id
}

type SysDictTypeById struct {
	dto.ObjectById
	common.ControlBy
}

func (s *SysDictTypeById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictTypeById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysDictTypeById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysDictType{}, nil
}
