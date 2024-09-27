package config

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
)

type dbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Timezone string `json:"timezone"`
}

var C dbConfig

func init() {
	var f = flag.String("f", "UTokyoMeru_Back/dbconfig.yaml", "config file")
	flag.Parse()
	conf.MustLoad(*f, &C)
}

func (c *dbConfig) ToDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.Dbname, c.Port, c.SSLMode, c.Timezone)
}