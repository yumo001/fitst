package main

import (
	"github.com/yumo001/fitst/global"
	"github.com/yumo001/fitst/initialize"
	"github.com/yumo001/fitst/rpcClient"
	"google.golang.org/grpc"
	"log"
)

func init() {
	log.Println("开始配置初始化......")
	initialize.Viper()
	initialize.Nacos()
	initialize.Mysql()
	initialize.CreatElasticClient()
	initialize.SynchronizationUser()
	log.Println("配置初始化完成......")
}

func main() {
	err := rpcClient.GrpcActivate(global.SevConf.RpcPort, func(s *grpc.Server) {
		rpcClient.RegisterGrpc(s)
		initialize.Consul(global.SevConf.RpcPort)
		log.Println("服务启动，监听端口:" + global.SevConf.RpcPort + "......")
	})
	if err != nil {
		log.Println("")
		panic(err)
	}
}
