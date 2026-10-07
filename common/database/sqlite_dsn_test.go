package database

import (
	"net/url"
	"testing"
)

func TestSqliteDSN(t *testing.T) {
	parse := func(s string) (string, url.Values) {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
		return u.Path, u.Query()
	}

	t.Run("plain path gets every default", func(t *testing.T) {
		path, q := parse(sqliteDSN("./go-admin-db.db"))
		if path != "./go-admin-db.db" {
			t.Errorf("path = %q", path)
		}
		for _, p := range sqlitePragmas {
			if q.Get(p[0]) != p[1] {
				t.Errorf("%s = %q, want %q", p[0], q.Get(p[0]), p[1])
			}
		}
	})

	t.Run("explicit values win", func(t *testing.T) {
		_, q := parse(sqliteDSN("a.db?_busy_timeout=500&_journal_mode=DELETE"))
		if q.Get("_busy_timeout") != "500" || q.Get("_journal_mode") != "DELETE" {
			t.Errorf("overridden: %v", q)
		}
		if q.Get("_txlock") != "immediate" {
			t.Errorf("missing default alongside overrides: %v", q)
		}
	})

	t.Run("files that only look like memory databases get the defaults", func(t *testing.T) {
		for _, s := range []string{"mode=memory.db", "/tmp/:memory:.db", "data/mode=memory/app.db"} {
			if got := sqliteDSN(s); got == s {
				t.Errorf("sqliteDSN(%q) was left untouched", s)
			}
		}
	})

	t.Run("in-memory and empty untouched", func(t *testing.T) {
		for _, s := range []string{"", ":memory:", "file::memory:?cache=shared", "file:x?mode=memory"} {
			if got := sqliteDSN(s); got != s {
				t.Errorf("sqliteDSN(%q) = %q", s, got)
			}
		}
	})
}
