package main

import (
	"gitee.com/yumo01/bag/global"
	"gitee.com/yumo01/bag/initialize"
	"gitee.com/yumo01/bag/logic"
	"gitee.com/yumo01/bag/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
)

func init() {
	zap.S().Info("正在初始化配置文件")
	initialize.Zap()
	initialize.Viper()
	initialize.Nacos()
	initialize.Mysql()
	//initialize.Consul()
	zap.S().Info("配置文件初始化完成")
}

var ser logic.SerS

func main() {

	server := grpc.NewServer()
	//实现具体服务
	pb.RegisterUserServer(server, ser)
	//注册健康检测服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//创建tpc监听
	lis, err := net.Listen("tcp", "0.0.0.0:"+global.SevConf.RpcPort)
	if err != nil {
		zap.S().Panic("创建rpc监听服务失败", err)
		return
	}
	zap.S().Info("创建rpc监听服务成功,监听端口:" + global.SevConf.RpcPort)
	defer lis.Close()

	//不断接收客户端连接
	//创建一个 gRPC 服务器并开始监听指定的网络地址（lis），以便接受客户端的 gRPC 连接请求。
	go func() {
		err = server.Serve(lis)
		if err != nil {
			zap.S().Panic("rpc服务启动失败", err)
			return
		}
	}()

	//阻塞等待客户端断开连接--优雅断开连接

	//创建了一个类型为 os.Signal 的通道 quit。
	//通道是 Go 语言中用于在 goroutine 之间进行通信的机制，
	//而 os.Signal 类型用于表示操作系统发送的信号，如中断信号。
	quit := make(chan os.Signal)
	//signal.Notify 函数用于将操作系统发出的信号通知到指定的通道。
	//在这里，它将操作系统的中断信号（os.Interrupt）发送到 quit 通道中。
	signal.Notify(quit, os.Interrupt)
	//<-quit：这一行代码是一个接收操作，它会阻塞当前的 goroutine，直到从 quit 通道中接收到一个值。这里的目的是等待操作系统发出的中断信号。
	<-quit

}
