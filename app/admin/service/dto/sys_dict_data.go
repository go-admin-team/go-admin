package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysDictDataSearch struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id" search:"type:exact;column:dict_code;table:sys_dict_data" comment:""`
	DictLabel      string `form:"dictLabel" search:"type:contains;column:dict_label;table:sys_dict_data" comment:""`
	DictValue      string `form:"dictValue" search:"type:contains;column:dict_value;table:sys_dict_data" comment:""`
	DictType       string `form:"dictType" search:"type:contains;column:dict_type;table:sys_dict_data" comment:""`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dict_data" comment:""`
}

func (m *SysDictDataSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysDictDataSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (m *SysDictDataSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysDictDataControl struct {
	Id        int    `uri:"dictCode" comment:""`
	DictSort  int    `json:"dictSort" comment:""`
	DictLabel string `json:"dictLabel" comment:""`
	DictValue string `json:"dictValue" comment:""`
	DictType  string `json:"dictType" comment:""`
	CssClass  string `json:"cssClass" comment:""`
	ListClass string `json:"listClass" comment:""`
	IsDefault string `json:"isDefault" comment:""`
	Status    string `json:"status" comment:""`
	Default   string `json:"default" comment:""`
	Remark    string `json:"remark" comment:""`
}

func (s *SysDictDataControl) Bind(ctx *gin.Context) error {
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

func (s *SysDictDataControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictDataControl) GenerateM() (common.ActiveRecord, error) {
	return &system.SysDictData{
		DictCode:  s.Id,
		DictSort:  s.DictSort,
		DictLabel: s.DictLabel,
		DictValue: s.DictValue,
		DictType:  s.DictType,
		CssClass:  s.CssClass,
		ListClass: s.ListClass,
		IsDefault: s.IsDefault,
		Status:    s.Status,
		Default:   s.Default,
		Remark:    s.Remark,
	}, nil
}

func (s *SysDictDataControl) GetId() interface{} {
	return s.Id
}

type SysDictDataById struct {
	Id  int   `uri:"dictCode"`
	Ids []int `json:"ids"`
}

func (s *SysDictDataById) Bind(ctx *gin.Context) error {
	if ctx.Request.Method == http.MethodDelete {
		err := ctx.ShouldBind(&s.Ids)
		if err != nil {
			return err
		}
		if len(s.Ids) > 0 {
			return nil
		}
	}
	return ctx.ShouldBindUri(s)
}

func (s *SysDictDataById) GetId() interface{} {
	if len(s.Ids) > 0 {
		return s.Ids
	}
	return s.Id
}

func (s *SysDictDataById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictDataById) GenerateM() (common.ActiveRecord, error) {
	return &system.SysDictData{}, nil
}
