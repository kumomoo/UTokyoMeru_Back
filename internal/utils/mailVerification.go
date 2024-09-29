package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

// GenerateRandomCode 生成指定长度的随机验证码
func GenerateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(to, code string) error {
	e := email.NewEmail()
	e.From = "yamanashiluna@gmail.com"
	e.To = []string{to}
	e.Subject = "邮箱验证码"
	e.Text = []byte(fmt.Sprintf("您的验证码是：%s", code))

	// 配置SMTP服务器信息
	smtpHost := "smtp.example.com"
	smtpPort := 587
	smtpUsername := "your-username"
	smtpPassword := "your-password"

	return e.Send(fmt.Sprintf("%s:%d", smtpHost, smtpPort), smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost))
}
