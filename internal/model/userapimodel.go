package model

type GetUserResponse struct {
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

type PostUserReceive struct{
	Name        string
	MailAddress string
	Password    string
	Avatar      string
	UserClass   string
	Gender      string
	PhoneNumber string
	MailCode    string
	Address     string
}