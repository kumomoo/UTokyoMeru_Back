package router

import (
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {

	Router.POST("/signup", SignUpHandler)
	Router.POST("/signup/verification", VerificationHandler)
	Router.POST("/login/verification", VerificationHandler)
	Router.POST("/login/password", LoginHandler)
	Router.POST("/login/code", LoginByCodeHandler)
	Router.POST("/login/reset", ResetPasswordHandler)

	goodsAuth := Router.Group("/goods")
	goodsAuth.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		goodsAuth.POST("/", CreateGood)
		goodsAuth.PUT("/:id", UpdateGood)
		goodsAuth.DELETE("/:id", DeleteGood)
	}

	goodsUnauth := Router.Group("/goods")
	{
		goodsUnauth.GET("/", GetAllGoods)
		goodsUnauth.GET("/:id", GetGoodById)
	}

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "API not foundg",
		})
	})
}
