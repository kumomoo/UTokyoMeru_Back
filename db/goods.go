package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	GoodsId      uint   `gorm:"not null;autoIncrement"`
	GoodsName    string `gorm:"not null"`
	AuthorId     string `gorm:"not null"`
	AuthorName   string `gorm:"not null"`
	AuthorAvatar string
	Floor        uint   `gorm:"not null; default:1"`
	Describe     string `gorm:"type:text;not null"`
	Images       string `gorm:"type:text"`
	Price        uint   `gorm:"not null"`
	Location     string `gorm:"type:text"`
	Type         string `gorm:"not null"`
	Views        uint   `gorm:"default:0"`
	IsInvisible  bool   `gorm:"default:False"`
	User         User   `gorm:"foreignKey:AuthorId,AuthorName,AuthorAvatar;references:UserId,UserName,Avatar"`
}

type GoodsGet struct {
	GoodsId     uint
	GoodsName   string
	AuthorId    string
	AuthorName  string
	Describe    string
	Images      string
	Price       uint
	Location    string
	CreatedTime time.Time
	UpdatedTime time.Time
	Avatar      string
	Type        string
	Views       uint
}

type GoodsRequest struct {
	GoodsName  string
	AuthorId   string
	AuthorName string
	Describe   string
	Images     string
	Price      uint
	Location   string
	Type       string
	Views      uint
}

type GoodsCRUD struct{}

func (crud GoodsCRUD) CreateByObject(g *Goods) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	if g == nil {
		return errors.New("Goods is nil")
	}
	result := db.Create(g)
	if result.Error != nil {
		return result.Error
	}

	r := &Goods{
		GoodsId:    g.GoodsId,
		GoodsName:  g.GoodsName,
		AuthorId:   g.AuthorId,
		AuthorName: g.AuthorName,
		Describe:   g.Describe,
		Images:     g.Images,
		Price:      g.Price,
		Location:   g.Location,
		Type:       g.Type,
	}
	result = db.Create(r)
	return result.Error
}

func (crud GoodsCRUD) FindById(id uint) (*Goods, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res Goods
	result := db.First(&res, id)
	return &res, result.Error
}

func (crud GoodsCRUD) FindAll() ([]Goods, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res []Goods
	result := db.Find(&res)
	return res, result.Error
}

func (crud GoodsCRUD) UpdateByObject(g *Goods) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	result := db.Save(&g)
	return result.Error
}

func (crud GoodsCRUD) FindAllOrdered() ([]Goods, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var top []Goods
	var res []Goods
	topGoods := db.Where("is_top = ?", true).Find(&top)
	if topGoods.Error != nil {
		return nil, topGoods.Error
	}

	result := db.Order("updated_at desc").Find(&res)
	res = append(res, top...)
	return res, result.Error
}

func (crud GoodsCRUD) DeleteById(GoodsId string) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	result := db.Delete(&Goods{}, GoodsId)
	return result.Error
}
