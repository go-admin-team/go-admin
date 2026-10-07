package database

import (
	"net/url"
	"strings"

	"gorm.io/gorm"
)

// sqlitePragmas are the connection parameters a file-backed sqlite database
// needs to behave under a server that issues concurrent requests. Each one is
// applied only when the configured source does not already set it.
//
// Without them sqlite runs in rollback-journal mode with no busy timeout: a
// writer blocks every reader and a second writer fails at once with
// "database is locked", instead of waiting for the first to finish.
//
//	_busy_timeout    wait up to 10s for a lock rather than failing immediately
//	_journal_mode    WAL lets readers proceed while a write is in progress
//	_synchronous     NORMAL is the pairing WAL is designed for
//	_txlock          a transaction takes the write lock when it begins. The
//	                 default defers it, and a deferred transaction that reads
//	                 and then writes is refused with SQLITE_BUSY straight away -
//	                 busy_timeout does not apply to that upgrade.
var sqlitePragmas = [][2]string{
	{"_busy_timeout", "10000"},
	{"_journal_mode", "WAL"},
	{"_synchronous", "NORMAL"},
	{"_txlock", "immediate"},
}

// sqliteDSN returns source with the defaults in sqlitePragmas added.
// In-memory databases are returned untouched: WAL has no meaning there, and
// every connection to one is its own database anyway.
func sqliteDSN(source string) string {
	if source == "" || strings.Contains(source, ":memory:") || strings.Contains(source, "mode=memory") {
		return source
	}

	path, rawQuery := source, ""
	if i := strings.IndexByte(source, '?'); i >= 0 {
		path, rawQuery = source[:i], source[i+1:]
	}
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		// Not something to rewrite blind; the driver will report it.
		return source
	}
	for _, p := range sqlitePragmas {
		if _, set := q[p[0]]; !set {
			q.Set(p[0], p[1])
		}
	}
	return path + "?" + q.Encode()
}

// withSqliteDefaults wraps open so every connection string handed to a sqlite3
// dialector - the primary and any replicas - carries the defaults above.
func withSqliteDefaults(driver string, open func(string) gorm.Dialector) func(string) gorm.Dialector {
	if driver != "sqlite3" {
		return open
	}
	return func(dsn string) gorm.Dialector { return open(sqliteDSN(dsn)) }
}
