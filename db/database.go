package db

import (
	"sync"

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
		dsn := "Meru.db" // SQLite 数据库文件名
		dbInstance, dbErr = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	})
	return dbInstance, dbErr
}
