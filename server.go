package main

import (
	_ "backend/internal/db"
	"backend/internal/utils/logger"
	"backend/internal/router"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	defer logger.Logger.Sync()

	// 设置 Gin 模式
	gin.SetMode(gin.DebugMode)

	// 启动服务器
	err := router.Router.Run(":8100")
	if err != nil {
		logger.Logger.Error("服务器启动失败", zap.Error(err))
	}
}

// func test(){
// 	testUser := model.User{
// 		Name:        "测试用户",
// 		MailAddress: "test@example.com",
// 		Password:    "password123",
// 		Avatar:      "https://picx.zhimg.com/v2-4a3fd89b61cfeb8f2fe7bc85de5b5438_1440w.jpg?source=7e7ef6e2&needBackground=1",
// 		UserClass:   "user",
// 		Gender:      "男",
// 		PhoneNumber: "12345678901",
// 		Address:     "测试地址",
// 	}

// 	var userCRUD db.UsersCRUD

// 	userCRUD.CreateByObject(testUser)
// }
