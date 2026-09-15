package models

import (
	"time"

	"go-admin/common/models"
)

// The values sys_app.status takes.
//
// Three states rather than a single "installed", because an install that
// stopped partway has to be an observable row rather than the absence of one:
// the versions an app installs are separate migration files, and on MySQL a
// DDL statement commits the transaction around it - taking an outer
// transaction and every savepoint under it with it - so they cannot be
// wrapped in one.
//
// AppInstalling is also what a row reads as after the process was killed
// mid-install, which is why it is not treated as "installed" by anything.
const (
	AppInstalling = 1
	AppInstalled  = 2
	AppFailed     = 3
)

// SysApp is the sys_app row model: one row per installed application (PRD
// 008 F2). It deliberately does not embed models.ModelTime - see the design
// doc (docs-prd/008-应用清单与安装器/数据库变更.md) §1.1 for why an
// installed-app registry does not need the millisecond soft-delete marker
// every other sys_* table follows. Uninstalling an app deletes its row
// outright; a later reinstall creates a fresh one.
type SysApp struct {
	models.Model // Id int, primary key, autoincrement

	// AppCode is the app.Manifest.Code / migration.ForApp / seed.SeedMenus
	// identity, already lower-cased by migration.NormalizeAppCode before
	// anything reaches this table. Unique: row existence alone answers G2
	// ("is app X installed").
	AppCode string `json:"appCode" gorm:"type:varchar(64);not null;uniqueIndex:uk_sys_app_app_code;comment:app code"`

	Name string `json:"name" gorm:"size:128;not null;comment:display name, from Manifest.Name"`
	// Version is the version this row currently reflects - attempted or
	// confirmed, disambiguated by Status. It does not drive which
	// migrations run next; sys_migration's per-version rows do that (see
	// design doc §1.5's resume flow). This field is descriptive, refreshed
	// from the manifest on every install/upgrade/resume attempt.
	Version     string `json:"version" gorm:"size:32;not null;comment:version this row currently reflects, see Status"`
	Description string `json:"description" gorm:"size:255;not null;default:'';comment:from Manifest.Description"`
	Author      string `json:"author" gorm:"size:128;not null;default:'';comment:from Manifest.Author"`

	// Requires is a comma-separated list of app codes this app declared as
	// dependencies (Manifest.Requires). Stored as plain VARCHAR CSV, not
	// JSON - see design doc §1.3 for why. F8 (P1) is what validates and
	// orders on this; this batch only stores what the manifest declared.
	Requires string `json:"requires" gorm:"size:255;not null;default:'';comment:declared dependency app codes, comma separated"`

	// Pricing/License are reserved passthrough fields (PRD 003; PRD 008
	// open question 1). This batch stores whatever the manifest carries and
	// does not interpret either one.
	Pricing string `json:"pricing" gorm:"size:64;not null;default:'';comment:reserved, not interpreted by this batch"`
	License string `json:"license" gorm:"size:64;not null;default:'';comment:reserved, not interpreted by this batch"`

	// Status: 1=installing 2=installed 3=failed. Three states, not a
	// single "1=installed", because a partial, stuck install has to be an
	// observable row rather than "the row doesn't exist yet" - see design
	// doc §1.5 for why cross-migration-file atomicity is not available on
	// MySQL (implicit commit on DDL).
	Status int `json:"status" gorm:"size:4;not null;default:1;comment:1=installing 2=installed 3=failed"`

	// FailedVersion and LastError are DIAGNOSTIC TEXT ONLY - what a human
	// looking at this row is told about the last failure, nothing more. No
	// code anywhere may read either one to decide what to do next.
	//
	// The question "where should a resume pick up" has exactly one
	// authoritative answer, and it is not these two columns: subtract
	// sys_migration's applied rows for this app_code from what the app's
	// own compiled-in code has registered (migration.Snapshot()/ForApp -
	// the same set F7's `migrate status` already walks). That answer can
	// never go stale, because it is not stored anywhere to go stale - it is
	// recomputed from sys_migration every time it is asked. FailedVersion
	// is a snapshot of what that computation returned at the moment of
	// failure, kept only so an operator does not have to go find the
	// process's logs; if it and a fresh recomputation from sys_migration
	// ever disagree, sys_migration is right and this column is stale, by
	// definition, and nothing should ever notice or care except a human
	// reading the row.
	FailedVersion string `json:"failedVersion" gorm:"size:64;not null;default:'';comment:diagnostic snapshot only, not a judgment basis; meaningful only when status=3"`
	LastError     string `json:"lastError" gorm:"size:255;not null;default:'';comment:diagnostic text only, not a judgment basis; meaningful only when status=3"`

	// InstalledAt is when this app first reached status=installed - set
	// once, never moved by a later upgrade (see design doc §1.4). Nullable,
	// unlike every other column here: a row can exist before it has a
	// value (a fresh install starts at status=installing). This is not the
	// deleted_at problem 1786700003000_soft_delete_marker.go fixed - that
	// column sat inside a unique index, where NULL <> NULL let two live
	// rows coexist under the same key. InstalledAt is in no index at all,
	// so nullability here opens no such hole.
	InstalledAt *time.Time `json:"installedAt" gorm:"comment:first successful install time; null until status first reaches installed"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"comment:last updated time"`

	models.ControlBy // CreateBy/UpdateBy: which operator triggered the attempt
}

func (*SysApp) TableName() string {
	return "sys_app"
}
