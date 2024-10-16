package router

import (
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {

	Router.POST("/signup", SignUpHandler)
	loginUnauth := Router.Group("/login")
	{
		loginUnauth.POST("/verification", VerificationHandler)
		loginUnauth.POST("/password", LoginHandler)
		loginUnauth.POST("/code", LoginByCodeHandler)
		loginUnauth.POST("/reset", ResetPasswordHandler)
	}

	goodsAuth := Router.Group("/goods")
	goodsAuth.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		goodsAuth.POST("/", CreateGoodHandler)
		goodsAuth.PUT("/:id", UpdateGoodHandler)
		goodsAuth.DELETE("/:id", DeleteGoodHandler)
		goodsAuth.PATCH("/like", LikeGoodHandler)
		goodsAuth.DELETE("/like", UnLikeGoodHandler)
		goodsAuth.POST("/buy", BuyGoodHandler)
	}

	goodsUnauth := Router.Group("/goods")
	{
		goodsUnauth.GET("/", GetAllGoods)
		goodsUnauth.GET("/:id", GetGoodById)
		goodsUnauth.GET("/search", SearchGoodsHandler)
	}

	userAuth := Router.Group("/user")
	userAuth.Use(middlewares.JWTAuthMiddleware())
	{
		userAuth.GET("/favolist", GetAllLikedGoodsHandler)
		userAuth.GET("/sales", GetAllSalesGoodsHandler)
		userAuth.GET("/data", GetAllUserDataHandler)
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
