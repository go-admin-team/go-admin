package dto

import (
	"errors"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysCategorySearch struct {
	dto.Pagination `search:"-"`
	Name           string `form:"name" search:"type:exact;column:name;table:sys_category" comment:"名称" vd:"?"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_category" comment:"状态" vd:"?"`
	CateId         int    `form:"cateId" search:"type:exact;column:cate_id;table:sys_category" comment:"分类id" vd:"?"`
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

	if err = vd.Validate(m); err != nil {
		log.Errorf("Validate error: %s", err.Error())
		return err
	}
	return err
}

func (m *SysCategorySearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysCategoryControl struct {
	ID     int    `uri:"Id" comment:"标识"`
	Name   string `json:"name" comment:"名称" vd:"len($)>0 && $!=' '; msg:'invalid name: 不能是空字符串'"`
	Img    string `json:"img" comment:"图标" vd:"?"`
	Sort   int    `json:"sort" comment:"排序" vd:"?"`
	Status int    `json:"status" comment:"状态" vd:"$>0; msg:'invalid status: 状态无效'"`
	Remark string `json:"remark" comment:"备注" vd:"?"`
}

func (s *SysCategoryControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Errorf("ShouldBindUri error: %s", err.Error())
		return errors.New("数据绑定出错")
	}
	err = ctx.ShouldBind(s)
	if err != nil {
		log.Errorf("ShouldBind error: %s", err.Error())
		err = errors.New("数据绑定出错")
	}
	if err1 := vd.Validate(s); err != nil {
		log.Errorf("Validate error: %s", err1.Error())
		return err1
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
	dto.ObjectById `vd:"?"`
}

func (s *SysCategoryById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysCategoryById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysCategory{}, nil
}