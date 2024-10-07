package router

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {

	Router.POST("/signup", SignUpHandler)
	Router.POST("/verification", VerificationHandler)
	Router.POST("/login", LoginHandler)

	goods := Router.Group("/goods")
	goods.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		goods.GET("/", GetAllGoods)
		goods.POST("/", CreateGood)
		goods.GET("/:id", GetGoodById)
		goods.PUT("/:id", UpdateGood)
		goods.DELETE("/:id", DeleteGood)
	}

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
}
