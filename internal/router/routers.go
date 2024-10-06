package router

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {

	v1 := Router.Group("/v1")

	v1.POST("/signup", SignUpHandler)
	v1.POST("/login", LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		v1.GET("/goods/", GetAllGoods)
		v1.POST("/goods/", CreateGood)
		v1.GET("/goods/:id", GetGoodById)
		v1.PUT("/goods/:id", UpdateGood)
		v1.DELETE("/goods/:id", DeleteGood)
	}

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
}
