package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
)

// SysMenuSearch 列表或者搜索使用结构体
type SysMenuSearch struct {
	dto.Pagination `search:"-"`
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"上级部门"` //上级部门
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`     //路径
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"部门名称"` //部门名称
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"排序"`            //排序
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"负责人"`       //负责人
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"手机"`          //手机
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"邮箱"`          //邮箱
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"状态"`        //状态
}

func (m *SysMenuSearch) GetNeedSearch() interface{} {
	return *m
}

// Bind 映射上下文中的结构体数据
func (m *SysMenuSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

// SysConfigControl 增、改使用的结构体
type SysMenuControl struct {
	MenuId     int    `uri:"id" comment:"编码"`             // 编码
	MenuName   string `form:"menuName" comment:"菜单name"`  //菜单name
	Title      string `form:"title" comment:"显示名称"`       //显示名称
	Icon       string `form:"icon" comment:"图标"`          //图标
	Path       string `form:"path" comment:"路径"`          //路径
	Paths      string `form:"paths" comment:"id路径"`       //id路径
	MenuType   string `form:"menuType" comment:"菜单类型"`    //菜单类型
	Action     string `form:"action" comment:"请求方式"`      //请求方式
	Permission string `form:"permission" comment:"权限编码"`  //权限编码
	ParentId   int    `form:"parentId" comment:"上级菜单"`    //上级菜单
	NoCache    bool   `form:"noCache" comment:"是否缓存"`     //是否缓存
	Breadcrumb string `form:"breadcrumb" comment:"是否面包屑"` //是否面包屑
	Component  string `form:"component" comment:"组件"`     //组件
	Sort       int    `form:"sort" comment:"排序"`          //排序
	Visible    string `form:"visible" comment:"是否显示"`     //是否显示
	IsFrame    string `form:"isFrame" comment:"是否frame"`  //是否frame
}

// Bind 映射上下文中的结构体数据
func (s *SysMenuControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
	}
	err = ctx.ShouldBindBodyWith(s, binding.JSON)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	var jsonStr []byte
	jsonStr, err = json.Marshal(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	ctx.Set("body", string(jsonStr))
	return err
}

// Generate 结构体数据转化 从 SysConfigControl 至 system.SysConfig 对应的模型
func (s *SysMenuControl) Generate() (*system.SysMenu, error) {
	return &system.SysMenu{
		MenuId:     s.MenuId,
		MenuName:   s.MenuName,
		Title:      s.Title,
		Icon:       s.Icon,
		Path:       s.Path,
		Paths:      s.Paths,
		MenuType:   s.MenuType,
		Action:     s.Action,
		Permission: s.Permission,
		ParentId:   s.ParentId,
		NoCache:    s.NoCache,
		Breadcrumb: s.Breadcrumb,
		Component:  s.Component,
		Sort:       s.Sort,
		Visible:    s.Visible,
		IsFrame:    s.IsFrame,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysMenuControl) GetId() interface{} {
	return s.MenuId
}

// SysConfigById 获取单个或者删除的结构体
type SysMenuById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysMenuById) Generate() *SysMenuById {
	cp := *s
	return &cp
}

func (s *SysMenuById) GetId() interface{} {
	return s.Id
}

func (s *SysMenuById) Bind(ctx *gin.Context) error {
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

func (s *SysMenuById) GenerateM() (*system.SysMenu, error) {
	return &system.SysMenu{}, nil
}

type MenuLabel struct {
	Id       int         `json:"id" gorm:"-"`
	Label    string      `json:"label" gorm:"-"`
	Children []MenuLabel `json:"children" gorm:"-"`
}

type MenuRole struct {
	system.SysMenus
	IsSelect bool `json:"is_select" gorm:"-"`
}

type SelectRole struct {
	RoleId int `uri:"roleId"`
}
