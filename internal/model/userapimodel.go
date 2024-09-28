package model

type GetUser struct {
	Name        string
	MailAddress string
	Password    string
	Avatar      string
	IsDeleted   bool
	IsBanned    bool
	UserClass   string
	Gender      string
	PhoneNumber string
	MailCode    string
	Address     string
}