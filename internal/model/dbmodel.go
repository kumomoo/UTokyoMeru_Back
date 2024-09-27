package model

import (
	"gorm.io/gorm"
)

type Good struct {
	gorm.Model
	Title       string   `gorm:"not null"`
	Description string   `gorm:"not null"`
	Images      []string `gorm:"not null"`
	Price       float64  `gorm:"not null"`
	Views       uint     `gorm:"default:0"`
	IsInvisible bool     `gorm:"default:False"`
	IsDeleted   bool     `gorm:"default:False"`
	IsBought    bool     `gorm:"default:False"`
	Tags        []string
	UserId      uint
	User        User   `gorm:"foreignKey:UserId"`
	Comments    []Comment `gorm:"foreignkey:GoodId"`
}

type User struct {
	gorm.Model
	Name        string `gorm:"not null"`
	MailAddress string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Avatar      string `gorm:"type:text"`
	Sales       []Good `gorm:"foreignkey:UserId"`
	Buys        []Good `gorm:"foreignkey:UserId"`
	Favorites   []Good   `gorm:"foreignkey:UserId"`
	Comments    []Comment `gorm:"foreignkey:UserId"`
	IsDeleted   bool      `gorm:"not null;default:0"`
	IsBanned    bool      `gorm:"not null;default:0"`
	UserClass   string `gorm:"not null;default:0"`
	Gender      string `gorm:"not null;default:Others"`
	PhoneNumber string `gorm:"not null;default:null"`
	MailCode    string `gorm:"not null;default:null"`
	Address     string `gorm:"not null;default:null"`
}

type Comment struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Content     string `gorm:"not null"`
	ReplyTo     uint   `gorm:"default:null"`
	IsInvisible bool   `gorm:"default:false"`
	IsDeleted   bool   `gorm:"default:false"`
	UserId      uint
	GoodId      uint
	User        User   `gorm:"foreignKey:UserId"`
	Good        Good   `gorm:"foreignKey:GoodId"`
}
