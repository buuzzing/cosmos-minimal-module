// Package checkers 键定义
// keys.go 定义了模块中用到的存储键，这是模块存储数据的唯一标识
// 存储键是模块保存和查询链上数据的地址，类似于数据库的表名和字段名
//
// 通常包括
//   - 存储空间的定义
//   - 存储键的命名
package checkers

import "cosmossdk.io/collections"

const (
	// ModuleName 定义模块的名称
	ModuleName = "checkers"

	// Version 定义 IBC 模块支持的版本
	Version = "checkers-1"

	// PortId 定义模块的默认绑定端口
	PortId = "checkers"
)

// MaxIndexLength 定义游戏状态索引的最大长度
const MaxIndexLength = 256

var (
	ParamsKey      = collections.NewPrefix("Params")
	StoredGamesKey = collections.NewPrefix("StoredGames/value/")
	RecordKey      = collections.NewPrefix("Record/value/")
)
