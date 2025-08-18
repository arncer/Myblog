package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var cacheKey = "articles"

// CreateArticle 创建文章
func CreateArticle(ctx *gin.Context) {
	//模型绑定,绑定到json格式上
	var article models.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//新增文章时,同步更新缓存
	if err := global.RedisDb.Del(cacheKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusCreated, article)
}

// GetArticle 获取数据库中的所有文章
func GetArticle(ctx *gin.Context) {
	cacheData, err := global.RedisDb.Get(cacheKey).Result()
	//如果缓存中没有数据,则从数据库中获取
	if err == redis.Nil {
		var articles []models.Article
		if err := global.Db.Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//序列化为json格式的数组
		artcicleJson, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		//将文章列表存入缓存中,设置过期时间为10分钟
		if err := global.RedisDb.Set(cacheKey, artcicleJson, 10*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		//如果缓存有数据,直接从缓存中获取数据,直接将缓存的数据反序列化文章列表,返回给客户端

		var articles []models.Article
		if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//如果反序列化成功,则直接返回文章列表
		ctx.JSON(http.StatusOK, articles)
	}
}

// GetArticleByID 获取指定ID的文章
func GetArticleByID(ctx *gin.Context) {

	//获取路由中的参数,比如 /articles/:id
	id := ctx.Param("id")
	//找到对应的文章
	var article models.Article
	//判断错误类型
	if err := global.Db.Where("id = ?", id).Find(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	//如果没有问题,返回文章
	ctx.JSON(http.StatusOK, article)
}
