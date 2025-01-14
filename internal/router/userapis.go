package router

import (
	"backend/internal/db"
	"backend/internal/logic"
	mw "backend/internal/middlewares"
	"backend/internal/model"
	"backend/internal/utils"
	"backend/internal/utils/logger"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func VerificationHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamVerify
	logger.Logger.Info("开始发送验证码",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		logger.Logger.Error("请求参数有误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	to := p.MailAddress                 // 收件人邮箱
	code := utils.GenerateRandomCode(6) // 验证码

	err := utils.SendEmail(to, code)
	if err != nil {
		logger.Logger.Error("发送验证码失败",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "sending code failed", "error": err})
	}

	//将验证码存入redis
	err = db.SetVerificationCode(p.MailAddress+p.VerificationCodeType, code)
	if err != nil {
		logger.Logger.Error("存储验证码失败",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "storaging code failed"})
		return
	}

	logger.Logger.Info("发送验证码成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
		zap.String("mailAddress", to),
		zap.String("code", code),
	)
	//返回响应
	c.JSON(200, gin.H{"message": "sending code success"})
}

func SignUpHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamSignup
	logger.Logger.Info("开始注册",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		logger.Logger.Error("请求参数有误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + p.VerificationCodeType)
	if err == redis.Nil {
		logger.Logger.Info("验证码过期或不存在",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"error": "Verificationcode expired or not exist."})
		return
	} else if err != nil {
		logger.Logger.Error("获取验证码失败",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"error": "getting verificationcode failed"})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		logger.Logger.Info("验证码错误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"error": "VerificationCode error"})
		return
	}

	//业务处理
	user, err := logic.SignUp(&p)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			logger.Logger.Info("用户已存在",
				zap.String("path", c.FullPath()),
				zap.Any("params", c.Params),
			)
			c.JSON(400, gin.H{"message": "User exist", "error": err})
			return
		}
		logger.Logger.Error("无法注册",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Server busy", "error": err})
		return
	}
	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}
	//返回响应
	logger.Logger.Info("注册成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func LoginHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p model.ParamLogin
	logger.Logger.Info("开始登录",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		logger.Logger.Error("请求参数有误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}

	//2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			logger.Logger.Info("用户不存在",
				zap.String("path", c.FullPath()),
				zap.Any("params", c.Params),
			)
			c.JSON(404, gin.H{"message": "User not exist", "error": err})
			return
		}
		logger.Logger.Info("密码错误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Invalid password", "error": err})
		return
	}

	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Token:       user.Token,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}

	//返回响应
	logger.Logger.Info("登录成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func LoginByCodeHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamLoginByCode
	logger.Logger.Info("开始登录",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		logger.Logger.Error("请求参数有误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}

	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + p.VerificationCodeType)
	if err == redis.Nil {
		logger.Logger.Info("验证码过期或不存在",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"message": "Verificationcode expired or not exist.", "error": err})
		return
	} else if err != nil {
		logger.Logger.Error("获取验证码失败",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Failed to get verification code", "error": err})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		logger.Logger.Info("验证码错误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"message": "Verification code error", "error": "VerificationCode error"})
		return
	}

	//业务处理
	user, err := logic.LoginByCode(&p)
	if err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			logger.Logger.Info("用户不存在",
				zap.String("path", c.FullPath()),
				zap.Any("params", c.Params),
			)
			c.JSON(400, gin.H{"message": "User not exist", "error": err})
			return
		}
		logger.Logger.Error("无法登录",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Internal error", "error": err})
		return
	}

	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Token:       user.Token,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}

	//返回响应
	logger.Logger.Info("登录成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func ResetPasswordHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamResetPassword
	logger.Logger.Info("开始重置密码",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		logger.Logger.Error("请求参数有误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}

	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + "reset")
	if err == redis.Nil {
		logger.Logger.Info("验证码过期或不存在",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(404, gin.H{"message": "Verificationcode expired or not exist.", "error": err})
		return
	} else if err != nil {
		logger.Logger.Error("获取验证码失败",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Failed to get verification code", "error": err})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		logger.Logger.Info("验证码错误",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
		)
		c.JSON(400, gin.H{"message": "VerificationCode error", "error": errors.New("VerificationCode error")})
		return
	}

	//业务处理
	if err := logic.ResetPassword(&p); err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			logger.Logger.Info("用户不存在",
				zap.String("path", c.FullPath()),
				zap.Any("params", c.Params),
			)
			c.JSON(404, gin.H{"message": "User not exist", "error": err})
			return
		}
		logger.Logger.Error("无法重置密码",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Invalid password", "error": err})
		return
	}

	//3.返回响应
	logger.Logger.Info("重置密码成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{"message": "reset password"})
}

// 需要userid因为可能需要看其他用户的出售商品
func GetAllSalesGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	logger.Logger.Error("开始获取用户所有出售商品信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		logger.Logger.Error("无法获取用户ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Failed to get params", "error": err})
		return
	}
	goods, err := crud.FindGoodsByFK(uint(userID), "Sales")
	if err != nil {
		logger.Logger.Error("无法获取信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Failed to get all sales goods", "error": err})
		return
	}
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	logger.Logger.Info("获取用户所有出售商品信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

// 需要userid因为可能需要看其他用户的
func GetAllSellingGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	logger.Logger.Error("开始获取用户正在出售商品信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		logger.Logger.Error("无法获取用户ID",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "Failed to get params", "error": err})
		return
	}
	goods, err := crud.FindGoodsByFK(uint(userID), "Sales")
	if err != nil {
		logger.Logger.Error("无法获取信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Failed to get all sales goods", "error": err})
		return
	}
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		if !good.IsBought {
			response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
		}
	}

	logger.Logger.Info("获取用户正在出售商品信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func GetAllLikedGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	logger.Logger.Error("开始获取用户所有收藏商品信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		logger.Logger.Error("无法获取用户信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.FavoList
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	logger.Logger.Info("获取用户所有收藏商品信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func GetAllBoughtGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	logger.Logger.Error("开始获取用户所有购买商品信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		logger.Logger.Error("无法获取用户信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.Buys
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	logger.Logger.Info("获取用户所有购买商品信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func GetAllSoldGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	logger.Logger.Error("开始获取用户所有已售商品信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		logger.Logger.Error("无法获取用户信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.Sales
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		if good.IsBought {
			response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
		}
	}

	logger.Logger.Info("获取用户所有已售商品信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, response)
}

func GetAllGoodsStatsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	//获取用户ID
	logger.Logger.Error("开始获取用户所有商品统计信息",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		logger.Logger.Error("无法获取用户信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}

	goods, err := crud.FindAllGoodsFK(user.ID)
	if err != nil {
		logger.Logger.Error("无法获取信息",
			zap.String("path", c.FullPath()),
			zap.Any("params", c.Params),
			zap.Error(err),
		)
		c.JSON(500, gin.H{"message": "Failed to get all goods", "error": err})
		return
	}

	var lengths [4]int
	for i := 0; i < 3; i++ {
		lengths[i] = len(goods[i])
	}

	isBoughtCount := 0
	for _, sale := range goods[0] {
		if sale.IsBought {
			isBoughtCount++
		}
	}
	lengths[3] = isBoughtCount

	logger.Logger.Info("获取用户所有商品统计信息成功",
		zap.String("path", c.FullPath()),
		zap.Any("params", c.Params),
	)
	c.JSON(200, gin.H{
		"sale_number":  lengths[0] - lengths[3],
		"sold_number":  lengths[3],
		"buy_number":   lengths[1],
		"favor_number": lengths[2],
	})
}
