package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	/*
		在 Go 语言中，{ 是代码块的起始符号，通常用于函数、方法或控制结构（如 if、for 等）。
		但在这里，r.Group("api/auth") 返回的是一个 *gin.RouterGroup 对象，
		而不是一个需要代码块的控制结构。
		因此，{ 的使用是为了包裹一组路由定义，而不是 Go 语言语法本身的要求。这种写法是开发者的风格选择，
		目的是让代码更清晰地表达路由组的范围。
	*/
	auth := r.Group("api/auth")
	{ //表示这个路由组的范围从这里开始,如果不加的话也行,但是不清晰
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)

	}
	api := r.Group("/api")
	api.GET("exchangeRates", controllers.GetExchangeRate)
	//使用中间件
	api.Use(middlewares.AuthMidlleware())
	{
		//创建对应汇率
		api.POST("/exchangeRates", controllers.CreateExchangeRate)

		//获取文章
		api.GET("/articles", controllers.GetArticle)

		//通过文章id获取文章
		api.GET("/articles/:id", controllers.GetArticleByID)

		//获取点赞数目
		api.GET("articles/:id/likes", controllers.GetArticleLikes)

		//给文章点赞
		api.POST("articles/:id/likes", controllers.LikeArticle)

		//取消文章点赞
		api.POST("articles/:id/unlikes", controllers.UnLikeArticle)
	}

	return r
}
