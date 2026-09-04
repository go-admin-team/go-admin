package dto

import contractdto "github.com/go-admin-team/go-admin-core/v2/sdk/contract/dto"

// AutoForm and the types below describe a form built by go-admin-ui's form
// designer. They are thin aliases of go-admin-core's sdk/contract/dto (PRD
// 006 F2/F5).
type (
	AutoForm = contractdto.AutoForm
	Config   = contractdto.Config
	Option   = contractdto.Option
	Slot     = contractdto.Slot
	Field    = contractdto.Field
	Style    = contractdto.Style
)
