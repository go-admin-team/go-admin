package config

import (
	"os"
	"regexp"
	"testing"
)

// The seed files write the built-in admin role's data_scope inline in a SQL
// INSERT, not through Go code, so nothing else in the test suite exercises
// this value. It has to be one of the five scopes actions.Permission
// recognizes: PRD 006 F14/H2 made every other value match no rows, and an
// empty string - which is what these files shipped before that fix - is one
// such value. Without this the shipped admin account would silently lose
// all visibility the moment a deployment turns EnableDP on.
func TestSeedAdminRoleHasAValidDataScope(t *testing.T) {
	cases := map[string]*regexp.Regexp{
		"db.sql": regexp.MustCompile(
			`INSERT INTO sys_role VALUES \(1, '系统管理员', '2', 'admin', 1, '', '', true, '([^']*)'`),
		"db-sqlserver.sql": regexp.MustCompile(
			`\(1, '系统管理员', '2', 'admin', 1, '', '', 1, '([^']*)'`),
	}
	valid := map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true}

	for file, pattern := range cases {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		m := pattern.FindSubmatch(data)
		if m == nil {
			t.Fatalf("%s: admin role INSERT not found; the regex may be out of date", file)
		}
		if scope := string(m[1]); !valid[scope] {
			t.Errorf("%s: admin role data_scope = %q, want one of 1-5", file, scope)
		}
	}
}
