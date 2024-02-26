package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"gitee.com/yumo01/bag/global"
	"gitee.com/yumo01/bag/pb"
	"gitee.com/yumo01/bag/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gomail "gopkg.in/gomail.v2"
)

//这是具体的服务端实现功能的logic层

// 创建本地服务结构体类型
// 模仿grpc.go中的方法实例
type SerS struct {
	pb.UnimplementedUserServer
}

func (ser SerS) Ping(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	ping := req.Ping
	if ping == "ping" {
		return &pb.Response{
			Pong: "pong",
		}, nil
	}

	return &pb.Response{
		Pong: "",
	}, nil
}

func (ser SerS) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	u := req.U
	var count int64
	err := global.MysqlDB.Model(&u).Where("username = ?", req.U.Username).Count(&count).Error
	if err != nil {
		return &pb.RegisterResponse{}, status.Errorf(codes.NotFound, "数据库查询失败")
	}

	if count > 0 {
		return &pb.RegisterResponse{}, status.Errorf(codes.AlreadyExists, "该用户名已被占用")
	}

	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(req.U.Password)))
	err = global.MysqlDB.Create(&u).Error
	if err != nil {
		return &pb.RegisterResponse{}, status.Errorf(codes.Canceled, "数据库写入失败")
	}

	return &pb.RegisterResponse{}, nil

}

func (ser SerS) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var u pb.Users
	var count int64
	err := global.MysqlDB.Model(&u).Where("username = ?", req.U.Username).Count(&count).Error
	if err != nil {
		return &pb.LoginResponse{}, status.Errorf(codes.NotFound, "mysql查询失败")
	}

	if count <= 0 {
		return &pb.LoginResponse{}, status.Errorf(codes.Canceled, "该用户不存在")
	}
	var pp string
	err = global.MysqlDB.Raw("select `password` from users where username = ?", req.U.Username).Scan(&pp).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "mysql查询失败")
	}

	if fmt.Sprintf("%x", md5.Sum([]byte(req.U.Password))) != pp {
		return &pb.LoginResponse{}, status.Errorf(codes.Canceled, "密码错误")
	}

	return &pb.LoginResponse{}, nil
}

// 找回密码
func (ser SerS) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	usernmae := req.Username

	//获取手机号与邮箱
	var u pb.Users
	err := global.MysqlDB.Raw("select * from users where username = ?", usernmae).Scan(&u).Error
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "数据库查询失败")
	}

	//通过邮箱验证
	//err = sendEmail(u.Email)
	//if err != nil {
	//	return nil, status.Errorf(codes.Canceled, "发送邮件失败")
	//}

	//通过手机号验证
	err = utils.HuYiSms(u.Tel)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "发送短信失败")
	}

	return &pb.ForgotPasswordResponse{}, nil
}

// 发送至邮箱
func sendEmail(email string) error {

	// 创建邮件对象
	msg := gomail.NewMessage()
	// 设置发件人
	msg.SetHeader("From", "yumo040410@163.com")
	// 设置收件人
	msg.SetHeader("To", email)
	// 设置邮件主题
	msg.SetHeader("Subject", "hello,subject!")
	// 设置邮件正文
	msg.SetBody("text/html", "<span>This is the body of the mail</span>")
	// msg.Attach("/home/User/cat.jpg")

	n := gomail.NewDialer("smtp.163.com", 25, "yumo040410@163.com", "WWASOEEHAQMQUSSD")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {

		return err
	}
	return nil
}

// 修改密码
func (ser SerS) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := global.MysqlDB.Exec("update users set password = ? WHERE username = ?", req.Password, req.Username).Error
	if err != nil {
		return &pb.ResetPasswordResponse{}, status.Errorf(codes.Canceled, "修改密码失败")
	}

	return &pb.ResetPasswordResponse{}, nil
}

func (ser SerS) UserList(ctx context.Context, rep *pb.UserListRequest) (*pb.UserListResponse, error) {

	return &pb.UserListResponse{
		Us: []*pb.Users{},
	}, nil
}
