package models

import (
	orm "go-admin/database"
	"go-admin/pkg"
	"log"
	"strings"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser() (user SysUser, role SysRole, e error) {

	e = orm.Eloquent.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
	if e != nil {
		if strings.Contains(e.Error(), "record not found") {
			pkg.AssertErr(e, "账号或密码错误(代码204)", 500)
		}
		log.Print(e)
		return
	}
	_, e = pkg.CompareHashAndPassword(user.Password, u.Password)
	if e != nil {
		if strings.Contains(e.Error(), "hashedPassword is not the hash of the given password") {
			pkg.AssertErr(e, "账号或密码错误(代码201)", 500)
		}
		log.Print(e)
		return
	}
	e = orm.Eloquent.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		log.Print(e.Error())
		return
	}
	return
}
