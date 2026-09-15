package database

import (
	"net/url"
	"regexp"
	"strings"
)

// A DSN carries the database password, and the startup line used to print it
// whole. Anyone with the log file, a log shipper, or a screenshot of a terminal
// had the credential.
//
// Both shapes this project accepts are covered:
//
//	mysql      user:password@tcp(host:3306)/db?params
//	postgres   postgres://user:password@host:5432/db?params
var mysqlDSN = regexp.MustCompile(`^([^:/@]+):([^@]*)@`)

// redactDSN returns a DSN safe to log: everything but the password.
func redactDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	// URL form, used by postgres and sqlserver.
	if strings.Contains(dsn, "://") {
		u, err := url.Parse(dsn)
		if err != nil {
			// Unparseable and possibly holding a password: say nothing about it.
			return "[dsn]"
		}
		if u.User != nil {
			if _, hasPassword := u.User.Password(); hasPassword {
				u.User = url.UserPassword(u.User.Username(), "***")
			}
		}
		return u.String()
	}

	// user:password@tcp(...) form, used by mysql.
	if m := mysqlDSN.FindStringSubmatch(dsn); m != nil {
		return m[1] + ":***@" + dsn[len(m[0]):]
	}

	// No credential recognised - a sqlite path, for instance.
	return dsn
}
