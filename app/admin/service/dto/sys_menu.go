package dto

import (
	"go-admin/app/admin/models"
	common "go-admin/common/models"

	"go-admin/common/dto"
)

// SysMenuSearch 列表或者搜索使用结构体
type SysMenuSearch struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title" search:"type:contains;column:title;table:sys_menu" comment:"菜单名称"`  // 菜单名称
	Visible        int    `form:"visible" search:"type:exact;column:visible;table:sys_menu" comment:"显示状态"` // 显示状态
}

func (m *SysMenuSearch) GetNeedSearch() interface{} {
	return *m
}

// SysMenuControl 增、改使用的结构体
type SysMenuControl struct {
	MenuId     int             `uri:"id" comment:"编码"`            // 编码
	MenuName   string          `form:"menuName" comment:"菜单name"` //菜单name
	Title      string          `form:"title" comment:"显示名称"`      //显示名称
	Icon       string          `form:"icon" comment:"图标"`         //图标
	Path       string          `form:"path" comment:"路径"`         //路径
	Paths      string          `form:"paths" comment:"id路径"`      //id路径
	MenuType   string          `form:"menuType" comment:"菜单类型"`   //菜单类型
	SysApi     []models.SysApi `form:"sysApi"`
	Apis       []int           `form:"apis"`
	Action     string          `form:"action" comment:"请求方式"`      //请求方式
	Permission string          `form:"permission" comment:"权限编码"`  //权限编码
	ParentId   int             `form:"parentId" comment:"上级菜单"`    //上级菜单
	NoCache    bool            `form:"noCache" comment:"是否缓存"`     //是否缓存
	Breadcrumb string          `form:"breadcrumb" comment:"是否面包屑"` //是否面包屑
	Component  string          `form:"component" comment:"组件"`     //组件
	Sort       int             `form:"sort" comment:"排序"`          //排序
	Visible    string          `form:"visible" comment:"是否显示"`     //是否显示
	IsFrame    string          `form:"isFrame" comment:"是否frame"`  //是否frame
	common.ControlBy
}

func (s *SysMenuControl) SetCreateBy(id int) {
	s.CreateBy = id
}

func (s *SysMenuControl) SetUpdateBy(id int) {
	s.UpdateBy = id
}

// Generate 结构体数据转化 从 Control 至 model 对应的模型
func (s *SysMenuControl) Generate(model *models.SysMenu) {
	if s.MenuId != 0 {
		model.MenuId = s.MenuId
	}
	model.MenuName = s.MenuName
	model.Title = s.Title
	model.Icon = s.Icon
	model.Path = s.Path
	model.Paths = s.Paths
	model.MenuType = s.MenuType
	model.Action = s.Action
	model.SysApi = s.SysApi
	model.Permission = s.Permission
	model.ParentId = s.ParentId
	model.NoCache = s.NoCache
	model.Breadcrumb = s.Breadcrumb
	model.Component = s.Component
	model.Sort = s.Sort
	model.Visible = s.Visible
	model.IsFrame = s.IsFrame
	if s.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
	if s.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
}

// GetId 获取数据对应的ID
func (s *SysMenuControl) GetId() interface{} {
	return s.MenuId
}

// SysMenuById 获取单个或者删除的结构体
type SysMenuById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
	common.ControlBy
}

func (s *SysMenuById) Generate() *SysMenuById {
	cp := *s
	return &cp
}

func (s *SysMenuById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysMenuById) GenerateM() (*models.SysMenu, error) {
	return &models.SysMenu{}, nil
}

type MenuLabel struct {
	Id       int         `json:"id,omitempty" gorm:"-"`
	Label    string      `json:"label,omitempty" gorm:"-"`
	Children []MenuLabel `json:"children,omitempty" gorm:"-"`
}

type MenuRole struct {
	models.SysMenu
	IsSelect bool `json:"is_select" gorm:"-"`
}

type SelectRole struct {
	RoleId int `uri:"roleId"`
}