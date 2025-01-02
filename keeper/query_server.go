// Package keeper 模块查询服务
// RPC 服务和方法定义参见 query.proto
// 这里需要做:
//   - 实现 GetGame 函数
package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/buzzing/checkers"
)

type queryServer struct {
	k Keeper
}

var _ checkers.QueryServer = queryServer{}

// NewQueryServerImpl 返回一个实现了 checkers.QueryServer 接口的对象
func NewQueryServerImpl(keeper Keeper) checkers.QueryServer {
	return queryServer{k: keeper}
}

// GetGame QueryGetGameRequest 消息的 handler，获取游戏内容
func (qs queryServer) GetGame(ctx context.Context, req *checkers.QueryGetGameRequest) (*checkers.QueryGetGameResponse, error) {
	game, err := qs.k.StoredGames.Get(ctx, req.Index)
	// 如果找到，则返回游戏
	if err != nil {
		return &checkers.QueryGetGameResponse{Game: &game}, err
	}
	// 如果未找到，则返回 nil
	if errors.Is(err, collections.ErrNotFound) {
		return &checkers.QueryGetGameResponse{Game: nil}, nil
	}

	// 否则返回错误
	return nil, status.Error(codes.Internal, err.Error())
}
