package db

import (
	"backend/config"
	"backend/internal/model"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
	dbErr      error
	dsn        string
)

var rdb *redis.Client

var c = context.Background()

// 初始化redis
func RedisInit() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			config.C.RedisConfig.Host,
			config.C.RedisConfig.Port,
		),
		Password: config.C.RedisConfig.Password,  // Redis 密码
		DB:       config.C.RedisConfig.Db,        // Redis 数据库编号
		PoolSize: config.C.RedisConfig.Pool_size, // 连接池大小
	})

	_, err = rdb.Ping(c).Result() // 通过 Ping 来测试连接
	fmt.Printf("connect redis success!\n")
	return
}

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

	err := dbInstance.AutoMigrate(&model.User{}, &model.Good{}, &model.Comment{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	if err := RedisInit(); err != nil {
		log.Fatal("failed to initialize redis: ", err)
	}

	fmt.Println("Database initialized")
}
