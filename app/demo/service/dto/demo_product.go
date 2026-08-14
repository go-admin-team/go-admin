package dto

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/demo/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

// DemoProductSearch 列表查询条件
//
// search tag 决定 MakeCondition 拼出的 WHERE：
//
//	exact 精确匹配 / icontains 忽略大小写模糊 / gte 大于等于 …
//
// 未打 search tag 的字段不参与查询，可避免无意间开放过滤维度。
type DemoProductSearch struct {
	dto.Pagination `search:"-"`

	Name   string `form:"name" search:"type:icontains;column:name;table:demo_product"`
	Code   string `form:"code" search:"type:exact;column:code;table:demo_product"`
	Status string `form:"status" search:"type:exact;column:status;table:demo_product"`

	DemoProductOrder
}

// DemoProductOrder 排序字段单独成组，避免与查询字段混在一起
type DemoProductOrder struct {
	CreatedAtOrder string `form:"createdAtOrder" search:"type:order;column:created_at;table:demo_product"`
}

func (m *DemoProductSearch) GetNeedSearch() interface{} { return *m }

func (m *DemoProductSearch) Bind(ctx *gin.Context) error {
	return ctx.ShouldBind(m)
}

func (m *DemoProductSearch) Generate() dto.Index {
	o := *m
	return &o
}

// DemoProductControl 新增与修改共用的入参
//
// 通用 Action（Create / Update）通过 GenerateM 拿到落库对象，
// 因此这里不直接暴露 Model，字段校验用 validate tag 声明。
type DemoProductControl struct {
	Id     int     `json:"id" comment:"主键"`
	Name   string  `json:"name" comment:"名称" validate:"required"`
	Code   string  `json:"code" comment:"编码" validate:"required"`
	Price  float64 `json:"price" comment:"单价" validate:"gte=0"`
	Status string  `json:"status" comment:"状态"`
	Remark string  `json:"remark" comment:"备注"`
}

func (s *DemoProductControl) Bind(ctx *gin.Context) error {
	return ctx.ShouldBind(s)
}

func (s *DemoProductControl) Generate() dto.Control {
	o := *s
	return &o
}

func (s *DemoProductControl) GetId() interface{} { return s.Id }

// GenerateM 组装落库对象。CreateBy / UpdateBy 由通用 Action 在此之后注入，
// 此处不要手动赋值。
func (s *DemoProductControl) GenerateM() (common.ActiveRecord, error) {
	return &models.DemoProduct{
		Model:  common.Model{Id: s.Id},
		Name:   s.Name,
		Code:   s.Code,
		Price:  s.Price,
		Status: s.Status,
		Remark: s.Remark,
	}, nil
}

// DemoProductById 详情与删除共用，支持单个 id 与批量 ids
type DemoProductById struct {
	dto.ObjectById
}

// Bind 与 GetId 由内嵌的 dto.ObjectById 提供：它已处理好 uri 绑定、
// DELETE 时的批量 ids 合并与参数校验，无需在此重复实现。

func (s *DemoProductById) Generate() dto.Control {
	o := *s
	return &o
}

func (s *DemoProductById) GenerateM() (common.ActiveRecord, error) {
	return &models.DemoProduct{}, nil
}
