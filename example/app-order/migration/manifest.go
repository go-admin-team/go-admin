package migration

import (
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/app"
)

// Version is what this application calls itself. A host's installer records
// it, compares it against what is already installed to tell an upgrade from a
// downgrade, and shows it in `migrate status`.
//
// It is not the same thing as the migration version above, and the two move
// independently: adding a migration file without changing what the
// application is called is normal, and so is a release that changes no
// schema. The migration versions decide what runs; this decides what the
// installed row says.
const Version = "1.0.0"

// The manifest is registered from this package rather than one of its own
// because this is the package a host has to import for the application to
// exist at all - the migrations register here too. A second package would be
// a second thing to remember to import, and forgetting it would leave an
// application whose migrations run and which no installer can name.
func init() {
	app.Register(app.Manifest{
		Code:        AppCode,
		Name:        "Order example",
		Version:     Version,
		Description: "A worked example of an application that ships its own tables, menus and APIs.",
		Author:      "go-admin",
		// Nothing yet. When an application does declare dependencies, a host
		// refuses to install it until they are installed - it does not
		// install them for you, because the blast radius of installing one
		// application should not be "and everything it happens to name".
		Requires: nil,
		// Reserved. A host stores both and interprets neither.
		Pricing: "",
		License: "MIT",
	})
}
