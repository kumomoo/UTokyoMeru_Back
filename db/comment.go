package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	GoodsId     uint   `gorm:"not null"`
	GoodsName   string `gorm:"not null"`
	AuthorId    string `gorm:"not null"`
	AuthorName  string `gorm:"not null"`
	Content     string `gorm:"not null"`
	Avatar      string `gorm:"type:text"`
	ReplyTo     string `gorm:"default:null"`
	Floor       uint   `gorm:"not null"`
	IsInvisible bool   `gorm:"default:false"`
	User        User   `gorm:"foreignKey:AuthorId,AuthorName;references:UserId,UserName"`
	Goods       Goods  `gorm:"foreignKey:GoodsId,GoodsName;references:GoodsId,GoodsName"`
}

type CommentGet struct {
	ID          uint
	CreatedTime time.Time
	UpdatedTime time.Time
	GoodsId     uint
	GoodsName   string
	AuthorId    string
	AuthorName  string
	Content     string
	Avatar      string
	ReplyTo     string
}

type CommentRequest struct {
	GoodsId  uint
	AuthorId string
	Content  string
	ReplyTo  string
}

type CommentCRUD struct{}

func (crud CommentCRUD) CreateByObject(c *Comment) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("post is nil")
	}

	g := &GoodsCRUD{}
	goods, err := g.FindById(c.GoodsId)
	if err != nil {
		return err
	}

	if goods.Floor > 1 {
		goods.Floor++
		c.Floor = goods.Floor
	} else {
		c.Floor = 1
	}

	err = g.UpdateByObject(goods)
	if err != nil {
		return err
	}

	result := db.Create(c)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (crud CommentCRUD) UpdateByObject(c *Comment) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	result := db.Save(&c)
	return result.Error
}

func (crud CommentCRUD) DeleteById(id uint) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	result := db.Delete(&Goods{}, id)
	return result.Error
}

func (crud CommentCRUD) FindAllByGoodsId(GoodsId uint) ([]Comment, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res []Comment
	result := db.Where("GoodsId = ? AND IsInvisible = ?", GoodsId, false).Order("Floor").Find(&res)
	return res, result.Error
}
func (crud CommentCRUD) FindById(id uint) (*Comment, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res Comment
	result := db.First(&res, id)
	return &res, result.Error
}

func (crud CommentCRUD) SafeDeleteById(id uint) error {
	result, err := crud.FindById(id)
	if err != nil {
		return err
	}

	result.IsInvisible = true
	err = crud.UpdateByObject(result)
	if err != nil {
		return err
	}
	return nil
}
