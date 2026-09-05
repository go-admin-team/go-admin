package dto

import (
	"gorm.io/gorm"

	contractdto "github.com/go-admin-team/go-admin-core/v2/sdk/contract/dto"
)

// GeneralDelDto and GeneralGetDto are thin aliases of go-admin-core's
// sdk/contract/dto (PRD 006 F2/F5).
type (
	GeneralDelDto = contractdto.GeneralDelDto
	GeneralGetDto = contractdto.GeneralGetDto
)

// MakeCondition and Paginate forward to go-admin-core's sdk/contract/dto
// (PRD 006 F2/F5). This file used to read go-admin/common/global.Driver to
// pick the SQL dialect MakeCondition resolves search tags against; the
// lowered version instead reads db.Dialector.Name() from inside the closure
// it returns, which is always the driver the caller's own *gorm.DB is bound
// to - correct even when a multi-tenant host has more than one database
// open with different drivers, which a single package-level variable could
// never be. global.Driver itself is untouched and still readable, but
// nothing in this package reads it anymore.
func MakeCondition(q interface{}) func(db *gorm.DB) *gorm.DB {
	return contractdto.MakeCondition(q)
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return contractdto.Paginate(pageSize, pageIndex)
}
