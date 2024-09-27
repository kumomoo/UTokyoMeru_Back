package model

import "time"

// send this when frontend requests goods
type GetGoods struct {
	GoodsId     uint
	CreatedTime time.Time
	UpdatedTime time.Time
	Title       string
	Description string
	Images      []string
	Price       float64
	Views       uint
	IsInvisible bool
	IsDeleted   bool
	IsBought    bool
	Tags        []string
	UserId      uint
	User        User
	Comments    []Comment
}

// receive this when frontend posts goods
type PostGoodsReceive struct {
	Title       string
	Description string
	Images      []string
	Price       float64
	Tags        []string
	UserId      uint
}

// send this when frontend posts goods
type PostGoodsResponse struct {
	Success  bool
	Message  string
	GoodInfo Good
}

type GetUser struct {
	Name        string
	MailAddress string
	Password    string
	Avatar      string
	IsDeleted   bool
	IsBanned    bool
	UserClass   string
	Gender      string
	PhoneNumber string
	MailCode    string
	Address     string
}
