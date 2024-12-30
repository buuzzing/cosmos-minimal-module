// Package module 模块 CLI (command line interface) 配置
// autocli.go 用于自动生成模块的命令行接口
// 也可以绑定模块的 gRPC 服务或消息处理器到 CLI
package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

// AutoCLIOptions 实现 autocli.HasAutoCLIConfig 接口
// 返回内容 *autocliv1.ModuleOptions 描述了模块的 CLI 配置
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: nil,
		Tx:    nil,
	}
}
