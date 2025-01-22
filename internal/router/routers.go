package router

import (
	"backend/internal/middlewares"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	// 创建路由实例
	Router = gin.New()

	// 配置全局中间件
	Router.Use(gin.Recovery())               // 恢复中间件
	Router.Use(middlewares.CORSMiddleware()) // CORS中间件

	// 配置日志中间件
	Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// 路由组配置
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
		userAuth.GET("/data", GetAllGoodsStatsHandler)
		userAuth.GET("/bought", GetAllBoughtGoodsHandler)
		userAuth.GET("/sold", GetAllSoldGoodsHandler)
		userAuth.GET("/:user_id", GetUserCommonInfoHandler)
	}

	userUnauth := Router.Group("/user")
	{
		userUnauth.GET("/sales", GetAllSalesGoodsHandler)
		userUnauth.GET("/selling", GetAllSellingGoodsHandler)
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
