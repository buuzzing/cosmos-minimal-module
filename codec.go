// Package checkers 编解码器
// codec.go 负责注册模块中使用的数据类型的编码和解码规则
// 确保模块的消息、查询相应、状态存储等内容能够正确地序列化和反序列化
//
// 通常包括
//   - ProtoBuf 编码注册 RegisterCodec(cdc *codec.ProtoCodec)
//     注册模块中定义的所有消息和结构体
//   - 接口注册 RegisterInterfaces(registry types.InterfaceRegistry)
//     注册模块中的接口类型和它们的具体实现
//   - Codec 实例
//     定义一个全局的 Codec 实例用于其他模块调用
package checkers

import (
	types "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {}
