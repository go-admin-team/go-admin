package dto

import (
	"testing"

	vd "github.com/bytedance/go-tagexpr/v2/validator"
)

// api.Bind calls vd.Validate unconditionally on every request, regardless of
// which binding stage ran, so a vd tag on DataScope is enough to reject
// anything actions.Permission's fail-closed default would otherwise have to
// deal with. PRD 006 F14/H2 named this the real trigger for the default
// branch: SysRoleInsertReq.DataScope had no validation at all, so leaving
// dataScope out of a create-role request wrote an empty string straight to
// sys_role.
func TestDataScopeRejectsWhatPermissionCannotRecognize(t *testing.T) {
	invalid := []string{"", "0", "6", "all", " 1", "1 "}
	valid := []string{"1", "2", "3", "4", "5"}

	t.Run("SysRoleInsertReq", func(t *testing.T) {
		for _, s := range invalid {
			req := SysRoleInsertReq{RoleName: "r", RoleKey: "r", DataScope: s}
			if err := vd.Validate(&req); err == nil {
				t.Errorf("DataScope %q was accepted, want rejected", s)
			}
		}
		for _, s := range valid {
			req := SysRoleInsertReq{RoleName: "r", RoleKey: "r", DataScope: s}
			if err := vd.Validate(&req); err != nil {
				t.Errorf("DataScope %q was rejected: %v", s, err)
			}
		}
	})

	t.Run("SysRoleUpdateReq", func(t *testing.T) {
		for _, s := range invalid {
			req := SysRoleUpdateReq{RoleName: "r", RoleKey: "r", DataScope: s}
			if err := vd.Validate(&req); err == nil {
				t.Errorf("DataScope %q was accepted, want rejected", s)
			}
		}
		for _, s := range valid {
			req := SysRoleUpdateReq{RoleName: "r", RoleKey: "r", DataScope: s}
			if err := vd.Validate(&req); err != nil {
				t.Errorf("DataScope %q was rejected: %v", s, err)
			}
		}
	})

	t.Run("RoleDataScopeReq", func(t *testing.T) {
		for _, s := range invalid {
			req := RoleDataScopeReq{RoleId: 1, DataScope: s}
			if err := vd.Validate(&req); err == nil {
				t.Errorf("DataScope %q was accepted, want rejected", s)
			}
		}
		for _, s := range valid {
			req := RoleDataScopeReq{RoleId: 1, DataScope: s}
			if err := vd.Validate(&req); err != nil {
				t.Errorf("DataScope %q was rejected: %v", s, err)
			}
		}
	})
}
