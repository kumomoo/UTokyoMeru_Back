package router

import (
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllGoods(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	usercrud := &db.UsersCRUD{}

	// 获取请求参数
	pageNumStr := c.Query("PageNum")
	itemNumStr := c.Query("ItemNum")

	// 检查参数是否存在
	var useLimit bool
	var limit int
	if pageNumStr != "" && itemNumStr != "" {
		// 转换参数为整数
		pageNum, err := strconv.Atoi(pageNumStr)
		if err != nil || pageNum < 1 {
			pageNum = 1 // 如果参数无效，设置默认页码为1
		}

		itemNum, err := strconv.Atoi(itemNumStr)
		if err != nil || itemNum < 1 {
			itemNum = 10 // 如果参数无效，设置每页默认数量为10
		}

		// 计算需要获取的商品数量
		limit = pageNum * itemNum
		useLimit = true
	}

	// 查找所有商品并按顺序排列
	result, err := crud.FindAllOrdered()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot Find Goods", "error": err})
		return
	}

	// 处理需要获取的商品数量
	var limitedResult []model.Good
	if useLimit {
		// 确保 limit 不超过商品总数
		if limit > len(result) {
			limit = len(result)
		}
		limitedResult = result[:limit]
	} else {
		// 如果没有分页参数，返回所有商品
		limitedResult = result
	}

	// 构建响应
	posts := make([]model.GetGoodsResponse, len(limitedResult))
	for i := range limitedResult {
		theUser, err := usercrud.FindById(limitedResult[i].SellerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot Find User", "error": err})
			return
		}
		posts[i] = gt.FindGoodsByIdDb2ResponseModel(limitedResult[i], *theUser)
	}

	c.JSON(http.StatusOK, posts)
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

func CreateGoodHandler(c *gin.Context) {
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

func UpdateGoodHandler(c *gin.Context) {
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

func DeleteGoodHandler(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}
	crud.DeleteById(uint(id))
	c.JSON(200, gin.H{"message": "Good Deleted"})
}

func LikeGoodHandler(c *gin.Context) {
	u := &db.UsersCRUD{}
	uid, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid user ID", "error": err})
		return
	}
	gid, err := strconv.Atoi(c.Query("goodID"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid good ID", "error": err})
		return
	}
	err = u.AddFavorite(uint(uid), uint(gid))
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot Like Good", "error": err})
		return
	}
	c.JSON(200, gin.H{"message": "Good Liked"})
}

func UnLikeGoodHandler(c *gin.Context) {
	u := &db.UsersCRUD{}
	uid, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid user ID", "error": err})
		return
	}
	gid, err := strconv.Atoi(c.Query("goodID"))
	if err != nil {
		c.JSON(500, gin.H{"message": "Invalid good ID", "error": err})
		return
	}
	err = u.RemoveFavorite(uint(uid), uint(gid))
	if err != nil {
		c.JSON(500, gin.H{"message": "Cannot UnLike Good", "error": err})
		return
	}
}