package db

import (
	"backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

type UsersCRUD struct{}

func (crud UsersCRUD) CreateByObject(u model.User) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Create(&u).Error
}

func (crud UsersCRUD) FindAll() ([]model.User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var users []model.User
	result := db.Find(&users)
	return users, result.Error
}

func (crud UsersCRUD) FindById(id uint) (*model.User, error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return nil, err
	}
	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (crud UsersCRUD) UpdateByObject(u model.User) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	return db.Save(&u).Error
}

func (crud UsersCRUD) DeleteById(id uint) error {
	obj, err := crud.FindById(id)
	if err != nil {
		return err
	}
	obj.IsDeleted = true
	return crud.UpdateByObject(*obj)
}

func (crud UsersCRUD) Login(user *model.User) (err error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	oPassword := user.Password // 用户输入的密码

	err = db.Where("mail_address = ?", user.MailAddress).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not exist")
		}
		return err // 查询数据库时出错
	}

	// 判断密码是否正确
	if oPassword != user.Password {
		return errors.New("invalid password")
	}
	return
}
