package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

//存放一些全局变量

var (
	Db      *gorm.DB
	RedisDb *redis.Client
)
