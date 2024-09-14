package router

import (
	"backend/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllGoods(c *gin.Context) {
	crud := &db.GoodsCRUD{}
	result, err := crud.FindAllOrdered()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot Find Posts"})
		return
	}
	posts := make([]db.GoodsGet, len(result))
	for i := range posts {
		posts[i] = db.GoodsGet{
			GoodsId:     result[i].GoodsId,
			CreatedTime: result[i].CreatedAt,
			UpdatedTime: result[i].UpdatedAt,
			GoodsName:   result[i].GoodsName,
			Describe:    result[i].Describe,
			AuthorId:    result[i].AuthorId,
			AuthorName:  result[i].AuthorName,
			Images:      result[i].Images,
			Price:       result[i].Price,
			Location:    result[i].Location,
			Avatar:      result[i].AuthorAvatar,
			Type:        result[i].Type,
			Views:       result[i].Views,
		}
	}
	c.JSON(200, posts)
	return
}

func GetAllComments(c *gin.Context) {
	crud := &db.CommentCRUD{}
	postId, err := strconv.Atoi(c.Param("PostId"))
	if err != nil {
		c.JSON(500, gin.H{"error": "postId is invalid"})
		return
	}
	if postId <= 0 {
		c.JSON(500, gin.H{"error": "postId is invalid"})
		return
	}

	result, err := crud.FindAllByGoodsId(uint(postId))
	if err != nil {
		c.JSON(404, gin.H{"error": "post not found"})
		return
	}
	comments := make([]db.CommentGet, len(result))
	for i := range comments {
		comments[i] = db.CommentGet{
			ID:          result[i].ID,
			CreatedTime: result[i].CreatedAt,
			UpdatedTime: result[i].UpdatedAt,
			GoodsId:     result[i].GoodsId,
			GoodsName:   result[i].GoodsName,
			AuthorId:    result[i].AuthorId,
			AuthorName:  result[i].AuthorName,
			Avatar:      result[i].Avatar,
			Content:     result[i].Content,
			ReplyTo:     result[i].ReplyTo,
		}
	}

	c.JSON(200, comments)
	return
}

func CreateGoods(c *gin.Context) {
	g := db.GoodsCRUD{}
	var goods db.GoodsRequest
	err := c.BindJSON(&goods)

	if err != nil {
		c.JSON(500, gin.H{"error": "Format of request body is invalid"})
		return
	}

	newGoods := &db.Goods{
		GoodsName: goods.GoodsName,
	}
	err = g.CreateByObject(newGoods)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot create post now"})
		return
	}

	c.JSON(200, gin.H{"success": true})
	return
}

func CreateReply(c *gin.Context) {
	r := db.ReplyCRUD{}
	var reply db.ReplyRequest
	err := c.BindJSON(&reply)
	if err != nil {
		c.JSON(500, gin.H{"error": "Format of request body is invalid"})
		return
	}
	newReply := &db.Reply{
		PostId:   reply.PostId,
		AuthorId: reply.AuthorId,
		Content:  reply.Content,
		ReplyTo:  reply.ReplyTo,
	}

	err = r.CreateByObject(newReply)
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot create reply"})
		return
	}
	c.JSON(200, gin.H{"success": true})
	return
}
