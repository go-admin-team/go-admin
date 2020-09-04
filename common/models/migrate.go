package models

import "time"

type Migration struct {
	Version   string    `gorm:"primary_key"`
	ApplyTime time.Time `gorm:"autoCreateTime"`
}

func (Migration) TableName() string {
	return "sys_migration"
}
