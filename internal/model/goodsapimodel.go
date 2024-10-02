package model

// send this when frontend requests goods
type GetGoodsResponse struct {
	GoodID      uint
	Title       string
	Images      []string
	Price       float64
	Description string
	User        UserForGetGoodsResponse
}
type UserForGetGoodsResponse struct {
	Name   string
	Avatar string
}

// receive this when frontend posts or updates goods
type PostGoodsReceive struct {
	Title       string
	Description string
	Images      []string
	Price       float64
	Tags        []string
	SellerID    uint
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
	ID       uint
	SellerID uint
}

// send this when frontend deletes goods
type DeleteGoodsResponse struct {
	Message string
}
