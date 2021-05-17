package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysApiSearch struct {
	dto.Pagination `search:"-"`
	Name           string `form:"name"  search:"type:exact;column:name;table:sys_api" comment:"名称"`
	Title          string `form:"title"  search:"type:exact;column:title;table:sys_api" comment:"标题"`
	Path           string `form:"path"  search:"type:exact;column:path;table:sys_api" comment:"地址"`
	Action         string `form:"action"  search:"type:exact;column:action;table:sys_api" comment:"类型"`
	ParentId       string `form:"parentId"  search:"type:exact;column:parent_id;table:sys_api" comment:"按钮id"`
	SysApiOrder
}

type SysApiOrder struct {
	HandleOrder string `search:"type:order;column:handle;table:sys_api" form:"handle_order"`
}

func (m *SysApiSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysApiSearch) Generate() dto.Index {
	o := *m
	return &o
}

func (m *SysApiSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Errorf("ShouldBind error: %s", err.Error())
	}
	return err
}

type SysApiControl struct {
	Id       int    `uri:"id" comment:"编码"` // 编码
	Handle   string `json:"handle" comment:"handle"`
	Title    string `json:"title" comment:"标题"`
	Path     string `json:"path" comment:"地址"`
	Paths    string `json:"paths" comment:""`
	Action   string `json:"action" comment:"类型"`
	ParentId int    `json:"parentId" comment:"按钮id"`
	Sort     int    `json:"sort" comment:"排序"`
}

func (s *SysApiControl) Generate() dto.Control {
	o := *s
	return &o
}

func (s *SysApiControl) Bind(ctx *gin.Context) error {
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
	return err
}

func (s *SysApiControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysApi{
		Id:       s.Id,
		Handle:   s.Handle,
		Title:    s.Title,
		Path:     s.Path,
		Paths:    s.Paths,
		Action:   s.Action,
		ParentId: s.ParentId,
		Sort:     s.Sort,
	}, nil
}

func (s *SysApiControl) GetId() interface{} {
	return s.Id
}

type SysApiById struct {
	dto.ObjectById
}

func (s *SysApiById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysApiById) Bind(ctx *gin.Context) error {
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
	return err
}

func (s *SysApiById) SetUpdateBy(id int) {

}

func (s *SysApiById) Generate() dto.Control {
	o := *s
	return &o
}

func (s *SysApiById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysApi{}, nil
}
