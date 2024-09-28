package router

import (

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init(){
	goods:=Router.Group("/goods")
	{
		goods.GET("/",GetAllGoods)
		goods.POST("/",CreateGood)
		goods.GET("/:id",GetGoodById)
		goods.PUT("/:id",UpdateGood)
		goods.DELETE("/:id",DeleteGood)
	}
}