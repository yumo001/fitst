package global

import (
	"gorm.io/gorm"
)

// 配置文件的映射结构
type SeverConfig struct {
	ServiceName   string        `mapstructure:"serverName"`
	RpcPort       string        `yaml:"rpc_port"`
	NacosConfig   NacosConfig   `mapstructure:"nacosConfig"`
	Mysql         MysqlConfig   `yaml:"mysql"`
	Elastic       ElasticConfig `yaml:"elastic"`
	JwtSigningKey string        `yaml:"JwtSigningKey"`
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
	Root     string `yaml:"root"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type ElasticConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	MysqlDB = &gorm.DB{}
	// 配置文件结构体实例
	SevConf = SeverConfig{}
)
