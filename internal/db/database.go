package db

import (
	"sync"
	"backend/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
	dbErr      error
)

func GetDatabaseInstance() (*gorm.DB, error) {
	dbOnce.Do(func() {
		dsn := config.C.ToDSN()
		dbInstance, dbErr = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	})
	return dbInstance, dbErr
}
