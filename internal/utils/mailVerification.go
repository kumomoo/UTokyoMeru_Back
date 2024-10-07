package utils

import (
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
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
func SendEmail(to string, code string) error {
	// SMTP服务器配置
	host := "smtp.qq.com"           // 替换为您的SMTP服务器地址
	port := 465                     // SMTP服务器端口
	username := "1677575560@qq.com" // 发件人邮箱
	password := "msuqmezhglhbebch"  // 发件人邮箱密码

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/plain", fmt.Sprintf("您的验证码是：%s", code))

	d := gomail.NewDialer(host, port, username, password)
	d.SSL = true

	return d.DialAndSend(m)
}
