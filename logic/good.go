package logic

import (
	"context"
	goods "github.com/yumo001/fitst/pb/goods"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 创建本地服务结构体类型
// 模仿grpc.go中的方法实例
type GoodsServer struct {
	goods.UnimplementedGoodsServer
}

func (ser GoodsServer) AddGood(ctx context.Context, in *goods.AddGoodRequest) (*goods.AddGoodResponse, error) {
	//测试阶段-避免数据冗余
	return &goods.AddGoodResponse{}, status.Errorf(codes.OK, "添加成功")
}
