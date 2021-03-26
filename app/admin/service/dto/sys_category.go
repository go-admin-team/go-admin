package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysCategorySearch struct {
	dto.Pagination `search:"-"`
	Name           string `form:"name" search:"type:exact;column:name;table:sys_category" comment:"名称"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_category" comment:"状态"`
	CateId         int    `form:"cateId" search:"type:exact;column:cate_id;table:sys_category" comment:"分类id"`
}

func (m *SysCategorySearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysCategorySearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (m *SysCategorySearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysCategoryControl struct {
	ID     int    `uri:"Id" comment:"标识"`
	Name   string `json:"name" comment:"名称"`
	Img    string `json:"img" comment:"图标"`
	Sort   int    `json:"sort" comment:"排序"`
	Status int    `json:"status" comment:"状态"`
	Remark string `json:"remark" comment:"备注"`
}

func (s *SysCategoryControl) Bind(ctx *gin.Context) error {
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

func (s *SysCategoryControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysCategoryControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysCategory{
		Model:  common.Model{Id: s.ID},
		Name:   s.Name,
		Img:    s.Img,
		Sort:   s.Sort,
		Status: s.Status,
		Remark: s.Remark,
	}, nil
}

func (s *SysCategoryControl) GetId() interface{} {
	return s.ID
}

type SysCategoryById struct {
	dto.ObjectById
}

func (s *SysCategoryById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysCategoryById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysCategory{}, nil
}