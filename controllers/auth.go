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

	hashPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		//statusInternalServerError 表示服务器内部错误，通常用于处理服务器无法完成请求的情况。
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = hashPwd

	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	//如果没有问题
	ctx.JSON(http.StatusOK, gin.H{"Token": token})

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
}

//登录

func login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
