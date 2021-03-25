package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin/app/admin/models/system"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysUserSearch struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"UserId" search:"type:exact;column:user_id;table:sys_user" comment:"用户ID"`
	Username       string `form:"username" search:"type:contains;column:username;table:sys_user" comment:"用户名"`
	NickName       string `form:"nickName" search:"type:contains;column:nick_name;table:sys_user" comment:"昵称"`
	Phone          string `form:"phone" search:"type:contains;column:phone;table:sys_user" comment:"手机号"`
	RoleId         string `form:"roleId" search:"type:exact;column:role_id;table:sys_user" comment:"角色ID"`
	Sex            string `form:"sex" search:"type:exact;column:sex;table:sys_user" comment:"性别"`
	Email          string `form:"email" search:"type:contains;column:email;table:sys_user" comment:"邮箱"`
	DeptId         string `form:"deptId" search:"type:exact;column:dept_id;table:sys_user" comment:"部门"`
	PostId         string `form:"postId" search:"type:exact;column:post_id;table:sys_user" comment:"岗位"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_user" comment:"状态"`
}

func (m *SysUserSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysUserSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (m *SysUserSearch) Generate() dto.Index {
	o := *m
	return &o
}

type SysUserControl struct {
	UserId   int    `json:"userId" comment:"用户ID"` // 用户ID
	Username string `json:"username" comment:"用户名"`
	Password string `json:"password" comment:"密码"`
	NickName string `json:"nickName" comment:"昵称"`
	Phone    string `json:"phone" comment:"手机号"`
	RoleId   int    `json:"roleId" comment:"角色ID"`
	Avatar   string `json:"avatar" comment:"头像"`
	Sex      string `json:"sex" comment:"性别"`
	Email    string `json:"email" comment:"邮箱"`
	DeptId   int    `json:"deptId" comment:"部门"`
	PostId   int    `json:"postId" comment:"岗位"`
	Remark   string `json:"remark" comment:"备注"`
	Status   string `json:"status" comment:"状态"`
}

func (s *SysUserControl) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	//err := ctx.ShouldBindUri(s)
	//if err != nil {
	//	log.Debugf("ShouldBindUri error: %s", err.Error())
	//	return err
	//}
	err := ctx.ShouldBind(s)
	if err != nil {
		log.Debugf("ShouldBind error: %s", err.Error())
	}
	return err
}

func (s *SysUserControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysUserControl) GenerateM() (common.ActiveRecord, error) {
	return &system.SysUser{
		UserId:   s.UserId,
		Username: s.Username,
		Password: s.Password,
		NickName: s.NickName,
		Phone:    s.Phone,
		RoleId:   s.RoleId,
		Avatar:   s.Avatar,
		Sex:      s.Sex,
		Email:    s.Email,
		DeptId:   s.DeptId,
		PostId:   s.PostId,
		Remark:   s.Remark,
		Status:   s.Status,
	}, nil
}

func (s *SysUserControl) GetId() interface{} {
	return s.UserId
}

type SysUserById struct {
	dto.ObjectById
}

func (s *SysUserById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysUserById) GenerateM() (common.ActiveRecord, error) {
	return &system.SysUser{}, nil
}

// PassWord 密码
type PassWord struct {
	NewPassword string `json:"newPassword" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
}
