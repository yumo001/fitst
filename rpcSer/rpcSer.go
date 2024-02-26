package rpcSer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

func RegisterGrpg(port int, f func(ser *grpc.Server)) error {
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Println("创建rpc监听服务失败", err)
		return err
	}
	server := grpc.NewServer()
	//反射接口
	reflection.Register(server)
	f(server)

	log.Println("创建rpc监听服务成功,监听端口:" + strconv.Itoa(port))

	err = server.Serve(lis)
	if err != nil {
		log.Println("rpc服务启动失败", err)
		return err
	}

	return nil
}
