package models

import contractmodels "github.com/go-admin-team/go-admin-core/v2/sdk/contract/models"

// Directory, Menu and Button are the menu type enum values used by
// sys_menu.menu_type, referenced directly from go-admin-core's
// sdk/contract/models rather than restated as literals: PRD 006's hard
// constraint 4 requires `const X = pkg.X` for exactly this reason - two
// independently written copies of the same value can be edited out of step,
// where a direct reference cannot.
const (
	Directory = contractmodels.Directory
	Menu      = contractmodels.Menu
	Button    = contractmodels.Button
)
