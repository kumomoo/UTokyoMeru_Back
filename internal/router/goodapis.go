package router

import (
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllGoods(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	result, err := crud.FindAllOrdered()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Posts"})
		return
	}
	posts := make([]model.GetGoods, len(result))
	for i := range posts {
		posts[i] = gt.GoodTransformToApiModel(result[i])
	}
	c.JSON(200, posts)
	return
}

func GetGoodsById(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := crud.FindById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Post"})
		return
	}
	post := gt.GoodTransformToApiModel(*result)
	c.JSON(200, post)
	return
}

func CreateGood(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}
	dbGood := gt.GoodTransformToDbModel(good)
	err := crud.CreateByObject(dbGood)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Create Post"})
		return
	}
	post := model.PostGoodsResponse{
		Success: true,
		Message: "Good Created",
		GoodInfo: dbGood,
	}
	c.JSON(200, post)
	return
}
