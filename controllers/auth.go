package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		//ctx.ShouldBindJSON(&user) 的作用是从 HTTP 请求的 JSON 数据中解析并绑定到 user 变量中。
		//它会根据 user 结构体的字段定义，
		//将请求体中的 JSON 数据自动映射到对应的字段。  如果解析成功，user 变量会被填充；
		//如果解析失败（例如 JSON 格式错误或字段类型不匹配），
		//会返回一个错误，供后续处理。  这是 Gin 框架中用于处理 JSON 请求体的常用方法之一。
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}
	//加密密码
	hashPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		//statusInternalServerError 表示服务器内部错误，通常用于处理服务器无法完成请求的情况。
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//将加密后的密码存入数据库中
	user.Password = hashPwd

	//AutoMigrate 的作用是根据传入的模型（如 User）自动创建或更新数据库表结构
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//将用户信息存储到数据库中
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//生成token
	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	//如果没有问题
	ctx.JSON(http.StatusOK, gin.H{"Token": token})
}

// 登录
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//反序列化
	if err := ctx.ShouldBindJSON(&input); err != nil { //JSON 数据解析并绑定到 Go 语言结构体
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	//查询用户是否在数据库
	if err := global.Db.Where("username=?", input.Username).First(&user).Error; err != nil {
		//HTTP 401 状态码表示 "Unauthorized"，即未经授权。它通常用于指示客户端请求需要身份验证，但未提供有效的凭据，或者提供的凭据无效。
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	//判断密码是否正确
	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "wrong credential"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
