package middlewares

import (
	"exchangeapp/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMidlleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Authorization 是 HTTP 协议中定义的标准头字段，专门用于传递认证信息。使用标准字段可以确保代码的可读性和兼容性。
		token := ctx.GetHeader("Authorization")

		//判断token是否存在
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Authorization header",
			})
			//Abort()方法会停止当前请求的处理链,并组织后续的处理函数或者中间件继续执行
			//调用abort之后GIN会立即停止后续的中间件或者处理函数的执行
			ctx.Abort()
			return
		}
		username, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			ctx.Abort()
			return
		}
		ctx.Set("username", username)
	}
}
