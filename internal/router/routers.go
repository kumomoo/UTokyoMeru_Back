package router

import (
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {

	Router.POST("/signup", SignUpHandler)
	Router.POST("/login/verification", VerificationHandler)
	Router.POST("/login/password", LoginHandler)
	Router.POST("/login/code", LoginByCodeHandler)
	Router.POST("/login/reset", ResetPasswordHandler)

	goodsAuth := Router.Group("/goods")
	goodsAuth.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		goodsAuth.POST("/", CreateGoodHandler)
		goodsAuth.PUT("/:id", UpdateGoodHandler)
		goodsAuth.DELETE("/:id", DeleteGoodHandler)
		goodsAuth.POST("/like", LikeGoodHandler)
		goodsAuth.DELETE("/like", UnLikeGoodHandler)
	}

	goodsUnauth := Router.Group("/goods")
	{
		goodsUnauth.GET("/", GetAllGoods)
		goodsUnauth.GET("/:id", GetGoodById)
		goodsUnauth.GET("/search", SearchGoodsHandler)
	}

	admin := Router.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.GET("/users", GetAllUsersHandler)
		admin.GET("/users/:id", GetUserInfoByIdHandler)
		admin.PUT("/users/ban", BanUserHandler)
		admin.PUT("/users/unban", UnbanUserHandler)
		admin.PUT("/users/", UpdateUserHandler)
	}

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "API not found",
		})
	})
}
