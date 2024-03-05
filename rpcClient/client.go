package rpcClient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func GrpcActivate(port string, f func(s *grpc.Server)) error {

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Println("监听端口失败:" + port)
		panic(err)
	}
	//创建grpc服务器
	s := grpc.NewServer()

	//注册健康检测服务
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(s, healthcheck)

	//反射接口支持查询
	reflection.Register(s)
	f(s)

	//创建一个 gRPC 服务器并开始监听指定的网络地址（lis），以便接受客户端的 gRPC 连接请求。
	err = s.Serve(lis)
	if err != nil {
		log.Println("rpc服务启动失败", err)
		panic(err)
	}

	return nil
}
