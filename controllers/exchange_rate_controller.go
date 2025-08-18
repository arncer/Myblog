package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CreateExchangeRate
func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//确定汇率的时间
	exchangeRate.Date = time.Now()
	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//创建成功
	ctx.JSON(http.StatusCreated, exchangeRate)

}

func GetExchangeRate(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate
	//Find() 方法本身不会直接返回 err，而是通过 gorm.DB 对象的 Error 字段来存储最近一次操作的错误。
	//这是 gorm 的设计模式，目的是让方法链调用更方便。
	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, exchangeRates)
}
