package main

import (
	_ "backend/internal/db"
	"backend/internal/middlewares"
	"backend/internal/router"
	"backend/internal/utils/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	logger.InitLogger()
	defer logger.Logger.Sync()

	gin.SetMode(gin.DebugMode)

	r := router.Router
	r.Use(middlewares.CORSMiddleware())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	err := r.Run(":8100")
	if err != nil {
		fmt.Println(err)
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
