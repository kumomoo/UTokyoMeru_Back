package router

import (
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllGoods(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	usercrud := &db.UsersCRUD{}
	result, err := crud.FindAllOrdered()
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Find Goods", "error": err})
		return
	}
	posts := make([]model.GetGoodsResponse, len(result))
	for i := range posts {
		theUser, err := usercrud.FindById(result[i].SellerID)
		if err != nil {
			c.JSON(500, gin.H{"message": "Cannot Find User", "error": err})
			return
		}
		posts[i] = gt.FindGoodsByIdDb2ResponseModel(result[i], *theUser)
	}
	c.JSON(200, posts)
}

func GetGoodById(c *gin.Context) {
	// 数据库的CRUD
	crud := &db.GoodsCRUD{}
	// 数据转换
	gt := &utils.GoodTransform{}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}
	result, err := crud.FindById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Find Good", "error": err})
		return
	}
	usercrud := &db.UsersCRUD{}
	theUser, err := usercrud.FindById(result.SellerID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Find User", "error": err})
		return
	}
	post := gt.FindGoodsByIdDb2ResponseModel(*result, *theUser)

	c.JSON(200, post)
}

func CreateGood(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}
	dbGood := gt.Post2DbModel(good)
	err := crud.CreateByObject(dbGood)
	if err != nil {
		fmt.Println(err, dbGood)
		c.JSON(500, gin.H{"message": "Cannot Create Good", "error": err})
		return
	}
	post := model.PostGoodsResponse{
		Message:  "Good Created",
		GoodInfo: dbGood,
	}
	c.JSON(200, post)
}

func UpdateGood(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}

	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(400, gin.H{"message": "Invalid Input", "error": err})
		return
	}

	dbGood, err := crud.FindById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Find Good", "error": err})
		return
	}

	dbGood.Title = good.Title
	dbGood.Description = good.Description
	dbGood.Images = good.Images
	dbGood.Price = good.Price
	dbGood.Tags = good.Tags
	dbGood.IsInvisible = good.IsInvisible
	dbGood.IsDeleted = good.IsDeleted
	dbGood.IsBought = good.IsBought

	err = crud.UpdateByObject(*dbGood)
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Update Good", "error": err})
		return
	}
	updatedGood := model.PostGoodsResponse{
		Message:  "Good Updated",
		GoodInfo: *dbGood,
	}
	c.JSON(200, updatedGood)
}

func DeleteGood(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}
	crud.DeleteById(uint(id))
	c.JSON(200, gin.H{"message": "Good Deleted"})
}
