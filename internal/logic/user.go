package logic

import (
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/pkg/jwt"
)

func SignUp(p *model.ParamSignup) (user *model.User, err error) {
	crud := &db.UsersCRUD{}
	//构造一个user实例
	user = &model.User{
		Name:        p.Username,
		MailAddress: p.MailAddress,
		Password:    p.Password,
		Gender:      p.Gender,
		PhoneNumber: p.PhoneNumber,
		Address:     p.Address.PostalCode + p.Address.Prefecture + p.Address.City + p.Address.AddressDetail,
	}
	//3.保存进数据库
	if err := crud.CreateByObject(*user); err != nil {
		return nil, err
	}
	token, err := jwt.GenToken(user.MailAddress, user.Name)
	if err != nil {
		return
	}
	user.Token = token
	return
}

func Login(p *model.ParamLogin) (user *model.User, err error) {
	crud := &db.UsersCRUD{}
	user = &model.User{
		MailAddress: p.MailAddress,
		Password:    p.Password,
	}
	if err := crud.Login(user); err != nil {
		return nil, err
	}
	//生成jwt
	token, err := jwt.GenToken(user.MailAddress, user.Name)
	if err != nil {
		return
	}
	user.Token = token
	return
}

func LoginByCode(p *model.ParamLoginByCode) (user *model.User, err error) {
	crud := &db.UsersCRUD{}
	user = &model.User{
		MailAddress: p.MailAddress,
	}
	if err := crud.LoginByCode(user); err != nil {
		return nil, err
	}
	//生成jwt
	token, err := jwt.GenToken(user.MailAddress, user.Name)
	if err != nil {
		return
	}
	user.Token = token
	return
}

func ResetPassword(p *model.ParamResetPassword) (err error) {
	crud := &db.UsersCRUD{}
	user := &model.User{
		MailAddress: p.MailAddress,
		Password:    p.Password,
	}
	if err := crud.ResetPassword(*user); err != nil {
		return err
	}
	return
}
