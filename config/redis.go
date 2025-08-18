package config

import (
	"exchangeapp/global"
	"github.com/go-redis/redis"
	"log"
)

// InitRedis 初始化redis
func initRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		DB:       0,                // 使用默认DB
		Password: "111111",
	}) //options是填入的配置项

	//使用ping方法检查是否连接到redis服务
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis:%v", err)
	}
	//将RedisClient赋值给全局变量
	global.RedisDb = RedisClient
}
