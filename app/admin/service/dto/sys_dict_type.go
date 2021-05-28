package dto

import (
	"github.com/gin-gonic/gin"

	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/admin/models/system"
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

func (m *SysDictTypeSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysDictTypeSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (m *SysDictTypeSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysDictTypeControl struct {
	Id int `uri:"id" comment:""` //

	DictName string `json:"dictName" comment:""`

	DictType string `json:"dictType" comment:""`

	Status string `json:"status" comment:""`

	Remark string `json:"remark" comment:""`
}

func (s *SysDictTypeControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysDictTypeControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictTypeControl) GenerateM() (common.ActiveRecord, error) {
	return &system.SysDictType{

		ID:       s.Id,
		DictName: s.DictName,
		DictType: s.DictType,
		Status:   s.Status,
		Remark:   s.Remark,
	}, nil
}

func (s *SysDictTypeControl) GetId() interface{} {
	return s.Id
}

type SysDictTypeById struct {
	dto.ObjectById
}

func (s *SysDictTypeById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictTypeById) GenerateM() (common.ActiveRecord, error) {
	return &system.SysDictType{}, nil
}
