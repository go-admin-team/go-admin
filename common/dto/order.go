package dto

import (
	"gorm.io/gorm"

	contractdto "github.com/go-admin-team/go-admin-core/v2/sdk/contract/dto"
)

// OrderDest forwards to go-admin-core's sdk/contract/dto (PRD 006 F2/F5). A
// function cannot be aliased the way a type can, so this is a pure
// pass-through rather than a `func X = pkg.X` form Go does not have.
func OrderDest(sort string, bl bool) func(db *gorm.DB) *gorm.DB {
	return contractdto.OrderDest(sort, bl)
}
