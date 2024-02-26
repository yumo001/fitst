package rpcSer

import (
	"gitee.com/yumo01/bag/initialize"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func rpcInit() {
	zap.S().Info("正在初始化配置文件")
	initialize.Zap()
	initialize.Viper()
	initialize.Nacos()
	initialize.Mysql()
	//initialize.Consul()
	zap.S().Info("配置文件初始化完成")
}

func RegisterGrpg(port string, f func(ser *grpc.Server)) error {
	rpcInit()

	//创建tpc监听
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Println("创建rpc监听服务失败", err)
		return err
	}
	log.Println("创建rpc监听服务成功,监听端口:" + port)

	server := grpc.NewServer()
	//反射接口
	reflection.Register(server)
	f(server)

	err = server.Serve(lis)
	if err != nil {
		log.Println("rpc服务启动失败", err)
		return err
	}

	defer lis.Close()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	return nil
}
