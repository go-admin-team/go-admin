package migrate

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/app"
	"gorm.io/gorm"

	adminmodels "go-admin/app/admin/models"
	"go-admin/cmd/migrate/migration"
)

// engine is the part of the migration engine the installer drives.
//
// An interface rather than *migration.Migration because the concrete type is
// a package-level singleton with no exported constructor, so a test that took
// it would be sharing one registry with every other test in the process.
type engine interface {
	SetDb(*gorm.DB)
	Status() ([]migration.StatusEntry, error)
	MigrateApp(string) error
}

// installReport is what an install did, for the command to print.
type installReport struct {
	Code string
	// Version is the manifest version this run recorded.
	Version string
	// Previous is the version sys_app held before this run, empty when this
	// is the first install.
	Previous string
	// Applied lists the versions this run brought in, in the order they were
	// applied. Empty on a no-op, and also empty on a run that only corrected
	// sys_app - the difference is NoOp.
	Applied []string
	// NoOp says nothing was left to do: the app is recorded as installed, at
	// this same version, with no migration outstanding.
	NoOp bool
}

// install brings one application up to the version its manifest declares.
//
// Three phases, each committing on its own. They are not one transaction and
// cannot be: an application's versions are separate migration files, and a
// DDL statement inside any of them commits the transaction around it on
// MySQL, which destroys an outer transaction and every savepoint taken from
// it. So this does not
// promise that a half-installed application cannot happen. It promises that
// one is visible when it does: phase A writes "installing" before anything
// that can fail, and phase C turns that into "installed" or "failed".
//
// What is left to apply comes from sys_migration, never from sys_app.
// sys_app is a derived view - a summary for a human, and the answer to "which
// version does this app think it is at". If it were the authority, then an
// operator who deleted sys_migration rows by hand would be told an app is
// installed while its schema is not, which is worse than not knowing.
func install(db *gorm.DB, eng engine, m app.Manifest) (installReport, error) {
	code := migration.NormalizeAppCode(m.Code)
	rep := installReport{Code: code, Version: m.Version}
	if code == "" {
		return rep, errors.New("the manifest declares no app code")
	}
	if code == migration.FrameworkAppCode {
		// Installing the framework is what `migrate` is, and the framework
		// has no manifest and no sys_app row. Saying so beats writing a row
		// that nothing else in this batch expects to exist.
		return rep, fmt.Errorf("%q is the framework's own migrations, not an application; run `migrate` for those", code)
	}
	if !db.Migrator().HasTable(&adminmodels.SysApp{}) {
		return rep, errors.New("sys_app does not exist; run `migrate` first to bring the framework's own tables up to date")
	}

	eng.SetDb(db)

	row, found, err := loadApp(db, code)
	if err != nil {
		return rep, err
	}
	// sameVersion is only meaningful when found; it stays false otherwise.
	// The comparison happens here, before phase A, so an unparseable
	// recorded version is refused while it is still readable rather than
	// after being overwritten.
	sameVersion := false
	if found {
		rep.Previous = row.Version
		cmp, err := app.Compare(m.Version, row.Version)
		if err != nil {
			return rep, fmt.Errorf("comparing %s against the recorded %s: %w", m.Version, row.Version, err)
		}
		if cmp < 0 {
			return rep, fmt.Errorf("%s is recorded at %s; installing %s would be a downgrade, which is not supported",
				code, row.Version, m.Version)
		}
		sameVersion = cmp == 0
	}

	pending, err := pendingFor(eng, code)
	if err != nil {
		return rep, err
	}

	// Nothing outstanding, recorded as installed, at this same version. All
	// three, and the first one comes from sys_migration: a row that says
	// installed while a migration of its has never run is exactly the case
	// sys_app must not be believed about. AppInstalling is not installed -
	// it is what a row reads as after the process was killed partway.
	if found && sameVersion && row.Status == adminmodels.AppInstalled && len(pending) == 0 {
		rep.NoOp = true
		return rep, nil
	}

	// Phase A: the attempt is on disk before anything that can fail.
	now := time.Now()
	if err := beginInstall(db, &row, m, code, found, now); err != nil {
		return rep, err
	}

	// Phase B: no atomicity across these, by the nature of the thing.
	runErr := eng.MigrateApp(code)

	// Phase C.
	if runErr != nil {
		failed := ""
		var vf *migration.VersionFailure
		if errors.As(runErr, &vf) {
			failed = vf.Version
		}
		if err := markFailed(db, code, failed, runErr, time.Now()); err != nil {
			return rep, errors.Join(runErr, fmt.Errorf("recording the failure on sys_app: %w", err))
		}
		return rep, runErr
	}
	if err := markInstalled(db, code, row.InstalledAt, time.Now()); err != nil {
		return rep, err
	}
	rep.Applied = pending
	return rep, nil
}

// loadApp reads the sys_app row for code. A missing row is not an error: it
// is what a first install looks like.
func loadApp(db *gorm.DB, code string) (adminmodels.SysApp, bool, error) {
	var row adminmodels.SysApp
	err := db.Where("app_code = ?", code).First(&row).Error
	if err == nil {
		return row, true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return adminmodels.SysApp{}, false, nil
	}
	return adminmodels.SysApp{}, false, fmt.Errorf("reading sys_app for %q: %w", code, err)
}

// loadApps reads every sys_app row, keyed by app code.
//
// A database that has never had 1786700007000 applied has no such table, and
// that is not an error here: `migrate status` has to keep working on a
// database that has not been migrated at all, which is when it is most wanted.
// A nil map is the honest answer there, and the caller prints what it always
// printed.
func loadApps(db *gorm.DB) (map[string]adminmodels.SysApp, error) {
	if !db.Migrator().HasTable(&adminmodels.SysApp{}) {
		return nil, nil
	}
	var rows []adminmodels.SysApp
	if err := db.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("reading sys_app: %w", err)
	}
	out := make(map[string]adminmodels.SysApp, len(rows))
	for _, r := range rows {
		out[r.AppCode] = r
	}
	return out, nil
}

// pendingFor is the authoritative answer to "what is left to apply", and it
// is recomputed every time rather than stored: what is registered in this
// process, minus what sys_migration says has run. sys_app.failed_version is a
// snapshot of what this returned once and may be stale by now; nothing may
// read it to decide this.
func pendingFor(eng engine, code string) ([]string, error) {
	entries, err := eng.Status()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.AppCode == code && e.Registered && !e.Applied {
			out = append(out, e.Version)
		}
	}
	sort.Strings(out)
	return out, nil
}

// beginInstall is phase A. It refreshes every descriptive column from the
// manifest, because those are the manifest's to say and the row is only a
// copy, and it clears the two diagnostic columns so a stale failure from a
// previous attempt cannot be read as this one's.
func beginInstall(db *gorm.DB, row *adminmodels.SysApp, m app.Manifest, code string, found bool, now time.Time) error {
	row.AppCode = code
	row.Name = m.Name
	row.Version = m.Version
	row.Description = m.Description
	row.Author = m.Author
	row.Requires = strings.Join(m.Requires, ",")
	row.Pricing = m.Pricing
	row.License = m.License
	row.Status = adminmodels.AppInstalling
	row.FailedVersion = ""
	row.LastError = ""
	row.UpdatedAt = now
	if !found {
		if err := db.Create(row).Error; err != nil {
			return fmt.Errorf("recording the install attempt for %q: %w", code, err)
		}
		return nil
	}
	if err := db.Save(row).Error; err != nil {
		return fmt.Errorf("recording the install attempt for %q: %w", code, err)
	}
	return nil
}

// markInstalled is the success half of phase C. installed_at is set once and
// never moved: an upgrade keeps the time of the first install, which is what
// the column is for.
//
// Computed here rather than with COALESCE so the statement is the same on all
// four drivers this repository supports.
func markInstalled(db *gorm.DB, code string, installedAt *time.Time, now time.Time) error {
	updates := map[string]any{
		"status":     adminmodels.AppInstalled,
		"updated_at": now,
	}
	if installedAt == nil {
		updates["installed_at"] = now
	}
	err := db.Model(&adminmodels.SysApp{}).Where("app_code = ?", code).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("recording %q as installed: %w", code, err)
	}
	return nil
}

// markFailed is the other half. Both columns it writes are diagnostic text
// for whoever reads the row; no code may branch on either one.
func markFailed(db *gorm.DB, code, failedVersion string, cause error, now time.Time) error {
	updates := map[string]any{
		"status":         adminmodels.AppFailed,
		"failed_version": truncate(failedVersion, 64),
		"last_error":     truncate(cause.Error(), 255),
		"updated_at":     now,
	}
	return db.Model(&adminmodels.SysApp{}).Where("app_code = ?", code).Updates(updates).Error
}

// truncate cuts s to at most n runes, not bytes: these columns are declared in
// characters, and a message that is partly Chinese would otherwise be cut in
// the middle of one and stored as an invalid sequence.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// reportInstall prints what happened, and says that the data is in place but
// the code is not.
//
// That last sentence is not a pleasantry. Go links its applications at build
// time and Vite resolves its import globs at build time, so installing an
// application writes its menus, its APIs and its permissions and cannot make
// one line of its code run. An operator who is not told that sees the menus
// appear and reasonably concludes the thing is live.
func reportInstall(w io.Writer, rep installReport) {
	if rep.NoOp {
		fmt.Fprintf(w, "%s %s is already installed; nothing to do\n", rep.Code, rep.Version)
		return
	}
	switch {
	case rep.Previous == "":
		fmt.Fprintf(w, "installed %s %s\n", rep.Code, rep.Version)
	case rep.Previous == rep.Version:
		fmt.Fprintf(w, "brought %s %s the rest of the way\n", rep.Code, rep.Version)
	default:
		fmt.Fprintf(w, "upgraded %s from %s to %s\n", rep.Code, rep.Previous, rep.Version)
	}
	if len(rep.Applied) > 0 {
		fmt.Fprintf(w, "applied %d migration(s): %s\n", len(rep.Applied), strings.Join(rep.Applied, ", "))
	} else {
		fmt.Fprintln(w, "no migration was outstanding; only sys_app was brought up to date")
	}
	fmt.Fprintln(w, "the database is up to date, but the application's code is not running yet:")
	fmt.Fprintln(w, "rebuild and restart the server before expecting its routes to answer.")
}

// manifestFor finds the manifest an application registered for code.
//
// A code nothing registered is an error naming what is registered, for the
// same reason exitUnlessAppRegistered exists: the alternative is telling an
// operator who typed `install ordr` that there was nothing to do.
func manifestFor(code string) (app.Manifest, error) {
	want := migration.NormalizeAppCode(code)
	all := app.Snapshot()
	if m, ok := all[want]; ok {
		return m, nil
	}
	codes := make([]string, 0, len(all))
	for c := range all {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	if len(codes) == 0 {
		// Worth its own sentence: no application is compiled into this
		// binary at all, which is a different thing from having typed the
		// wrong one of several.
		return app.Manifest{}, fmt.Errorf(
			"no application registers a manifest in this binary, so %q cannot be installed; "+
				"an application has to be compiled in before it can be installed", want)
	}
	return app.Manifest{}, fmt.Errorf("no application registers the code %q; registered: %s",
		want, strings.Join(codes, ", "))
}
