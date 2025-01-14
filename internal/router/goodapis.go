package router

import (
	"backend/internal/db"
	"backend/internal/model"
	"backend/internal/utils"
	"backend/internal/utils/logger"
	
	"net/http"
	"strconv"

	"go.uber.org/zap"
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
	result, err := crud.FindAllOrdered("updated_at", db.OrderDesc)
	if err != nil {
		logger.Logger.Error("无法找到商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
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
			logger.Logger.Error("无法找到用户",
				zap.String("path", c.FullPath()),
				zap.Any("params", c.Params),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot Find User", "error": err})
			return
		}
		posts[i] = gt.FindGoodsByIdDb2ResponseModel(limitedResult[i], *theUser)
	}
	logger.Logger.Debug("构建响应成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
		zap.Any("posts", posts),
	)

	c.JSON(http.StatusOK, posts)
}

func GetGoodById(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}

	logger.Logger.Debug("开始获取商品详情",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("无效ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}
	result, err := crud.FindById(uint(id))
	if err != nil {
		logger.Logger.Error("无法找到商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot Find Good", "error": err})
		return
	}
	usercrud := &db.UsersCRUD{}
	user, err := usercrud.FindById(result.SellerID)
	if err != nil {
		logger.Logger.Error("无法找到用户",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot Find User", "error": err})
		return
	}
	post := gt.FindGoodsByIdDb2ResponseModel(*result, *user)
	

	// 增加点击量
	err = crud.UpdateByField("Views", result.Views+1, *result)
	if err != nil {
		logger.Logger.Error("无法更新点击量",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
	}
	logger.Logger.Debug("获取商品详情成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, post)
}

func CreateGoodHandler(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		logger.Logger.Error("无效输入",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Any("body", c.Request.Body),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"error": "Invalid Input"})
		return
	}
	dbGood := gt.Post2DbModel(good)
	id, err := crud.CreateByObject(&dbGood)
	if err != nil {
		logger.Logger.Error("无法创建商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Any("good", good),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot Create Good", "error": err})
		return
	}
	res, err := crud.FindById(id)
	if err != nil {
		logger.Logger.Error("无法找到商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Database Error", "error": err})
		return
	}
	post := model.PostGoodsResponse{
		Message:  "Good Created",
		GoodInfo: *res,
	}
	logger.Logger.Debug("创建商品成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
		zap.Any("post", post),
	)
	c.JSON(200, post)
}

func UpdateGoodHandler(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("无效ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}

	var good model.PostGoodsReceive
	if err := c.ShouldBindJSON(&good); err != nil {
		logger.Logger.Error("无效输入",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Any("body", c.Request.Body),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid Input", "error": err})
		return
	}

	dbGood, err := crud.FindById(uint(id))
	if err != nil {
		logger.Logger.Error("无法找到商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
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
		logger.Logger.Error("无法更新商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Any("good", dbGood),
			zap.Error(err),
		)
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
		logger.Logger.Error("无效ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid ID", "error": err})
		return
	}
	crud.DeleteById(uint(id))
	logger.Logger.Debug("删除商品成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{"message": "Good Deleted"})
}

func LikeGoodHandler(c *gin.Context) {
	u := &db.UsersCRUD{}
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		logger.Logger.Error("无效用户ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid user ID. Please enter a number", "error": err})
		return
	}
	gid, err := strconv.Atoi(c.Query("good_id"))
	if err != nil {
		logger.Logger.Error("无效商品ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid good ID. Please enter a number", "error": err})
		return
	}
	err = u.AddFavorite(uint(uid), uint(gid))
	if err != nil {
		logger.Logger.Error("无法收藏",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot Like Good", "error": err})
		return
	}
	logger.Logger.Debug("收藏成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{"message": "Good Liked"})
}

func UnLikeGoodHandler(c *gin.Context) {
	u := &db.UsersCRUD{}
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		logger.Logger.Error("无效用户ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid user ID. Please enter a number", "error": err})
		return
	}
	gid, err := strconv.Atoi(c.Query("good_id"))
	if err != nil {
		logger.Logger.Error("无效商品ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid good ID. Please enter a number", "error": err})
		return
	}
	err = u.RemoveFavorite(uint(uid), uint(gid))
	if err != nil {
		logger.Logger.Error("无法取消收藏",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot UnLike Good", "error": err})
		return
	}
	logger.Logger.Debug("取消收藏成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{"message": "Good Unliked"})
}

func SearchGoodsHandler(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	gt := &utils.GoodTransform{}
	keyword := c.Query("keyword")
	orderBy := c.Query("orderBy")
	order := c.Query("order")

	goods, err := crud.Search(
		db.WithKeyword(keyword),
		db.WithOrderBy(orderBy),
		db.WithOrder(order),
	)
	if err != nil {
		logger.Logger.Error("无法搜索商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Cannot Search Good", "error": err})
		return
	}
	response := make([]model.GetGoodsResponse, len(goods))
	for i := range goods {
		response[i] = gt.FindGoodsByIdDb2ResponseModel(goods[i], goods[i].Seller)
	}
	logger.Logger.Debug("搜索商品成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
		zap.Any("response", response),
	)
	c.JSON(200, response)
}

func BuyGoodHandler(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		logger.Logger.Error("无效用户ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid user ID. Please enter a number", "error": err})
		return
	}
	gid, err := strconv.Atoi(c.Query("good_id"))
	if err != nil {
		logger.Logger.Error("无效商品ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid good ID. Please enter a number", "error": err})
		return
	}
	
	good, err := crud.FindById(uint(gid))
	if err != nil {
		logger.Logger.Error("无法找到商品",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Cannot Find Good", "error": err})
		return
	}
	if good.IsBought {
		logger.Logger.Error("商品已购买",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"message": "Good has already been bought"})
		return
	}
	good.IsBought = true
	good.BuyerID = uint(uid)
	crud.UpdateByObject(*good)
	logger.Logger.Debug("购买商品成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{"message": "Good Bought"})
}
