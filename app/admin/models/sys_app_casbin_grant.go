package models

import "time"

// SysAppCasbinGrant is a ledger of casbin_rule rows an app install created,
// keyed by the exact natural key casbin_rule itself is unique on. It exists
// because casbin_rule is not a table this project owns (see design doc
// docs-prd/008-应用清单与安装器/数据库变更.md §2.2): we cannot add an
// app_code column to it without that column being silently zeroed the first
// time anything calls the gorm-adapter's SavePolicy/SavePolicyCtx. Recording
// the natural key here, instead of a foreign key into casbin_rule, is also
// what survives SysRole.Update's RemoveFilteredPolicy+re-add cycle for a
// role's policies (app/admin/service/sys_role.go): that cycle replaces the
// underlying row (a new auto-increment ID) but reproduces the same
// (ptype,v0,v1,v2) tuple from the same sys_menu/sys_api data, so a
// natural-key match here still finds it. What it does not survive is the
// role being renamed, or the tuple being rebuilt from a completely different
// source (a future SavePolicy call from outside this seeder) - in both cases
// the match legitimately fails, and business rule 3 says the uninstaller
// should report and skip, not delete something else that happens to look
// the same.
type SysAppCasbinGrant struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement"`

	AppCode string `json:"appCode" gorm:"type:varchar(64);not null;index:idx_sys_app_casbin_grant_app_code;comment:app code that created this grant"`

	// Column widths mirror gorm-adapter's own CasbinRule struct exactly, so
	// a value that fits into casbin_rule always fits here, and the unique
	// index below matches the one createTable() puts on casbin_rule itself.
	Ptype string `json:"ptype" gorm:"size:100;not null;uniqueIndex:uk_sys_app_casbin_grant_rule;comment:casbin ptype, 'p' today"`
	V0    string `json:"v0" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:role_key at grant time"`
	V1    string `json:"v1" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:api path"`
	V2    string `json:"v2" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:http method"`
	V3    string `json:"v3" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:unused today"`
	V4    string `json:"v4" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:unused today"`
	V5    string `json:"v5" gorm:"size:100;not null;default:'';uniqueIndex:uk_sys_app_casbin_grant_rule;comment:unused today"`

	CreatedAt time.Time `json:"createdAt" gorm:"comment:when this grant was recorded"`
}

func (*SysAppCasbinGrant) TableName() string {
	return "sys_app_casbin_grant"
}
