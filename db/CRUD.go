package db

type InterfaceOfGoodsCRUD interface {
	CreateByObject(g *Goods) error
	FindById(id uint) (*Goods, error)
	FindAll() ([]Goods, error)
	UpdateByObject(g *Goods) error
	FindAllOrdered() ([]Goods, error)
	DeleteById(GoodsId string) error
}

type InterfaceOfCommentCRUD interface {
	CreateByObject(c *Comment) error
	UpdateByObject(c *Comment) error
	DeleteById(id uint) error
	FindAllByGoodsId(GoodsId uint) ([]Comment, error)
	SafeDeleteById(id uint) error
	FindById(id uint) (*Comment, error)
}

type InterfaceOfUserCRUD interface {
	CreateByObject(u *User) error
	GetUserByName(name string) ([]User, error)
	UpdateByObject(u *User) error
	GetUserById(id string) (*User, error)
	FuzzyGetUserByName(name string) ([]User, error)
	GetAllUser() ([]User, error)
	GetAllUserOrdered() ([]User, error)
	DeleteUserbyName(name string) error
}
