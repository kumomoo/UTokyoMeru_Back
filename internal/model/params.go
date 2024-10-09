package model

import (
	"time"
)

// 定义请求的参数结构体

type ParamVerify struct {
	MailAddress          string `json:"mail_address" binding:"required,email"`
	VerificationCodeType string `json:"verification_code_type" binding:"required"`
}
type ParamSignup struct {
	Username             string  `json:"user_name" binding:"required"`
	MailAddress          string  `json:"mail_address" binding:"required,email"`
	VerificationCode     string  `json:"verification_code" binding:"required"`
	Password             string  `json:"password" binding:"required"`
	Gender               string  `json:"gender" binding:"required,oneof=0 1 2 3"`
	Birthday             time.Time `json:"birthday"`
	PhoneNumber          string  `json:"phone_number"`
	Address              Address `json:"address"`
	VerificationCodeType string  `json:"verification_code_type" binding:"required"`
}

type Address struct {
	PostalCode    string `json:"postal_code"`
	Prefecture    string `json:"prefecture"`
	City          string `json:"city"`
	AddressDetail string `json:"address_detail"`
}

type ParamLogin struct {
	MailAddress string `json:"mail_address" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
}

type ParamLoginByCode struct {
	MailAddress          string `json:"mail_address" binding:"required,email"`
	VerificationCode     string `json:"verification_code" binding:"required"`
	VerificationCodeType string `json:"verification_code_type" binding:"required"`
}

type ParamResetPassword struct {
	MailAddress      string `json:"mail_address" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required"`
	Password         string `json:"password" binding:"required"`
}
