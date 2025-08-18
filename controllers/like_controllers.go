package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// LikeArticle 处理文章点赞请求
func LikeArticle(ctx *gin.Context) {
	articleId := ctx.Param("id")

	//拼接键值对
	likeKey := "article:" + articleId + ":likes"

	//增加文章的点赞数
	if err := global.RedisDb.Incr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successful"})
}

// GetArticleLikes 获取文章点赞数目
func GetArticleLikes(ctx *gin.Context) {
	articleId := ctx.Param("id")
	//拼接key
	likeKey := "article:" + articleId + ":likes"
	likes, err := global.RedisDb.Get(likeKey).Result()
	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//成功则返回点赞数目
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}

// UnLikeArticle 取消当前点赞
func UnLikeArticle(ctx *gin.Context) {
	//获取文章id
	articleId := ctx.Param("id")

	//拼接key
	likeKey := "article:" + articleId + ":likes"

	//减少文章点赞数目
	if err := global.RedisDb.Decr(likeKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Unliked successfully"})

}
