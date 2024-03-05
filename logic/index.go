package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/yumo001/fitst/global"
	"github.com/yumo001/fitst/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		}, status.Errorf(codes.OK, "")
	}
	return &pb.Response{
		Pong: "",
	}, status.Errorf(codes.OK, "")
}

func (ser SerS) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	var count int64
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Count(&count).Error; err != nil {
		return &pb.UserRegisterResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if count > 0 {
		return &pb.UserRegisterResponse{}, status.Errorf(codes.AlreadyExists, "该用户已存在")
	}

	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	if err := global.MysqlDB.Table("users").Create(&in.U).Error; err != nil {
		return &pb.UserRegisterResponse{}, status.Errorf(codes.Canceled, "创建用户失败"+err.Error())
	}
	return nil, status.Errorf(codes.OK, "成功")
}

func (ser SerS) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var count int64
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Count(&count).Error; err != nil {
		return &pb.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if count <= 0 {
		return &pb.UserLoginResponse{}, status.Errorf(codes.NotFound, "该用户不存在")
	}

	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	var pp string
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Pluck("password", &pp).Error; err != nil {
		return &pb.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if in.U.Password != pp {
		return &pb.UserLoginResponse{}, status.Errorf(codes.Canceled, "密码错误")
	}

	return &pb.UserLoginResponse{}, status.Errorf(codes.OK, "")
}

func (ser SerS) List(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {

	var us []*pb.User
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Find(&us).Error; err != nil {
		return &pb.UserListResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	return &pb.UserListResponse{
		Users: us,
	}, nil
}
