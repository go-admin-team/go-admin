package version

import (
	"runtime"

	"gorm.io/gorm"

	jobmodels "go-admin/app/jobs/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
)

// Create sys_job_lease and seed the one row the scheduler competes for
// (issue #915).
//
// The row is seeded here rather than created on demand at startup. Two
// instances starting together would otherwise race to insert the very row
// they are each trying to claim, and the loser would have to tell a
// duplicate-key error apart from a real one in whichever driver it is
// running against. Seeding it makes the runtime path two UPDATE statements
// and nothing else.
//
// It is seeded free - no owner, and an expiry far enough in the past that
// the first instance to ask takes it - so that installing this migration
// does not leave the scheduler waiting out a TTL that nobody is holding.
//
// Ordered after 1786700003000 (the soft-delete conversion), so importing
// cmd/migrate/migration/models is banned here - see
// schema_coverage_test.go's TestPostConversionMigrationsAvoidFrozenSeedModels.
// sys_job_lease is AutoMigrate'd from its runtime model under
// app/jobs/models directly, and it is absent from 1786700003000's frozen
// softDeleteTables list because it embeds no common.ModelTime: a lease that
// could be soft-deleted would be a row that both does and does not hold the
// scheduler.
func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1786700009000JobSchedulerLease)
}

func _1786700009000JobSchedulerLease(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Migrator().AutoMigrate(new(jobmodels.SysJobLease)); err != nil {
			return err
		}

		// Seeded free: no owner, and an expiry of 0 - before every clock
		// reading there will ever be - so the first instance to ask takes
		// it rather than waiting out a TTL nobody is holding.
		lease := jobmodels.SysJobLease{
			Name:         jobmodels.SchedulerLeaseName,
			Owner:        "",
			AcquiredAtMs: 0,
			ExpiresAtMs:  0,
		}
		if err := tx.Create(&lease).Error; err != nil {
			return err
		}

		return tx.Create(&common.Migration{Version: version}).Error
	})
}
