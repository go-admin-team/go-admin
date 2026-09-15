package database

import (
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// openerFor resolves the gorm dialector registered for driver.
//
// Which drivers exist depends on how the binary was built: sqlite3 needs cgo
// and is only compiled in under the sqlite3 build tag. A configuration naming a
// driver this build does not carry used to reach gorm.Open with a nil function,
// which is a nil dereference deep in the stack rather than an answer.
func openerFor(driver string) (func(string) gorm.Dialector, error) {
	open, ok := opens[driver]
	if !ok {
		return nil, fmt.Errorf("unsupported database driver %q, this build supports %s",
			driver, strings.Join(supportedDrivers(), ", "))
	}
	return open, nil
}

// supportedDrivers lists what this build can open, in a stable order.
func supportedDrivers() []string {
	names := make([]string, 0, len(opens))
	for name := range opens {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
