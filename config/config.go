package config

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	DbConfig    DbConfig
	RedisConfig RedisConfig
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Timezone string `json:"timezone"`
}

type RedisConfig struct {
	Host      string `json:"host"`
	Port      int64  `json:"port"`
	Password  string `json:"password"`
	Db        int    `json:"db"`
	Pool_size int    `json:"pool_size"`
}

var C Config

func init() {
	var f = flag.String("f", "config.yaml", "config file")
	flag.Parse()
	conf.MustLoad(*f, &C)
}

func (c *Config) ToDSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DbConfig.Host, c.DbConfig.User, c.DbConfig.Password, c.DbConfig.Dbname, c.DbConfig.Port, c.DbConfig.SSLMode, c.DbConfig.Timezone)
	//fmt.Println("database config loaded:", dsn)
	return dsn
}
