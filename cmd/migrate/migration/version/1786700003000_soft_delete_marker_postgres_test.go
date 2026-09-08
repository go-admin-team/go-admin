package version

import (
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgresDSNEnv points these tests at a database. They are skipped without
// it, so a developer with no PostgreSQL running still gets a green run.
//
// The whole file exists because the rest of this package's tests run on
// SQLite, where the defect they cover cannot happen: dropping an index through
// gorm's migrator works there and produces unparseable SQL on PostgreSQL. A
// suite that only ever exercised SQLite reported success for a migration that
// failed on every PostgreSQL database it was pointed at - go-admin#919.
const postgresDSNEnv = "GO_ADMIN_TEST_POSTGRES_DSN"

func postgresDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv(postgresDSNEnv)
	if dsn == "" {
		// Skipping locally is the point; skipping in CI is the failure this
		// file exists to prevent. A workflow that renamed the variable or
		// dropped the service would otherwise go green while these tests
		// quietly did nothing - the same shape as the defect they cover.
		if os.Getenv("CI") != "" {
			t.Fatalf("%s is not set while CI is: the PostgreSQL migration tests must not skip here", postgresDSNEnv)
		}
		t.Skipf("%s is not set; skipping the PostgreSQL migration tests", postgresDSNEnv)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("connecting to %s: %v", postgresDSNEnv, err)
	}
	return db
}

// pgOldUser is the pre-migration shape: a nullable timestamp with an index on
// it, which is what makes dropping the column require dropping the index.
type pgOldUser struct {
	UserId    int64 `gorm:"column:user_id;primaryKey;autoIncrement"`
	Username  string
	DeletedAt *time.Time `gorm:"index"`
}

func (pgOldUser) TableName() string { return "sd_pg_user" }

// The conversion completes on PostgreSQL.
//
// It did not. dropIndexesOn went through Migrator().DropIndex, whose
// PostgreSQL driver falls back to an expression when it cannot resolve a
// schema - which is every call made here, because the migration passes a table
// name as a string:
//
//	DROP INDEX CURRENT_SCHEMA()."idx_sd_pg_user_deleted_at"
//
// DROP INDEX takes an identifier there, so it failed to parse and took the
// whole conversion with it. Every PostgreSQL deployment stopped at this
// migration, and the visible symptom was a login rejecting a correct password
// because deleted_at was still a timestamptz being compared to 0.
func TestConversionCompletesOnPostgres(t *testing.T) {
	db := postgresDB(t)
	t.Cleanup(func() { db.Migrator().DropTable(&pgOldUser{}) })

	db.Migrator().DropTable(&pgOldUser{})
	if err := db.AutoMigrate(&pgOldUser{}); err != nil {
		t.Fatalf("building the old shape: %v", err)
	}

	deleted := time.Now().Add(-time.Hour)
	db.Create(&pgOldUser{Username: "gone", DeletedAt: &deleted})
	db.Create(&pgOldUser{Username: "live"})

	if err := convertDeletedAt(db, "sd_pg_user"); err != nil {
		t.Fatalf("convertDeletedAt: %v", err)
	}

	var dataType string
	db.Raw(`SELECT data_type FROM information_schema.columns
	        WHERE table_name = 'sd_pg_user' AND column_name = 'deleted_at'`).Scan(&dataType)
	if dataType != "bigint" {
		t.Errorf("deleted_at is %q after the conversion, want bigint", dataType)
	}

	// The marker has to carry the timestamp across, or a row that was deleted
	// comes back live.
	var markers []int64
	db.Raw(`SELECT deleted_at FROM sd_pg_user ORDER BY user_id`).Scan(&markers)
	if len(markers) != 2 {
		t.Fatalf("read %d rows, want 2", len(markers))
	}
	if markers[0] == 0 {
		t.Error("the deleted row came back live")
	}
	if markers[1] != 0 {
		t.Errorf("the live row is marked deleted at %d", markers[1])
	}
}

// The index on deleted_at is gone afterwards, which is the step that failed.
//
// Asserted separately from the conversion because the conversion can succeed
// on a table with no index at all, and this migration exists for tables that
// have one.
func TestTheIndexOnDeletedAtIsDroppedOnPostgres(t *testing.T) {
	db := postgresDB(t)
	t.Cleanup(func() { db.Migrator().DropTable(&pgOldUser{}) })

	db.Migrator().DropTable(&pgOldUser{})
	if err := db.AutoMigrate(&pgOldUser{}); err != nil {
		t.Fatalf("building the old shape: %v", err)
	}

	var before int64
	db.Raw(`SELECT count(*) FROM pg_indexes
	        WHERE tablename = 'sd_pg_user' AND indexdef LIKE '%deleted_at%'`).Scan(&before)
	if before == 0 {
		t.Fatal("the old shape has no index on deleted_at, so this test asserts nothing")
	}

	if err := dropIndexesOn(db, "sd_pg_user", "deleted_at"); err != nil {
		t.Fatalf("dropIndexesOn: %v", err)
	}

	var after int64
	db.Raw(`SELECT count(*) FROM pg_indexes
	        WHERE tablename = 'sd_pg_user' AND indexdef LIKE '%deleted_at%'`).Scan(&after)
	if after != 0 {
		t.Errorf("%d index(es) on deleted_at survived", after)
	}
}
