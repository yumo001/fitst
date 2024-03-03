package main

import (
	"github.com/yumo001/fitst/initialize"
	"github.com/yumo001/fitst/logic"
	"github.com/yumo001/fitst/rpcClient"
	"google.golang.org/grpc"
	"log"
)

func init() {
	initialize.Viper()
	initialize.Nacos()
}

func main() {
	err := rpcClient.GrpcActivate("8080", func(s *grpc.Server) {
		logic.RegisterGrpc(s)
	})
	if err != nil {
		log.Println("")
		panic(err)
	}

}
