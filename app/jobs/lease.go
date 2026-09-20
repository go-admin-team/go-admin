package jobs

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models2 "go-admin/app/jobs/models"
)

// nowExprMs is the dialect's expression for the current time as
// milliseconds since the Unix epoch.
//
// The lease compares one instance's idea of "expired" against another
// instance's idea of "still mine", so both have to come from the same clock.
// Two processes whose wall clocks differ by more than the lease TTL would
// otherwise both hold it and both schedule - the exact situation the lease
// exists to prevent, and it would look like it was working, because each
// instance's own arithmetic is self-consistent.
//
// Milliseconds rather than a timestamp, because a timestamp does not survive
// the trip through a driver unchanged. MySQL's UTC_TIMESTAMP read over
// go-admin's own `parseTime=True&loc=Local` DSN arrives labelled as local
// time: on a UTC+8 host every lease is eight hours out, and a test that only
// checked the lease logic against itself passes anyway. An epoch integer has
// no timezone for a driver to apply.
func nowExprMs(dialect string) (string, error) {
	switch dialect {
	case "mysql":
		// UNIX_TIMESTAMP reads its argument in the session timezone and
		// NOW(3) is in the session timezone, so the two cancel and the
		// result is the absolute epoch regardless of what that zone is.
		return "CAST(ROUND(UNIX_TIMESTAMP(NOW(3)) * 1000) AS SIGNED)", nil
	case "postgres":
		return "CAST(EXTRACT(EPOCH FROM clock_timestamp()) * 1000 AS BIGINT)", nil
	case "sqlite":
		// julianday is the portable millisecond clock here: strftime('%s')
		// truncates to the second, and unixepoch('now','subsec') needs
		// SQLite 3.42.
		return "CAST((julianday('now') - 2440587.5) * 86400000.0 AS INTEGER)", nil
	case "sqlserver":
		return "DATEDIFF_BIG(millisecond, '1970-01-01T00:00:00', SYSUTCDATETIME())", nil
	}
	return "", fmt.Errorf("no epoch-milliseconds expression for dialect %q", dialect)
}

// dbNowMs reads the clock from the database rather than from this process.
//
// The read and the UPDATE that uses it are two statements, so the value is
// already slightly stale by the time it is compared - and that is the safe
// direction in both places it is used:
//
//   - as the expiry cutoff, a stale-old now makes this instance *less*
//     likely to decide another instance's lease has expired;
//   - as the basis for a new expiry, it makes this instance's own lease
//     expire sooner, so it renews sooner.
//
// Neither error makes two instances hold the lease at once.
//
// Zero is rejected rather than returned. It is what a failed conversion
// looks like, it is before every expiry there will ever be, and an
// implementation that passed it on would read every lease as expired, hand
// it to every instance, and restore the defect this lease fixes with a lease
// table sitting on top of it.
func dbNowMs(db *gorm.DB) (int64, error) {
	expr, err := nowExprMs(db.Dialector.Name())
	if err != nil {
		return 0, err
	}

	var ms int64
	if err := db.Raw("SELECT " + expr).Row().Scan(&ms); err != nil {
		return 0, fmt.Errorf("reading the database clock: %w", err)
	}
	if ms <= 0 {
		return 0, fmt.Errorf("the database clock read as %d from %q", ms, expr)
	}
	return ms, nil
}

// newOwnerID identifies this process in the lease row.
//
// Hostname and pid make a log line answer "which one is it" without a lookup;
// the random suffix is what actually makes it unique, because a container
// restarted under the same name can come back with the same hostname and the
// same pid 1.
func newOwnerID() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "unknown"
	}
	return fmt.Sprintf("%s-%d-%s", host, os.Getpid(), uuid.New().String()[:8])
}

// lease is one instance's claim on scheduling one database's jobs.
type lease struct {
	db    *gorm.DB
	owner string
	ttl   time.Duration
}

// acquire takes the lease or renews one this instance already holds, and
// reports whether this instance holds it when it returns.
//
// Renewal is tried first and is scoped to this owner, so it cannot take a
// lease another instance has meanwhile claimed. Only if that matches nothing
// does it try to take an expired one. Both are single UPDATE statements
// decided by RowsAffected: the database, not this process, arbitrates
// between two instances running this at the same moment.
//
// There is no insert path. The migration seeds the row, so a missing row is
// a broken installation rather than a state to recover from - and it is
// reported as one, instead of being papered over by an insert that two
// instances would race to win.
func (l *lease) acquire() (bool, error) {
	nowMs, err := dbNowMs(l.db)
	if err != nil {
		return false, err
	}
	expiresMs := nowMs + l.ttl.Milliseconds()

	renewed := l.db.Model(&models2.SysJobLease{}).
		Where("name = ? AND owner = ?", models2.SchedulerLeaseName, l.owner).
		Update("expires_at_ms", expiresMs)
	if renewed.Error != nil {
		return false, fmt.Errorf("renewing the scheduler lease: %w", renewed.Error)
	}
	if renewed.RowsAffected > 0 {
		return true, nil
	}

	taken := l.db.Model(&models2.SysJobLease{}).
		Where("name = ? AND expires_at_ms <= ?", models2.SchedulerLeaseName, nowMs).
		Updates(map[string]any{
			"owner":          l.owner,
			"acquired_at_ms": nowMs,
			"expires_at_ms":  expiresMs,
		})
	if taken.Error != nil {
		return false, fmt.Errorf("taking the scheduler lease: %w", taken.Error)
	}
	if taken.RowsAffected > 0 {
		return true, nil
	}

	// Neither statement matched. Either another instance holds an
	// unexpired lease - the ordinary case, and not an error - or the row
	// the migration seeds is gone, which is, and which would otherwise
	// present as jobs silently never running anywhere.
	var rows int64
	if err := l.db.Model(&models2.SysJobLease{}).
		Where("name = ?", models2.SchedulerLeaseName).
		Count(&rows).Error; err != nil {
		return false, fmt.Errorf("checking for the scheduler lease row: %w", err)
	}
	if rows == 0 {
		return false, fmt.Errorf("the %q lease row is missing from %s; run the migrations",
			models2.SchedulerLeaseName, (&models2.SysJobLease{}).TableName())
	}
	return false, nil
}

// release hands the lease back so a successor can take it now instead of
// waiting out the TTL. It is scoped to this owner: an instance that already
// lost the lease must not clear the row its successor is holding.
func (l *lease) release() error {
	res := l.db.Model(&models2.SysJobLease{}).
		Where("name = ? AND owner = ?", models2.SchedulerLeaseName, l.owner).
		Updates(map[string]any{"owner": "", "expires_at_ms": 0})
	if res.Error != nil {
		return fmt.Errorf("releasing the scheduler lease: %w", res.Error)
	}
	return nil
}
