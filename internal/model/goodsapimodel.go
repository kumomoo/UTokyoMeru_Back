package model

// send this when frontend requests goods
type GetGoodsResponse struct {
	GoodID      uint                    `json:"good_id"`
	Title       string                  `json:"title"`
	Images      []string                `json:"images"`
	Price       float64                 `json:"price"`
	Description string                  `json:"description"`
	User        UserForGetGoodsResponse `json:"user"`
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
	SellerID    uint    `json:"seller_id"`
	IsInvisible bool    `json:"is_invisible"`
	IsDeleted   bool    `json:"is_deleted"`
	IsBought    bool    `json:"is_bought"`
}

// send this when frontend posts or updates goods
type PostGoodsResponse struct {
	Message  string
	GoodInfo Good
}

// receive this when frontend deletes goods
type DeleteGoodsReceive struct {
	ID       uint
	SellerID uint `json:"seller_id"`
}

// send this when frontend deletes goods
type DeleteGoodsResponse struct {
	Message string
}
