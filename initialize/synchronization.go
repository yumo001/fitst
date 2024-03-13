package initialize

import (
	"context"
	user "github.com/yumo001/fitst/pb/user"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/yumo001/fitst/global"
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

// 同步用户表到es
func SynchronizationUser() {

	var users []user.User
	err := global.MysqlDB.Table("users").Find(&users).Error
	if err != nil {
		log.Panic("数据库查询失败", err)
		return
	}

	for _, u := range users {
		_, err = global.ElasticClient.Index().Index("users").Id(strconv.FormatInt(u.Id, 10)).BodyJson(&u).Do(context.Background())
		if err != nil {
			log.Panic("将数据填入es失败", err)
			return
		}
	}

	log.Println("es同步成功")
}
