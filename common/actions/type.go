package actions

import contractactions "github.com/go-admin-team/go-admin-core/v2/sdk/contract/actions"

// PermissionKey is a direct reference to go-admin-core's sdk/contract/actions
// constant, not a restated literal - see that package's PermissionKey doc
// comment. PRD 006's hard constraint 4 requires this form for exactly this
// symbol: PermissionAction (below) sets the gin context key it owns, and
// GetPermissionFromContext reads it back; an independently declared literal
// here would let the two silently drift apart if core's copy ever changed
// without this one following. common/actions/shim_test.go carries the
// regression test for that failure mode.
const PermissionKey = contractactions.PermissionKey
