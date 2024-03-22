package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/yumo001/fitst/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 创建es全局实例
func CreatElasticClient() {
	var err error
	//初始化es连接
	global.ElasticClient, err = elastic.NewClient(elastic.SetURL("http://"+global.SevConf.Elastic.Host+":"+global.SevConf.Elastic.Port), elastic.SetSniff(false))
	if err != nil {
		log.Panic("创建es实例失败", err)
		return
	}

	log.Println("elastic初始化成功")
}

// mysql闭包
func Mysql2(f func(mysqlDB *gorm.DB) error) {
	var err error
	mysqlDB, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", global.SevConf.Mysql.Root, global.SevConf.Mysql.Password, global.SevConf.Mysql.Host, global.SevConf.Mysql.Port, global.SevConf.Mysql.Database)), &gorm.Config{})
	if err != nil {
		zap.S().Panic("数据库连接失败", err)
		return
	}

	db, err := mysqlDB.DB()
	if err != nil {
		log.Panic("获取数据库连接对象失败", err)
		return
	}
	defer db.Close()

	f(mysqlDB)
}

// redis
func SendRedis() {
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     global.SevConf.Redis.Addr,
		Password: global.SevConf.Redis.Password,
		DB:       global.SevConf.Redis.DB,
	})
	defer global.RedisClient.Close()
}
