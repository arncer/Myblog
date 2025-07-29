package config

import (
	"exchangeapp/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func initDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil { //如果连接失败
		log.Fatalf("Failed to initialize the database,get a err:%v", err)
	}

	//用于从 *gorm.DB 实例中获取底层的 *sql.DB 对象。
	SqlDB, err := db.DB()

	//设置最大空连接数量为10,这些可以重复被使用,从而减少创建和销毁连接的开销;超出这个数量多余的空闲连接数量会被关闭
	SqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)

	//设置最大连接数为100,如果超过这个数量,则会阻塞等待
	SqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)

	//设置最大连接时间是1小时
	SqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configure the database, get a err:%v", err)
	}
	global.Db = db
}

//db *gorm.DB
