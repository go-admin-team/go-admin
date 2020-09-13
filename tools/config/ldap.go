package config

import "github.com/spf13/viper"

type Ldap struct {
	Enabled       bool
	Server        string
	BindDn        string
	BindPassword  string
	SearchBaseDns string
	SearchFilter  string
	Username      string
	Nickname      string
	Email         string
	RoleId        int
	DeptId        int
	Status        string
}

func InitLdap(cfg *viper.Viper) *Ldap {

	ldap := &Ldap{
		Enabled:       cfg.GetBool("enabled"),
		Server:        cfg.GetString("server"),
		BindDn:        cfg.GetString("bind_dn"),
		BindPassword:  cfg.GetString("bind_password"),
		SearchBaseDns: cfg.GetString("search_base_dns"),
		SearchFilter:  cfg.GetString("search_filter"),
		Username:      cfg.GetString("attribute_username"),
		Nickname:      cfg.GetString("attribute_nickname"),
		Email:         cfg.GetString("attribute_email"),
		RoleId:        cfg.GetInt("attribute_role_id"),
		DeptId:        cfg.GetInt("attribute_dept_id"),
		Status:        cfg.GetString("attribute_status"),
	}
	return ldap
}

var LdapConfig = new(Ldap)
