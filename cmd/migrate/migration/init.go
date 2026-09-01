package migration

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	common "go-admin/common/models"
)

var Migrate = newMigration()

func newMigration() *Migration {
	return &Migration{version: make(map[string]versionEntry)}
}

// versionEntry is one registered migration plus the app it belongs to. The
// empty app code means the framework itself, which is also what the
// sys_migration.app_code column defaults to, so history written before this
// field existed reads back correctly with no backfill.
type versionEntry struct {
	appCode string
	fn      func(db *gorm.DB, version string) error
}

type Migration struct {
	db      *gorm.DB
	version map[string]versionEntry
	mutex   sync.Mutex
}

func (e *Migration) GetDb() *gorm.DB {
	return e.db
}

func (e *Migration) SetDb(db *gorm.DB) {
	e.db = db
}

// SetVersion registers a migration owned by the framework. Signature and
// behaviour are unchanged: every existing call site in version/*.go keeps
// compiling and keeps writing common.Migration{Version: version} with no app
// code, which is the correct meaning of "framework".
func (e *Migration) SetVersion(k string, f func(db *gorm.DB, version string) error) {
	e.setVersion(k, "", f)
}

func (e *Migration) setVersion(k, appCode string, f func(db *gorm.DB, version string) error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.version[k] = versionEntry{appCode: appCode, fn: f}
}

// AppMigrationFunc is the signature of a migration registered through ForApp.
//
// It receives appCode explicitly because the migration - not the framework -
// writes its own completion row, normally as the last statement inside its own
// transaction. That is what makes "the schema change and the record of it
// commit together" true, and the framework cannot insert the row on the
// migration's behalf without giving that up. Handing the code to the function
// is what stops an app's migrations from silently recording themselves as the
// framework's.
type AppMigrationFunc func(db *gorm.DB, version, appCode string) error

// AppRegistrar is a per-app view over a registry.
type AppRegistrar struct {
	m       *Migration
	appCode string
}

// FrameworkAppCode is the name migrate status prints for migrations that belong
// to the framework rather than to an app, and the name --app accepts to select
// them. The stored app code for those is the empty string; this is only the
// spelling humans use. It is reserved - ForApp rejects it - so that every group
// heading status prints is also a value --app understands.
const FrameworkAppCode = "core"

// ForApp returns a registrar that records migrations under code.
//
// The code is lower-cased: sys_migration.version sorts as ASCII, so mixed case
// would order MyApp before crm for no reason a reader could guess, and the two
// spellings would group as two different apps in migrate status.
//
// An empty or reserved code panics rather than falling back to the framework.
// Registration happens in init(), so this fires the first time the binary runs
// anywhere, which is the point: an app whose migrations quietly file themselves
// under the framework is exactly the class of silent failure this work is meant
// to remove. Framework migrations call Migrate.SetVersion directly.
func ForApp(code string) *AppRegistrar { return Migrate.ForApp(code) }

// ForApp is the same on an explicit registry, which is what tests use.
func (e *Migration) ForApp(code string) *AppRegistrar {
	code = NormalizeAppCode(code)
	switch code {
	case "":
		panic("migration.ForApp: empty app code; framework migrations use Migrate.SetVersion")
	case FrameworkAppCode:
		panic("migration.ForApp: app code " + FrameworkAppCode + " is reserved for the framework")
	}
	return &AppRegistrar{m: e, appCode: code}
}

// AppCode reports the code this registrar files migrations under, after
// normalisation.
func (r *AppRegistrar) AppCode() string { return r.appCode }

// SetVersion registers an app-owned migration under k, which is the bare
// timestamp taken from the file name exactly as framework migrations do.
//
// What reaches sys_migration.version is the namespaced form; the version string
// handed to f is that same namespaced string, so a migration that writes
// common.Migration{Version: version, AppCode: appCode} records the key the
// registry will look for next time.
func (r *AppRegistrar) SetVersion(k string, f AppMigrationFunc) {
	key := namespacedKey(r.appCode, k)
	r.m.setVersion(key, r.appCode, func(db *gorm.DB, version string) error {
		return f(db, version, r.appCode)
	})
}

// namespacedKey scopes k to appCode so two apps cannot collide on the
// sys_migration.version primary key by minting the same millisecond timestamp.
// Framework migrations (appCode == "") stay bare, matching every version string
// already in production.
func namespacedKey(appCode, k string) string {
	if appCode == "" {
		return k
	}
	return appCode + "-" + k
}

// StatusEntry is one row of migrate status.
type StatusEntry struct {
	AppCode    string
	Version    string
	Registered bool
	Applied    bool
	ApplyTime  *time.Time
}

// Status merges the in-process registry with sys_migration, so it reports all
// three shapes at once: registered but not applied, registered and applied, and
// applied while nothing registers it any more - a row left behind by a
// migration file that was deleted, or by an app that was uninstalled.
//
// It only reads. Nothing here creates or alters a table, which is what lets
// both `status` and `--dry-run` run against a database without touching it.
func (e *Migration) Status() ([]StatusEntry, error) {
	if e.db == nil {
		return nil, fmt.Errorf("migration: no database configured")
	}

	e.mutex.Lock()
	registered := make(map[string]string, len(e.version))
	for k, v := range e.version {
		registered[k] = v.appCode
	}
	e.mutex.Unlock()

	applied := make(map[string]common.Migration)
	// A database that has never been migrated has no sys_migration table.
	// Reporting everything as pending is the honest answer there; erroring out
	// would make status useless in exactly the case it is most wanted.
	if e.db.Migrator().HasTable(&common.Migration{}) {
		var rows []common.Migration
		if err := e.db.Find(&rows).Error; err != nil {
			return nil, err
		}
		for _, r := range rows {
			applied[r.Version] = r
		}
	}

	versions := make(map[string]struct{}, len(registered)+len(applied))
	for k := range registered {
		versions[k] = struct{}{}
	}
	for k := range applied {
		versions[k] = struct{}{}
	}
	list := make([]string, 0, len(versions))
	for k := range versions {
		list = append(list, k)
	}
	sort.Strings(list)

	out := make([]StatusEntry, 0, len(list))
	for _, v := range list {
		entry := StatusEntry{Version: v}
		if code, ok := registered[v]; ok {
			entry.Registered = true
			entry.AppCode = code
		}
		if row, ok := applied[v]; ok {
			entry.Applied = true
			t := row.ApplyTime
			entry.ApplyTime = &t
			if !entry.Registered {
				// Nothing registers this version any more, so the database is
				// the only source left for what it belonged to.
				entry.AppCode = row.AppCode
			}
		}
		out = append(out, entry)
	}
	return out, nil
}

// Migrate applies every registered migration that has not been applied yet,
// across all apps. Existing callers are unaffected.
func (e *Migration) Migrate() { e.run(allApps) }

// MigrateApp applies only the migrations registered under appCode. Pass
// FrameworkAppCode for the framework's own migrations.
func (e *Migration) MigrateApp(appCode string) { e.run(AppFilter(appCode)) }

// NormalizeAppCode applies the same rule ForApp does, so a code typed on the
// command line matches one written in an init().
func NormalizeAppCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

// AppFilter turns a code as typed into the code stored in the registry, so
// "core" selects the framework's migrations, whose stored code is empty.
func AppFilter(code string) string {
	code = NormalizeAppCode(code)
	if code == FrameworkAppCode {
		return ""
	}
	return code
}

// DisplayAppCode is the inverse: what to print for a stored code.
func DisplayAppCode(code string) string {
	if code == "" {
		return FrameworkAppCode
	}
	return code
}

// AppCodes lists the app codes with at least one registered migration, framework
// included under its display name, sorted.
func (e *Migration) AppCodes() []string {
	e.mutex.Lock()
	seen := map[string]struct{}{}
	for _, v := range e.version {
		seen[DisplayAppCode(v.appCode)] = struct{}{}
	}
	e.mutex.Unlock()

	out := make([]string, 0, len(seen))
	for code := range seen {
		out = append(out, code)
	}
	sort.Strings(out)
	return out
}

func (e *Migration) run(appCode string) {
	e.mutex.Lock()
	versions := make([]string, 0, len(e.version))
	entries := make(map[string]versionEntry, len(e.version))
	for k, v := range e.version {
		if appCode != allApps && v.appCode != appCode {
			continue
		}
		versions = append(versions, k)
		entries[k] = v
	}
	e.mutex.Unlock()
	sort.Strings(versions)

	// A mistyped --app would otherwise select nothing and report "no
	// migrations to apply", which reads exactly like "already up to date".
	if appCode != allApps && len(versions) == 0 {
		log.Printf("no migrations are registered for app %q; registered: %s",
			DisplayAppCode(appCode), strings.Join(e.AppCodes(), ", "))
		return
	}

	var err error
	var count int64
	applied := 0
	for _, v := range versions {
		err = e.db.Table("sys_migration").Where("version = ?", v).Count(&count).Error
		if err != nil {
			log.Fatalln(err)
		}
		if count > 0 {
			// Already applied. This used to print the bare count, so a mature
			// database wrote a screen of "1" at every start.
			count = 0
			continue
		}
		log.Printf("applying migration %s", v)
		if err = entries[v].fn(e.db.Debug(), v); err != nil {
			log.Fatalf("migration %s failed: %v", v, err)
		}
		applied++
	}
	if applied == 0 {
		log.Println("no migrations to apply")
	} else {
		log.Printf("applied %d migration(s)", applied)
	}
}

// allApps is the sentinel run() takes to mean "do not filter". It is distinct
// from the empty app code, which selects the framework's own migrations.
const allApps = "\x00all"

func GetFilename(s string) string {
	s = filepath.Base(s)
	return s[:13]
}
