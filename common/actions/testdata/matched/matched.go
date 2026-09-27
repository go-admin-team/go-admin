// Package matched is mismatch with the pair corrected: the control that
// shows mismatch fails for the pairing, not for how it is built.
package matched

import (
	"go-admin/app/demo/models"
	"go-admin/app/demo/service/dto"
	"go-admin/common/actions"
)

var _ = actions.Create[models.DemoProduct, dto.DemoProductControl]
