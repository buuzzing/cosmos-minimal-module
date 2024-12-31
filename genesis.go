// Package checkers 创世纪状态
// genesis.go 定义了模块的初始状态
// 用于在区块链启动时加载模块的初始配置和数据
//
// 通常包括
//   - Genesis 数据结构的定义，这里位于 types.pb.go 中
//   - 初始化与导出逻辑
//     初始化用于加载 Genesis 文件
//     导出用于生成当前模块的状态快照
package checkers

func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// 验证创世状态中的所有游戏状态
	unique := make(map[string]bool)
	for _, indexedStoredGame := range gs.IndexedStoredGameList {
		// 索引长度验证
		if length := len([]byte(indexedStoredGame.Index)); length < 1 || length > MaxIndexLength {
			return ErrIndexTooLong
		}
		// 索引唯一性验证
		if _, ok := unique[indexedStoredGame.Index]; ok {
			return ErrDuplicateAddress
		}
		// 游戏状态验证
		if err := indexedStoredGame.StoredGame.Validate(); err != nil {
			return err
		}
		unique[indexedStoredGame.Index] = true
	}

	return nil
}
