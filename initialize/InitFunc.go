package initialize

//这是初始化文件

import (
	"encoding/json"
	"fmt"
	"gitee.com/yumo01/bag/global"
	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化zap日志服务
func Zap() {
	////创建logger实例
	Logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()

	//设置全局 Logger
	zap.ReplaceGlobals(Logger)

	zap.S().Info("初始化日志成功")
}

// 初始化配置文件读取
func Viper() {
	//实例化一个viper方法
	v := viper.New()

	//自动读取配置文件
	//viper.AutomaticEnv()

	//手动设置读取的文件路径
	v.SetConfigFile("./conf/config.yaml")

	//启用配置文件的动态监视,配置文件发生变化时自动重新加载配置
	v.WatchConfig()

	//读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Panic("读取配置文件失败")
		return
	}
	//把读取的配置文件信息拿出来
	err = v.Unmarshal(&global.SevConf)
	if err != nil {
		zap.S().Panic("解析yaml配置文件失败")
		return
	}

	//若配置文件发生了变化
	v.OnConfigChange(func(c fsnotify.Event) {
		//把读取的配置文件信息拿出来
		err = v.Unmarshal(&global.SevConf)
		if err != nil {
			zap.S().Panic("解析yaml配置文件失败")
			return
		}
		zap.S().Info("rpc配置发生变动")
		Mysql()
	})
	zap.S().Info("viper初始化完成")

}

// 初始化配置文件远程读取
func Nacos() {

	//服务端
	sc := []constant.ServerConfig{
		{
			IpAddr: global.SevConf.NacosConfig.ServerConfig.IpAddr,
			Port:   global.SevConf.NacosConfig.ServerConfig.Port,
		},
	}

	//客户端
	cc := constant.ClientConfig{
		NamespaceId: global.SevConf.NacosConfig.ClientConfig.NamespaceId,
		//NotLoadCacheAtStart: true,
		LogDir:   "tmp/log",
		CacheDir: "tmp/cache",
		LogLevel: "debug",
	}

	//创建客户端动态配置
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		zap.S().Panic(err)
	}

	//获取读取的
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.SevConf.NacosConfig.ConfigParam.DataId,
		Group:  global.SevConf.NacosConfig.ConfigParam.Group})

	err = json.Unmarshal([]byte(content), &global.SevConf)
	if err != nil {
		zap.S().Panic(err)
		return
	}
	zap.S().Info("nacos初始化完成")

}

// 数据库初始化
func Mysql() {
	var err error
	global.MysqlDB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", global.SevConf.Mysql.Root, global.SevConf.Mysql.Password, global.SevConf.Mysql.Host, global.SevConf.Mysql.Port, global.SevConf.Mysql.Database)), &gorm.Config{})
	if err != nil {
		zap.S().Panic("数据库连接失败", err)
		return
	}

	zap.S().Info("数据库初始化完成")
}

// consul健康检测服务
func Consul() {
	//创建默认客户端配置
	conf := api.DefaultConfig()

	//配置consul路由
	conf.Address = "43.143.123.142:8500" //默认

	//实例化客户端
	client, err := api.NewClient(conf)
	if err != nil {
		zap.S().Panic(err)
		return
	}

	//定义健康检查服务
	check := &api.AgentServiceCheck{
		GRPC:                           "10.2.171.69:8080",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//srvId:=fmt.Sprintf("%s",uuid.NewV4())
	//用于注册服务的结构体
	reg := api.AgentServiceRegistration{
		Address: "10.2.171.69",
		Port:    8080,
		Name:    "first_srv",
		//Tags:    []string{"tag1"},
		ID:    "first_id", //若不填默认是Name
		Check: check,
	}

	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		zap.S().Panic(err)
		return
	}

	zap.S().Info("consul注册服务成功")
}
