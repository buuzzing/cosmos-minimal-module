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
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces 注册模块中的接口类型和它们的具体实现
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// RegisterImplementations 将模块中的某些消息类型注册为特定接口的实现
	// sdk.Msg 是 cosmos SDK 中的标准接口，用于定义所有可以通过交易提交的消息
	// MsgCreateGame 为某个特定的交易消息，注册这个消息为 sdk.Msg 接口的实现
	// 使得框架能够识别和处理 MsgCreateGame 这个类型
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGame{},
		&MsgAddRecord{},
	)
	// cosmos SDK 使用接口注册机制来处理多态类型
	// 通过注册接口和具体实现的关系，使得框架能够正确地序列化和反序列化这些类型

	// RegisterMsgServiceDesc 将模块中定义的 gRPC 服务描述注册到接口注册表中
	// 注册后，模块的 gRPC 服务就能够通过 cosmos SDK 框架提供的统一接口被访问
	// _Msg_serviceDesc 为模块中定义的 gRPC 服务描述
	// 它是 proto 生成的类型为 grpc.ServiceDesc 的结构体
	// 描述了 gRPC 服务的定义（service Msg 定义）以及服务内的所有方法
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	// gRPC 服务的注册是为了支持模块的查询和操作接口（通过 gRPC 或 REST API 提供）
	// 通过注册服务描述，框架能够基于服务描述生成路由和响应逻辑
}
