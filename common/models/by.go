package models

import contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"

// ControlBy, Model and ModelTime are thin aliases of go-admin-core's
// sdk/contract/models (PRD 006 F1/F5). A type alias is the same type, not a
// new one, so every model that embeds these keeps its GORM tags, JSON tags
// and method set untouched.
type (
	ControlBy = contractmodels.ControlBy
	Model     = contractmodels.Model
	ModelTime = contractmodels.ModelTime
)
