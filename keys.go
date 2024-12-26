// Package checkers 键定义
// keys.go 定义了模块中用到的存储键，这是模块存储数据的唯一标识
// 存储键是模块保存和查询链上数据的地址，类似于数据库的表名和字段名
//
// 通常包括
// + 存储空间的定义
// + 存储键的命名
package checkers

import "cosmossdk.io/collections"

const ModuleName = "checkers"

var (
	ParamsKey = collections.NewPrefix("Params")
)
