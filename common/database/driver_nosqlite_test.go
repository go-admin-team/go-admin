//go:build !sqlite3

package database

import (
	"strings"
	"testing"
)

// A default build has no sqlite3, and a config naming it is the usual first
// run. The error has to say how to get the driver, not only that it is missing.
func TestOpenerForTellsHowToGetSqliteInADefaultBuild(t *testing.T) {
	_, err := openerFor("sqlite3")
	if err == nil {
		t.Fatal("sqlite3 was accepted by a build without the sqlite3 tag")
	}
	if !strings.Contains(err.Error(), "-tags sqlite3") {
		t.Errorf("error does not say how to enable sqlite3: %s", err)
	}
}
