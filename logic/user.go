package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/yumo001/fitst/global"
	user "github.com/yumo001/fitst/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 创建本地服务结构体类型
// 模仿grpc.go中的方法实例
type UserServer struct {
	user.UnimplementedThisServer
}

func (ser UserServer) Ping(ctx context.Context, in *user.Request) (*user.Response, error) {
	ping := in.Ping
	if ping == "ping" {
		return &user.Response{
			Pong: "pong",
		}, status.Errorf(codes.OK, "")
	}
	return &user.Response{
		Pong: "",
	}, status.Errorf(codes.OK, "")
}

func (ser UserServer) UserRegister(ctx context.Context, in *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {

	var count int64
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Count(&count).Error; err != nil {
		return &user.UserRegisterResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if count > 0 {
		return &user.UserRegisterResponse{}, status.Errorf(codes.AlreadyExists, "该用户已存在")
	}

	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	if err := global.MysqlDB.Table("users").Create(&in.U).Error; err != nil {
		return &user.UserRegisterResponse{}, status.Errorf(codes.Canceled, "创建用户失败"+err.Error())
	}
	return nil, status.Errorf(codes.OK, "成功")
}

func (ser UserServer) UserLogin(ctx context.Context, in *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	var count int64

	//var err error
	//initialize.Mysql2(func(mysqlDB *gorm.DB) error {
	//	if err = mysqlDB.Table("users").Where("username = ?", in.U.Username).Count(&count).Error; err != nil {
	//		return status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	//	}
	//	if count <= 0 {
	//		return status.Errorf(codes.NotFound, "该用户不存在")
	//	}
	//
	//	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	//	var pp string
	//	if err := mysqlDB.Table("users").Where("username = ?", in.U.Username).Pluck("password", &pp).Error; err != nil {
	//		return status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	//	}
	//	if in.U.Password != pp {
	//		return status.Errorf(codes.Canceled, "密码错误")
	//	}
	//	return nil
	//})

	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Count(&count).Error; err != nil {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if count <= 0 {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "该用户不存在")
	}

	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	var pp string
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Pluck("password", &pp).Error; err != nil {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if in.U.Password != pp {
		return &user.UserLoginResponse{}, status.Errorf(codes.Canceled, "密码错误")
	}

	return &user.UserLoginResponse{}, status.Errorf(codes.OK, "")
}

func (ser UserServer) UserLoginByPhone(ctx context.Context, in *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	var count int64
	if err := global.MysqlDB.Table("users").Where("phone = ?", in.U.Phone).Count(&count).Error; err != nil {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	if count <= 0 {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "该用户不存在")
	}

	in.U.Password = fmt.Sprintf("%x", md5.Sum([]byte(in.U.Password)))
	var pp string
	if err := global.MysqlDB.Table("users").Where("phone = ?", in.U.Phone).Pluck("password", &pp).Error; err != nil {
		return &user.UserLoginResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}

	if in.U.Password != pp {
		return &user.UserLoginResponse{}, status.Errorf(codes.Canceled, "密码错误")
	}
	return &user.UserLoginResponse{}, status.Errorf(codes.OK, "")
}

func (ser UserServer) UserList(ctx context.Context, in *user.UserListRequest) (*user.UserListResponse, error) {

	var us []*user.User
	//if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).Find(&us).Error; err != nil {
	//	return &user.UserListResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	//}

	query := elastic.NewMatchAllQuery()
	res, err := global.ElasticClient.Search().Index().Query(query).Do(ctx)
	if err != nil {
		return &user.UserListResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}
	fmt.Println(res.Hits.Hits)
	return &user.UserListResponse{
		Users: us,
	}, nil
}

func (ser UserServer) UserPasswordRecovery(ctx context.Context, in *user.UserPasswordRecoveryRequest) (*user.UserPasswordRecoveryResponse, error) {

	var u *user.User
	if err := global.MysqlDB.Table("users").Where("username = ?", in.U.Username).First(&u).Error; err != nil {
		return &user.UserPasswordRecoveryResponse{}, status.Errorf(codes.NotFound, "数据库查询失败"+err.Error())
	}

	return &user.UserPasswordRecoveryResponse{
		U: u,
	}, status.Errorf(codes.OK, "")
}
