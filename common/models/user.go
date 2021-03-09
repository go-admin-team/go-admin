package models

import (
	"gorm.io/gorm"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
)

// BaseUser 密码登录基础用户
type BaseUser struct {
	Username     string `json:"username" gorm:"type:varchar(100);comment:用户名"`
	Salt         string `json:"-" gorm:"type:varchar(255);comment:加盐;<-"`
	PasswordHash string `json:"-" gorm:"type:varchar(128);comment:密码hash;<-"`
	Password     string `json:"password" gorm:"-"`
}

// SetPassword 设置密码
func (u *BaseUser) SetPassword(value string) {
	u.Password = value
	u.generateSalt()
	u.PasswordHash = u.GetPasswordHash()
}

// GetPasswordHash 获取密码hash
func (u *BaseUser) GetPasswordHash() string {
	passwordHash, err := pkg.SetPassword(u.Password, u.Salt)
	if err != nil {
		return ""
	}
	return passwordHash
}

// generateSalt 生成加盐值
func (u *BaseUser) generateSalt() {
	u.Salt = pkg.GenerateRandomKey16()
}

// Verify 验证密码
func (u *BaseUser) Verify(db *gorm.DB, tableName string) bool {
	db.Table(tableName).Where("username = ?", u.Username).First(u)
	return u.GetPasswordHash() == u.PasswordHash
}
