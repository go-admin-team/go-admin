package dto

import (
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

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
	common.ControlBy
}

func (s *SysDictDataControl) Generate(model *models.SysDictData) {

	model.DictCode = s.Id
	model.DictSort = s.DictSort
	model.DictLabel = s.DictLabel
	model.DictValue = s.DictValue
	model.DictType = s.DictType
	model.CssClass = s.CssClass
	model.ListClass = s.ListClass
	model.IsDefault = s.IsDefault
	model.Status = s.Status
	model.Default = s.Default
	model.Remark = s.Remark
}

func (s *SysDictDataControl) GetId() interface{} {
	return s.Id
}

type SysDictDataById struct {
	Id  int   `uri:"dictCode"`
	Ids []int `json:"ids"`
	common.ControlBy
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
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysDictDataById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysDictDataById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysDictData{}, nil
}