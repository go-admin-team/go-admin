package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
)

// SysConfigSearch 列表或者搜索使用结构体
type SysDeptSearch struct {
	dto.Pagination `search:"-"`
	DeptId         int    `form:"deptId" search:"type:exact;column:dept_id;table:sys_dept" comment:"id"`       //id
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"上级部门"` //上级部门
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`     //路径
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"部门名称"` //部门名称
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"排序"`            //排序
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"负责人"`       //负责人
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"手机"`          //手机
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"邮箱"`          //邮箱
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"状态"`        //状态
}

func (m *SysDeptSearch) GetNeedSearch() interface{} {
	return *m
}

// Bind 映射上下文中的结构体数据
func (m *SysDeptSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

// SysConfigControl 增、改使用的结构体
type SysDeptControl struct {
	DeptId   int    `uri:"id" comment:"编码"`          // 编码
	ParentId int    `form:"parentId" comment:"上级部门"` //上级部门
	DeptPath string `form:"deptPath" comment:""`     //路径
	DeptName string `form:"deptName" comment:"部门名称"` //部门名称
	Sort     int    `form:"sort" comment:"排序"`       //排序
	Leader   string `form:"leader" comment:"负责人"`    //负责人
	Phone    string `form:"phone" comment:"手机"`      //手机
	Email    string `form:"email" comment:"邮箱"`      //邮箱
	Status   string `form:"status" comment:"状态"`     //状态
}

// Bind 映射上下文中的结构体数据
func (s *SysDeptControl) Bind(ctx *gin.Context) error {
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
func (s *SysDeptControl) Generate() (*system.SysDept, error) {
	return &system.SysDept{
		DeptId:   s.DeptId,
		DeptName: s.DeptName,
		ParentId: s.ParentId,
		DeptPath: s.DeptPath,
		Sort:     s.Sort,
		Leader:   s.Leader,
		Phone:    s.Phone,
		Email:    s.Email,
		Status:   s.Status,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysDeptControl) GetId() interface{} {
	return s.DeptId
}

// SysConfigById 获取单个或者删除的结构体
type SysDeptById struct {
	Id  int   `uri:"id"`
	Ids []int `json:"ids"`
}

func (s *SysDeptById) Generate() *SysDeptById {
	cp := *s
	return &cp
}

func (s *SysDeptById) GetId() interface{} {
	return s.Id
}

func (s *SysDeptById) Bind(ctx *gin.Context) error {
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

func (s *SysDeptById) GenerateM() (*system.SysDept, error) {
	return &system.SysDept{}, nil
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}
