package middlewares

import (
	"backend/internal/db"
	"backend/internal/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxxx.xxx.xxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{"message": "Need login"})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(403, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"message": "Token expired"})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(ContextUserIDKey, mc.MailAddress)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}



func AdminAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(403, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"message": "Token expired"})
			c.Abort()
			return
		}
		crud := db.UsersCRUD{}
		user, err := crud.FindOneByUniqueField("mail_address", mc.MailAddress)
		if err != nil {
			c.JSON(404, gin.H{"message": "Admin not exist", "error": err})
			return
		}
		if user.UserClass == "admin" {
			c.Next()
		} else {
			c.JSON(403, gin.H{"message": "Need user role: admin"})
			c.Abort()
		}
	}
}
