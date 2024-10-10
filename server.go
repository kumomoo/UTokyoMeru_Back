package main

import (
	_ "backend/internal/db"
	"backend/internal/middlewares"
	"backend/internal/router"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	r := router.Router
	r.Use(middlewares.CORSMiddleware())

	err := r.Run(":8101")
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
