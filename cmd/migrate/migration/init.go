package migration

import (
	"gorm.io/gorm"
	"log"
	"sort"
	"sync"
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
	sort.StringsAreSorted(versions)
	log.Println(versions)
	var err error
	var count int64
	for _, v := range versions {
		err = e.db.Debug().Table("sys_migration").Where("version = ?", v).Count(&count).Error
		if err != nil {
			log.Fatalln(err)
		}
		if count > 0 {
			log.Println(count)
			count = 0
			continue
		}
		err = (e.version[v])(e.db.Debug(), v)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
