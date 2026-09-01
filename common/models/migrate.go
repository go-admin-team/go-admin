package models

import "time"

type Migration struct {
	Version   string    `gorm:"primaryKey"`
	ApplyTime time.Time `gorm:"autoCreateTime"`

	// AppCode identifies which app registered this migration. The empty string
	// means the framework itself.
	//
	// NOT NULL DEFAULT '' rather than a nullable column, and the difference is
	// not cosmetic: on a nullable column the rows that already exist when
	// AutoMigrate adds it hold NULL, and the first SELECT scanning one into
	// this string field fails with "converting NULL to string is unsupported".
	// The default is what makes "existing history belongs to the framework"
	// true without a backfill script anyone could forget to run.
	AppCode string `gorm:"type:varchar(64);not null;default:'';index:idx_sys_migration_app_code;comment:AppCode"`
}

func (Migration) TableName() string {
	return "sys_migration"
}
