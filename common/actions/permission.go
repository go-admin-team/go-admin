package actions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	contractactions "github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"
)

// DataPermission is a thin alias of go-admin-core's sdk/contract/actions
// (PRD 006 F3/F5).
type DataPermission = contractactions.DataPermission

// The five values sys_role.data_scope can hold, referenced directly from
// go-admin-core's sdk/contract/actions rather than restated as literals -
// see that package's DataScope* doc comment and PRD 006's hard constraint 4.
const (
	DataScopeAll      = contractactions.DataScopeAll
	DataScopeCustom   = contractactions.DataScopeCustom
	DataScopeDept     = contractactions.DataScopeDept
	DataScopeDeptTree = contractactions.DataScopeDeptTree
	DataScopeSelf     = contractactions.DataScopeSelf
)

// PermissionAction, Permission, GetPermissionFromContext and
// IsValidDataScope forward to go-admin-core's sdk/contract/actions (PRD 006
// F3/F5). create.go/delete.go/index.go/update.go/view.go in this package
// (the generic CRUD actions, which do not move to core) call Permission and
// GetPermissionFromContext by these same names and are unchanged by the
// move: the names now resolve to forwards instead of local definitions, and
// the behaviour is identical either way.
func PermissionAction() gin.HandlerFunc {
	return contractactions.PermissionAction()
}

func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return contractactions.Permission(tableName, p)
}

func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return contractactions.GetPermissionFromContext(c)
}

// IsValidDataScope reports whether s is one of the five values Permission
// recognizes. See go-admin-core's sdk/contract/actions.IsValidDataScope.
func IsValidDataScope(s string) bool {
	return contractactions.IsValidDataScope(s)
}
