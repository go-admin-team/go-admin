package models

//casbin_rule
type CasbinRule struct {
	PType string `json:"p_type" gorm:"type:varchar(100);"`
	V0    string `json:"v0" gorm:"type:varchar(100);"`
	V1    string `json:"v1" gorm:"type:varchar(100);"`
	V2    string `json:"v2" gorm:"type:varchar(100);"`
	V3    string `json:"v3" gorm:"type:varchar(100);"`
	V4    string `json:"v4" gorm:"type:varchar(100);"`
	V5    string `json:"v5" gorm:"type:varchar(100);"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
