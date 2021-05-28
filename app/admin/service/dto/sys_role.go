package dto

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
)

// SysRoleSearch 列表或者搜索使用结构体
type SysRoleSearch struct {
	dto.Pagination `search:"-"`

	RoleId    int    `form:"roleId" search:"type:exact;column:role_id;table:sys_role" comment:"角色编码"`     // 角色编码
	RoleName  string `form:"roleName" search:"type:exact;column:role_name;table:sys_role" comment:"角色名称"` // 角色名称
	Status    string `form:"status" search:"type:exact;column:status;table:sys_role" comment:"状态"`        // 状态
	RoleKey   string `form:"roleKey" search:"type:exact;column:role_key;table:sys_role" comment:"角色代码"`   // 角色代码
	RoleSort  int    `form:"roleSort" search:"type:exact;column:role_sort;table:sys_role" comment:"角色排序"` // 角色排序
	Flag      string `form:"flag" search:"type:exact;column:flag;table:sys_role" comment:"标记"`            // 标记
	Remark    string `form:"remark" search:"type:exact;column:remark;table:sys_role" comment:"备注"`        // 备注
	Admin     bool   `form:"admin" search:"type:exact;column:admin;table:sys_role" comment:"是否管理员"`
	DataScope string `form:"dataScope" search:"type:exact;column:data_scope;table:sys_role" comment:"是否管理员"`
}

func (m *SysRoleSearch) GetNeedSearch() interface{} {
	return *m
}

// Bind 映射上下文中的结构体数据
func (m *SysRoleSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

// SysConfigControl 增、改使用的结构体
type SysRoleControl struct {
	RoleId    int    `uri:"id" comment:"角色编码"`        // 角色编码
	RoleName  string `form:"roleName" comment:"角色名称"` // 角色名称
	Status    string `form:"status" comment:"状态"`     // 状态
	RoleKey   string `form:"roleKey" comment:"角色代码"`  // 角色代码
	RoleSort  int    `form:"roleSort" comment:"角色排序"` // 角色排序
	Flag      string `form:"flag" comment:"标记"`       // 标记
	Remark    string `form:"remark" comment:"备注"`     // 备注
	Admin     bool   `form:"admin" comment:"是否管理员"`
	DataScope string `form:"dataScope"`
	MenuIds   []int  `json:"menuIds"`
}

// Bind 映射上下文中的结构体数据
func (s *SysRoleControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBindBodyWith(s, binding.JSON)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	err = ctx.ShouldBindUri(s)
	if err != nil {
		log.Debugf("ShouldBindUri error: %s", err.Error())
		return err
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
func (s *SysRoleControl) Generate() (*system.SysRole, error) {
	return &system.SysRole{
		RoleId:    s.RoleId,
		RoleName:  s.RoleName,
		Status:    s.Status,
		RoleKey:   s.RoleKey,
		RoleSort:  s.RoleSort,
		Flag:      s.Flag,
		Remark:    s.Remark,
		Admin:     s.Admin,
		DataScope: s.DataScope,
		MenuIds:   s.MenuIds,
	}, nil
}

// GetId 获取数据对应的ID
func (s *SysRoleControl) GetId() interface{} {
	return s.RoleId
}

// SysConfigById 获取单个或者删除的结构体
type SysRoleById struct {
	dto.ObjectById
}

func (s *SysRoleById) Generate() *SysRoleById {
	cp := *s
	return &cp
}

func (s *SysRoleById) GetId() interface{} {
	return s.Id
}

func (s *SysRoleById) GenerateM() (*system.SysRole, error) {
	return &system.SysRole{}, nil
}

// RoleDataScopeReq 角色数据权限修改
type RoleDataScopeReq struct {
	RoleId    int    `json:"roleId" binding:"required"`
	DataScope string `json:"dataScope" binding:"required"`
	DeptIds   []int  `json:"deptIds"`
}

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}