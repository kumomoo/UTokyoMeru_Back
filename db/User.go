package db

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string  `gorm:"not null"`
	PassWord   string  `gorm:"not null"`
	Avatar     string  `gorm:"type:text"`
	UserId     string  `gorm:"unique;not null"` //mailaddress
	Goods      []Goods `gorm:"foreignkey:AuthorId"`
	NumOfGoods uint    `gorm:"not null;default:0"`
	IsDeleted  bool    `gorm:"not null;default:0"`
	UserClass  string  `gorm:"not null;default:0"`
	Gender     string  `gorm:"not null;default:Others"`
}

type UserGet struct {
	UserName   string
	Avatar     string
	Goods      []Goods
	NumOfGoods uint
	IsDeleted  bool
	UserClass  string
	Gender     string
}

type UserLogin struct {
	UserId   string `json:"UserId" binding:"required"`
	PassWord string `json:"Password" binding:"required"`
}

type UserCRUD struct{}

func (crud UserCRUD) CreateByObject(u *User) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("User not exists")
	}
	result := db.Create(u)
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

func (crud UserCRUD) GetUserById(id string) (*User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res User
	result := db.First(&res, id)
	return &res, result.Error
}

func (crud UserCRUD) FuzzyGetUserByName(name string) ([]User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var Users []User
	result := db.Where("name LIKE ? AND IsDeleted = ?", "%"+name+"%", false).Find(&Users)
	return Users, result.Error
}

func (crud UserCRUD) GetUserByName(name string) (*User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res User
	result := db.First(&res, name)
	return &res, result.Error
}

func (crud UserCRUD) UpdateByObject(u *User) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	result := db.Save(&u)
	return result.Error
}

func (crud UserCRUD) GetAllUser() ([]User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var res []User
	result := db.Find(&res)
	return res, result.Error

}

func (crud UserCRUD) GetAllUserOrdered() ([]User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}

	var top []User
	var res []User
	topUsers := db.Where("is_top = ?", true).Find(&top)
	if topUsers.Error != nil {
		return nil, topUsers.Error
	}

	result := db.Order("updated_at desc").Find(&res)
	res = append(res, top...)
	return res, result.Error
}

func (crud UserCRUD) DeleteUserbyName(name string) error {
	result, err := crud.GetUserByName(name)
	if err != nil {
		return err
	}
	result.IsDeleted = true

	err = crud.UpdateByObject(result)
	if err != nil {
		return err
	}
	return nil
}
