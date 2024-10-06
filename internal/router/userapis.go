package router

import (
	"backend/internal/logic"
	"backend/internal/model"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p model.ParamSignup
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(500, gin.H{"message": "Invalid param", "error": err})
		return
	}

	fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(&p); err != nil {
		if strings.Contains(err.Error(), "23505") {
			c.JSON(500, gin.H{"message": "User exist", "error": err})
			return
		}
		c.JSON(500, gin.H{"message": "Server busy", "error": err})
		return
	}

	//3.返回响应
	c.JSON(200, nil)
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
