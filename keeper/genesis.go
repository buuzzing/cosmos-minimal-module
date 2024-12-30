// Package keeper 模块初始状态管理
// 该文件可以导出或导入模块的状态，确保模块状态在链重启或升级时的正确和安全
//
// keeper/genesis.go 主要作用
//   - 模块的初始状态加载
//     InitGenesis 负责从链的创世纪文件中读取初始数据，并初始化链上存储
//   - 模块的状态导出
//     ExportGenesis 负责从链的存储中读取当前状态，并生成可用作创世文件的数据
package keeper

import (
	"context"

	"github.com/buzzing/checkers"
)

// InitGenesis 从创世状态初始化模块
func (k *Keeper) InitGenesis(ctx context.Context, data *checkers.GenesisState) {
	// 初始化模块 Keeper 的 Params
	// 参见 keeper/keeper.go 中的 Params 字段
	if err := k.Params.Set(ctx, data.Params); err != nil {
		panic(err)
	}

	// 初始化模块 Keeper 的 Counter
	// 参见 keeper/keeper.go 中的 Counter 字段
	//TODO
}

// ExportGenesis 将模块状态导出到创世状态
func (k *Keeper) ExportGenesis(ctx context.Context) (*checkers.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &checkers.GenesisState{
		Params: params,
		//TODO
	}, nil
}
