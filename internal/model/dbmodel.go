package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Good struct {
	gorm.Model
	Title       string         `gorm:"not null"`
	Description string         `gorm:"not null"`
	Images      pq.StringArray `gorm:"type:text[]"`
	Price       float64        `gorm:"not null"`
	Views       uint           `gorm:"default:0"`
	Likes       uint           `gorm:"default:0"`
	IsInvisible bool           `gorm:"default:false"`
	IsDeleted   bool           `gorm:"default:false"`
	IsBought    bool           `gorm:"default:false"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	SellerID    uint           `gorm:"not null"`
	BuyerID     uint           `gorm:"default:null"`
	Comments    []Comment      `gorm:"foreignKey:GoodID"`
}

type User struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	MailAddress string    `gorm:"unique;not null"`
	Password    string    `gorm:"not null"`
	Avatar      string    `gorm:"type:text"`
	Sales       []Good    `gorm:"foreignKey:SellerID"`
	Buys        []Good    `gorm:"foreignKey:BuyerID"`
	FavoList    []Good    `gorm:"many2many:user_likes"`
	Comments    []Comment `gorm:"foreignKey:UserID"`
	IsDeleted   bool      `gorm:"not null;default:false"`
	IsBanned    bool      `gorm:"not null;default:false"`
	UserClass   string    `gorm:"not null;default:user"`
	Gender      string    `gorm:"not null;default:Others"`
	Birthday    time.Time `gorm:"default:null"`
	PhoneNumber string    `gorm:"default:null"`
	MailCode    string    `gorm:"default:null"`
	Address     string    `gorm:"default:null"`
	Rating      float32   `gorm:"default:0"`
	RatingCount float32   `gorm:"default:0"`
	Token       string
}

type Comment struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Content     string `gorm:"not null"`
	ReplyTo     uint   `gorm:"default:null"`
	IsInvisible bool   `gorm:"default:false"`
	IsDeleted   bool   `gorm:"default:false"`
	UserID      uint
	GoodID      uint
}
