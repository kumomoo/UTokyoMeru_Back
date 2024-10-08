package db

import (
	"backend/internal/model"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
)

type UsersCRUD struct{}

const secret = "whatcanisay"

func (crud UsersCRUD) CreateByObject(u model.User) error {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}
	u.Password = encryptPassword(u.Password)
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

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
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
	password := encryptPassword(oPassword)
	if password != user.Password {
		return errors.New("invalid password")
	}
	return
}

func (crud UsersCRUD) LoginByCode(user *model.User) (err error) {
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	err = db.Where("mail_address = ?", user.MailAddress).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not exist")
		}
		return err // 查询数据库时出错
	}
	return
}

func (crud UsersCRUD) ResetPassword(u model.User) error {
	// 获取数据库实例
	db, err := GetDatabaseInstance()
	if err != nil {
		return err
	}

	// 查找并更新指定的字段
	result := db.Model(&model.User{}).Where("mail_address = ?", u.MailAddress).Update("password", encryptPassword(u.Password))
	if result.Error != nil {
		return result.Error
	}

	// 检查是否有记录被更新
	if result.RowsAffected == 0 {
		return errors.New("no user found with the specified mail")
	}

	return nil
}
