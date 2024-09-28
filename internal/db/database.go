package db

import (
	"backend/config"
	"backend/internal/model"
	"sync"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
	dbErr      error
	dsn        string
)

func GetDatabaseInstance() (*gorm.DB, error) {
	dbOnce.Do(func() {
		dbInstance, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})
	return dbInstance, dbErr
}

func init() {
	dsn = config.C.ToDSN()
	dbInstance, dbErr = GetDatabaseInstance()
	if dbErr != nil {
		panic(dbErr)
	}

	dbInstance.AutoMigrate(&model.User{},&model.Good{},&model.Comment{})

	fmt.Println("Database initialized")
}
