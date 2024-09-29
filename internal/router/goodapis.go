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
	result, err := crud.FindAllOrdered()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Goods"})
		return
	}
	posts := make([]model.GetGoodsResponse, len(result))
	for i := range posts {
		posts[i] = gt.Db2ResponseModel(result[i])
	}
	c.JSON(200, posts)
}

func GetGoodById(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := crud.FindById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Good"})
		return
	}
	post := gt.Db2ResponseModel(*result)
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
		c.JSON(500, gin.H{"error": "Cannot Create Good"})
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
		c.JSON(500, gin.H{"error": "Invalid ID"})
		return
	}

	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}

	dbGood, err := crud.FindById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Post"})
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
		c.JSON(500, gin.H{"error": "Cannot Update Post"})
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
		c.JSON(500, gin.H{"error": "Invalid ID"})
		return
	}
	crud.DeleteById(uint(id))
	c.JSON(200, gin.H{"message": "Good Deleted"})
}
