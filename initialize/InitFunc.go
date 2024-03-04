package initialize

//这是初始化文件

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"

	"github.com/yumo001/fitst/global"
)

// 初始化配置文件读取
func Viper() {
	//实例化一个viper方法
	v := viper.New()

	//自动读取配置文件
	//viper.AutomaticEnv()

	//手动设置读取的文件路径
	v.SetConfigFile("./conf/conf.yaml")

	//启用配置文件的动态监视,配置文件发生变化时自动重新加载配置
	v.WatchConfig()

	//读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件失败")
		return
	}
	//把读取的配置文件信息拿出来
	err = v.Unmarshal(&global.SevConf)
	if err != nil {
		log.Fatal("解析yaml配置文件失败")
		return
	}

	//若配置文件发生了变化
	v.OnConfigChange(func(c fsnotify.Event) {
		//把读取的配置文件信息拿出来
		err = v.Unmarshal(&global.SevConf)
		if err != nil {
			log.Fatal("解析yaml配置文件失败")
			return
		}
		log.Println("rpc配置发生变动")
		Nacos()
	})
	log.Println("viper初始化完成")

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
		log.Fatal(err)
	}

	//读取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.SevConf.NacosConfig.ConfigParam.DataId,
		Group:  global.SevConf.NacosConfig.ConfigParam.Group,
	})
	err = yaml.Unmarshal([]byte(content), &global.SevConf)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("nacos初始化完成")
	Mysql()

}

// nacos服务发现
func NacosServicesDiscovery(cc constant.ClientConfig, sc []constant.ServerConfig) {
	//注册服务
	nClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	port, err := strconv.ParseUint(global.SevConf.RpcPort, 10, 64)
	OK, err := nClient.RegisterInstance(vo.RegisterInstanceParam{
		// 实例ID，如果不指定，则由 Nacos 自动生成
		Ip: global.SevConf.NacosConfig.ClientConfig.NamespaceId,
		// 实例端口
		Port: port,
		// 指定权重，用于负载均衡，默认值为 1。
		Weight: 10,
		// 是否启用临时实例，默认为 true。
		Enable: true,
		// 指定实例的健康状态，默认为 true。
		Healthy: true,
		// 健康检查端点，用于健康检查。
		Metadata: map[string]string{},
		// 指定集群名称，如果不指定则为默认值。
		ClusterName: "Cluster1",
		// 指定要注册的服务名。
		ServiceName: global.SevConf.ServiceName,
		//组名
		GroupName: "GROUP1",
		// 指定实例的上线状态，默认为 true
		Ephemeral: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	if !OK {
		log.Fatal("注册nacos服务失败")
	}

	//服务发现
	_, err = nClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		//指定服务实例所属的集群名称
		Clusters: []string{"Cluster1"},
		//定了要选择健康服务实例的服务名
		ServiceName: global.SevConf.ServiceName,
		//指定了服务的分组名称
		GroupName: "GROUP1",
	})
	if err != nil {
		log.Fatal(err)
	}

}

// 数据库初始化
func Mysql() {
	var err error
	global.MysqlDB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", global.SevConf.Mysql.Root, global.SevConf.Mysql.Password, global.SevConf.Mysql.Host, global.SevConf.Mysql.Port, global.SevConf.Mysql.Database)), &gorm.Config{})
	if err != nil {
		zap.S().Panic("数据库连接失败", err)
		return
	}

	log.Println("mysql数据库初始化完成")
}

// consul健康检测服务
func Consul() {
	//创建默认客户端配置
	conf := api.DefaultConfig()

	//配置consul路由
	conf.Address = "127.0.0.1:8500" //默认

	//实例化客户端
	client, err := api.NewClient(conf)
	if err != nil {
		log.Fatal(err)
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
	port, err := strconv.Atoi(global.SevConf.RpcPort)
	if err != nil {
		log.Fatal("数据类型转换失败", err)
	}
	reg := api.AgentServiceRegistration{
		Address: "10.2.171.69",
		Port:    port,
		Name:    "first _srv",
		//Tags:    []string{"tag1"},
		ID:    "first_id", //若不填默认是Name
		Check: check,
	}

	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("consul注册服务成功")
}
