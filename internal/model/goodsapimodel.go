package model

import "time"

// send this when frontend requests goods
type GetGoodsResponse struct {
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

// receive this when frontend posts or updates goods
type PostGoodsReceive struct {
	Title       string
	Description string
	Images      []string
	Price       float64
	Tags        []string
	UserId      uint
	IsInvisible bool
	IsDeleted   bool
	IsBought    bool
}

// send this when frontend posts or updates goods
type PostGoodsResponse struct {
	Message  string
	GoodInfo Good
}

// receive this when frontend deletes goods
type DeleteGoodsReceive struct {
	ID     uint
	UserId uint
}

// send this when frontend deletes goods
type DeleteGoodsResponse struct {
	Message string
}
