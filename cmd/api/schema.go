package api

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"gorm.io/gorm"

	"go-admin/cmd/migrate/migration"
	"go-admin/common/health"
	commonmodels "go-admin/common/models"
)

// schemaCheckName is what a failing schema reports itself as in /ready's body.
const schemaCheckName = "schema"

// registerSchemaCheck adds the pending-migration check to readiness.
//
// Readiness rather than a refusal to start, and rather than a log line alone.
// The two probes answer different questions: liveness is "restart me", and a
// process whose database is on the wrong schema comes back to the same schema,
// so restarting is not the answer. Readiness is "send me requests", and with a
// schema the binary does not match the answer is no.
//
// Issue #919 is what the absence of this looked like: the process started,
// both probes passed, and the first sign of trouble was a login failing with a
// driver-level encoding error. Refusing to start would have been the wrong fix
// - a process that exits tells an operator less than one that runs and says
// why, and under an orchestrator it crash-loops - while a log line alone is
// not something an orchestrator can act on.
func registerSchemaCheck() {
	health.Register(schemaCheckName, schemaCheck)
}

// schemaCheck fails while any tenant database is behind the migrations this
// binary registers.
//
// Any one of them, rather than only the tenant being served: migrations are
// applied to every database in one run, so one database behind means that run
// did not finish. Serving the rest would let a half-applied deploy look like a
// partial success.
//
// Evaluated per request rather than decided at start-up, so that running
// migrate clears it without a restart.
func schemaCheck(ctx context.Context) error {
	registered := migration.RegisteredVersions()
	if len(registered) == 0 {
		// Nothing registered means nothing can be pending, which is the honest
		// answer for a tree with no migrations. It is also what a broken build
		// would produce - the registry is filled by init() in packages the
		// binary has to link - so cmd/api's dependency test asserts the real
		// binary links them.
		return nil
	}

	behind := make([]string, 0, 2)
	for name, db := range sdk.Runtime.GetAllDb() {
		applied, err := appliedVersions(ctx, db)
		if err != nil {
			return fmt.Errorf("reading applied migrations for %q: %w", name, err)
		}
		if pending := pendingVersions(registered, applied); len(pending) > 0 {
			behind = append(behind, fmt.Sprintf("%s is %d behind, first pending %s",
				name, len(pending), pending[0]))
		}
	}
	if len(behind) == 0 {
		return nil
	}
	sort.Strings(behind)
	return fmt.Errorf("%s; run `go-admin migrate -c <config>` and see `go-admin migrate status`",
		strings.Join(behind, "; "))
}

// appliedVersions reads what sys_migration records for one database.
//
// A missing table is not an error: a database that has never been migrated has
// applied nothing, which is exactly what the caller needs to hear, and is the
// state a first deploy is in.
func appliedVersions(ctx context.Context, db *gorm.DB) (map[string]bool, error) {
	if db == nil {
		return nil, fmt.Errorf("no database")
	}
	db = db.WithContext(ctx)
	if !db.Migrator().HasTable(&commonmodels.Migration{}) {
		return map[string]bool{}, nil
	}
	var rows []commonmodels.Migration
	if err := db.Select("version").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(rows))
	for _, r := range rows {
		out[r.Version] = true
	}
	return out, nil
}

// pendingVersions returns the registered versions applied does not contain.
//
// Split out and taking both sides as arguments because the registry is
// process-wide and filled by init() in packages cmd/api does not import: a
// test in this package cannot arrange it, so the arranging part is the part
// that is not tested here.
func pendingVersions(registered []string, applied map[string]bool) []string {
	out := make([]string, 0)
	for _, v := range registered {
		if !applied[v] {
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
