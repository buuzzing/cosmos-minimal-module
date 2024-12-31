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
func (k *Keeper) InitGenesis(ctx context.Context, data *checkers.GenesisState) error {
	// 初始化模块 Keeper 的 Params
	// 参见 keeper/keeper.go 中的 Params 字段
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// 初始化模块 Keeper 的 StoredGame
	// 参见 keeper/keeper.go 中的 StoredGame 字段
	for _, indexedStoredGame := range data.IndexedStoredGameList {
		if err := k.StoredGames.Set(ctx, indexedStoredGame.Index, indexedStoredGame.StoredGame); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis 将模块状态导出到创世状态
func (k *Keeper) ExportGenesis(ctx context.Context) (*checkers.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var indexedStoredGames []checkers.IndexedStoredGame
	// 遍历所有的 StoredGame
	// collections.Map 中的 Walk 方法用于遍历 Map，每一个元素都会调用回调函数，并传入反序列化后的 key 和 value
	// 如果回调函数返回 true，则停止遍历
	// range 为 nil 时，遍历所有元素
	if err := k.StoredGames.Walk(ctx, nil, func(index string, storedGame checkers.StoredGame) (bool, error) {
		indexedStoredGames = append(indexedStoredGames, checkers.IndexedStoredGame{
			Index:      index,
			StoredGame: storedGame,
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &checkers.GenesisState{
		Params:                params,
		IndexedStoredGameList: indexedStoredGames,
	}, nil
}
