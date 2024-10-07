package router

import (
	"backend/internal/db"
	"backend/internal/logic"
	"backend/internal/model"
	"backend/internal/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func VerificationHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamVerify
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(500, gin.H{"message": "Invalid param", "error": err})
		return
	}
	to := p.MailAddress                 // 收件人邮箱
	code := utils.GenerateRandomCode(6) // 验证码

	err := utils.SendEmail(to, code)
	if err != nil {
		c.JSON(500, gin.H{"message": "sending code failed", "error": err})
	}

	//将验证码存入redis
	err = db.SetVerificationCode(p.MailAddress, code)
	if err != nil {
		c.JSON(500, gin.H{"message": "storaging code failed"})
		return
	}

	//返回响应
	c.JSON(200, gin.H{"message": "sending code success"})
}

func SignUpHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamSignup
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(500, gin.H{"message": "Invalid param", "error": err})
		return
	}
	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress)
	if err == redis.Nil {
		c.JSON(500, gin.H{"error": "Verificationcode expired or not exist."})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "getting verificationcode failed"})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		c.JSON(500, gin.H{"error": "VerificationCode error"})
		return
	}

	//业务处理
	user, err := logic.SignUp(&p)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			c.JSON(500, gin.H{"message": "User exist", "error": err})
			return
		}
		c.JSON(500, gin.H{"message": "Server busy", "error": err})
		return
	}

	//返回响应
	c.JSON(200, gin.H{
		"user_mailaddress": user.MailAddress,
		"user_name":        user.Name,
		"token":            user.Token,
	})
}

func LoginHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p model.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(500, gin.H{"message": "Invalid param", "error": err})
		return
	}
	fmt.Println(p)

	//2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			c.JSON(500, gin.H{"message": "User not exist", "error": err})
			return
		}
		c.JSON(500, gin.H{"message": "Invalid password", "error": err})
		return
	}

	//3.返回响应
	c.JSON(200, gin.H{
		"user_mailaddress": user.MailAddress,
		"user_name":        user.Name,
		"token":            user.Token,
	})
}
