package logic

import (
	"context"
	"github.com/yumo001/fitst/pb"
	"google.golang.org/grpc"
)

// 创建本地服务结构体类型
// 模仿grpc.go中的方法实例
type SerS struct {
	pb.UnimplementedThisServer
}

func RegisterGrpc(ser grpc.ServiceRegistrar) {
	pb.RegisterThisServer(ser, SerS{})
}

func (ser SerS) Ping(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	ping := in.Ping
	if ping == "ping" {
		return &pb.Response{
			Pong: "pong",
		}, nil
	}

	return &pb.Response{
		Pong: "",
	}, nil
}

func (ser SerS) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	return nil, nil
}

func (ser SerS) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {

	return &pb.UserLoginResponse{}, nil
}
