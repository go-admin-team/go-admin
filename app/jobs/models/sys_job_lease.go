package models

// SchedulerLeaseName is the name of the one lease row per database.
//
// One row, not one per tenant: a tenant is a separate database with its own
// sys_job table and its own scheduler, so the row that decides who schedules
// it lives in that database alongside the jobs it governs.
const SchedulerLeaseName = "scheduler"

// SysJobLease is the scheduler's single-writer lease over one database.
//
// app/jobs registers every enabled job into an in-process cron.Cron and keeps
// each job's scheduler handle in sys_job.entry_id. The scheduler is per
// process and entry_id is one shared column, so a second instance pointed at
// the same database does not divide the work - it overwrites it, and nothing
// logs that it did (issue #915). Only the holder of this lease calls
// jobs.Setup, which keeps the scheduler single-writer while the HTTP side
// still scales.
//
// It deliberately embeds neither models.ModelTime nor models.ControlBy. A
// lease is machine state, not a record a person creates, edits or
// soft-deletes: there is no author to attribute it to, and a deleted-but-
// present lease row would be a row that both does and does not hold the
// scheduler.
type SysJobLease struct {
	// Name is the lease being held. The migration seeds exactly one row,
	// SchedulerLeaseName, and the runtime only ever updates it - there is
	// no insert path, so two instances starting at once cannot race to
	// create the row they are both trying to claim.
	Name string `json:"name" gorm:"type:varchar(64);primaryKey"`

	// Owner identifies the process that holds the lease. Empty means the
	// lease is free, which is what the migration seeds.
	Owner string `json:"owner" gorm:"type:varchar(191);not null"`

	// AcquiredAtMs is when the current owner took the lease, not when it
	// last renewed: a leader that has held it for an hour and one that took
	// over a second ago are different situations, and only this column
	// tells them apart. Renewal moves ExpiresAtMs and leaves this alone.
	AcquiredAtMs int64 `json:"acquiredAtMs" gorm:"column:acquired_at_ms;not null"`

	// ExpiresAtMs is when another instance may take the lease.
	//
	// Milliseconds since the Unix epoch, in a BIGINT, rather than a
	// timestamp column. A timestamp crossing the driver boundary carries
	// timezone semantics that the driver applies on the way through: with
	// go-admin's own `parseTime=True&loc=Local` DSN, MySQL's UTC_TIMESTAMP
	// comes back labelled as local time, and a lease written in Asia/
	// Shanghai is then eight hours out - in whichever direction makes every
	// other instance's lease look expired. An integer has no timezone for
	// anything to apply, and the comparison that decides who schedules
	// becomes integer arithmetic that no DSN setting can reinterpret.
	ExpiresAtMs int64 `json:"expiresAtMs" gorm:"column:expires_at_ms;not null"`
}

func (*SysJobLease) TableName() string {
	return "sys_job_lease"
}
