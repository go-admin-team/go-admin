package migration

import (
	"log"
	"path/filepath"
	"sort"
	"sync"

	"gorm.io/gorm"
)

var Migrate = &Migration{
	version: make(map[string]func(db *gorm.DB, version string) error),
}

type Migration struct {
	db      *gorm.DB
	version map[string]func(db *gorm.DB, version string) error
	mutex   sync.Mutex
}

func (e *Migration) GetDb() *gorm.DB {
	return e.db
}

func (e *Migration) SetDb(db *gorm.DB) {
	e.db = db
}

func (e *Migration) SetVersion(k string, f func(db *gorm.DB, version string) error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.version[k] = f
}

func (e *Migration) Migrate() {
	versions := make([]string, 0)
	for k := range e.version {
		versions = append(versions, k)
	}
	if !sort.StringsAreSorted(versions) {
		sort.Strings(versions)
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
		if err = (e.version[v])(e.db.Debug(), v); err != nil {
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

func GetFilename(s string) string {
	s = filepath.Base(s)
	return s[:13]
}
