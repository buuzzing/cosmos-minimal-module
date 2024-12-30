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

	return nil
}
