package router

import (
	"strconv"
	"backend/internal/db"

	"github.com/gin-gonic/gin"
)


func GetUserInfoByIdHandler(c *gin.Context) {
	//获取参数与参数校验
	id := c.Param("id")
	crud := db.UsersCRUD{}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	user, err := crud.FindById(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"message": "User not exist", "error": err})
		return
	}
	response := user

	//返回响应
	c.JSON(200, response)
}

func GetAllUsersHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	users, err := crud.FindAll()
	if err != nil {
		c.JSON(404, gin.H{"message": "Users not exist", "error": err})
		return
	}
	response := users
	c.JSON(200, response)
}

func BanUserHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	crud := db.UsersCRUD{}
	user, err := crud.FindById(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"message": "User not exist", "error": err})
		return
	}
	user.IsBanned = true
	crud.UpdateByObject(*user)
	c.JSON(200, gin.H{"message": "Banned user"})
}

func UnbanUserHandler(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	crud := db.UsersCRUD{}
	user, err := crud.FindById(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"message": "User not exist", "error": err})
		return
	}
	user.IsBanned = false
	crud.UpdateByObject(*user)
	c.JSON(200, gin.H{"message": "Unbanned user"})
}
