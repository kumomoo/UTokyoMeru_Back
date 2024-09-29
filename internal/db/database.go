package db

import (
	"backend/config"
	"backend/internal/model"
	"fmt"
	"log"
	"sync"

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

	err := dbInstance.AutoMigrate(&model.User{},&model.Good{},&model.Comment{})
	if err!=nil {
		log.Fatal("failed to migrate database: ", err)
	}

	fmt.Println("Database initialized")
}
