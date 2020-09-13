package models

import (
	"errors"
	"fmt"
	orm "go-admin/global"
	"go-admin/tools"
	"go-admin/tools/config"

	"github.com/go-ldap/ldap"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

type LdapUserInfo struct {
	Username string
	Nickname string
	Email    string
	RoleId   int
	DeptId   int
	Status   string
}

func (u *Login) GetUser() (user SysUser, role SysRole, e error) {

	if config.LdapConfig.Enabled && u.Username != "admin" {
		lui, err := LdapLogin(u)
		if err != nil {
			return
		}
		// check 用户名
		var count int
		orm.Eloquent.Table("sys_user").Where("username = ?", lui.Username).Count(&count)
		if count == 0 {
			//添加用户数据
			var sysuser SysUser
			sysuser.Username = lui.Username
			sysuser.NickName = lui.Nickname
			sysuser.Email = lui.Email
			sysuser.RoleId = config.LdapConfig.RoleId
			sysuser.DeptId = config.LdapConfig.DeptId
			sysuser.Status = config.LdapConfig.Status
			_, err := sysuser.Insert()
			if err != nil {
				return
			}
		}
		e = orm.Eloquent.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
		if e != nil {
			return
		}
	} else {
		e = orm.Eloquent.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
		if e != nil {
			return
		}
		_, e = tools.CompareHashAndPassword(user.Password, u.Password)
		if e != nil {
			return
		}

	}
	e = orm.Eloquent.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		return
	}
	return
}

// LDAP 登录
func LdapLogin(u *Login) (lui LdapUserInfo, e error) {
	conn, err := ldap.DialURL(config.LdapConfig.Server)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Bind(config.LdapConfig.BindDn, config.LdapConfig.BindPassword)
	if err != nil {
		return
	}

	sql := ldap.NewSearchRequest(config.LdapConfig.SearchBaseDns,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		fmt.Sprintf(config.LdapConfig.SearchFilter, u.Username),
		[]string{config.LdapConfig.Username, config.LdapConfig.Nickname, config.LdapConfig.Email},
		nil)
	var cur *ldap.SearchResult
	cur, err = conn.Search(sql)
	if err != nil {
		return
	}
	if len(cur.Entries) != 1 {
		err = errors.New("账户不存在，或存在多个！")
		return
	}

	userDN := cur.Entries[0].DN
	err = conn.Bind(userDN, u.Password)
	if err != nil {
		return
	}

	lui = LdapUserInfo{
		Username: cur.Entries[0].GetAttributeValue(config.LdapConfig.Username),
		Nickname: cur.Entries[0].GetAttributeValue(config.LdapConfig.Nickname),
		Email:    cur.Entries[0].GetAttributeValue(config.LdapConfig.Email),
	}

	return lui, nil
}
