// Package module 模块 CLI (command line interface) 配置
// autocli.go 用于自动生成模块的命令行接口
// 也可以绑定模块的 gRPC 服务或消息处理器到 CLI
package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	checkersv1 "github.com/buzzing/checkers/api/v1"
)

// AutoCLIOptions 实现 autocli.HasAutoCLIConfig 接口
// 返回内容 *autocliv1.ModuleOptions 描述了模块的 CLI 配置
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: checkersv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "GetGame",
					// index 为所需的参数
					Use:   "get-game index",
					Short: "Get the current value of the game at the index",
					// proto 文件中定义的 GetGame 方法参数为 QueryGetGameRequest (参见 query.proto)
					// 其 QueryGetGameRequest 需要 index 参数
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "index"},
					},
				},
				{
					RpcMethod: "GetRecordList",
					Use:       "get-record-list",
					Short:     "Get the list of records",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			// Service 模块消息服务的 gRPC 名称
			Service: checkersv1.Msg_ServiceDesc.ServiceName,
			// RpcCommandOptions 指定了该服务中所有支持的 RPC 方法的 CLI 命令选项
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					// RpcMethod 指定了 gRPC 服务中的方法名称
					RpcMethod: "CreateGame",
					// Use 指定命令的使用方法
					// index black red 为所需的参数
					Use: "create index black red",
					// Short 指定了命令的简短描述
					Short: "Creates a new checkers game at the index for the black and red players",
					// PositionalArgs 定义命令的参数及其顺序
					// 每个参数通过 ProtoField 指定其对应的字段
					// proto 文件中定义的 CreateGame 方法参数为 MsgCreateGame (参见 tx.proto)
					// 其 MsgCreateGame 需要 index, black, red 三个参数
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "index"},
						{ProtoField: "black"},
						{ProtoField: "red"},
					},
				},
				{
					RpcMethod: "AddRecord",
					Use:       "add-record record_value",
					Short:     "Add a record to the chain storage",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "value"},
					},
				},
			},
		},
	}
}
