// Command e2e-apporder is a go-admin binary with the example application
// linked into it.
//
// It exists because there is otherwise nothing to install. `migrate install
// order` looks the code up in the registries an application fills from its
// own init(), and an application the binary was not built with registers
// nothing - so every path through the installer past its first check went
// untested, and so did the seeding, the ledger and the uninstall against
// anything but a hand-built fixture.
//
// This lives in its own module, not behind a build tag in the main one. A
// tagged import there would still be resolved by `go mod tidy`, which
// considers every build tag and would go looking for
// github.com/go-admin-team/example-app-order on the network - a repository
// that does not exist, because the example is a directory inside this one.
// A separate module with replace directives is invisible to the main
// module's tidy, its build and its tests, and needs no go.work.
package main

import (
	_ "github.com/go-admin-team/example-app-order/migration"

	"go-admin/cmd"
)

func main() {
	cmd.Execute()
}
