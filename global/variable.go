package global

import (
	"gorm.io/gorm"
)

// 配置文件的映射结构
type SeverConfig struct {
	ServiceName string      `mapstructure:"serverName"`
	RpcPort     string      `json:"rpc_port"`
	NacosConfig NacosConfig `mapstructure:"nacosConfig"`
	Mysql       MysqlConfig `mapstructure:"mysqlConfig"`
}

type NacosConfig struct {
	ServerConfig struct {
		IpAddr string `mapstructure:"IpAddr"`
		Port   uint64 `mapstructure:"Port"`
	}
	ClientConfig struct {
		NamespaceId string `mapstructure:"NamespaceId"`
	}
	ConfigParam struct {
		DataId string `mapstructure:"DataId"`
		Group  string `mapstructure:"Group"`
	}
}

type MysqlConfig struct {
	Root     string `mapstructure:"root"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

var (
	MysqlDB = &gorm.DB{}
	// 配置文件结构体实例
	SevConf = SeverConfig{}
)
