// Package module 依赖注入（depinject）工具配置
// depinject.go 文件定义模块需要注入的依赖信息（例如 Keeper 和编解码器等）
// 通过 depinject 可以将模块与其他模块的接口连接起来，实现透明地调用其他模块的功能
package module

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	modulev1 "github.com/buzzing/checkers/api/module/v1"
	"github.com/buzzing/checkers/keeper"
)

var _ appmodule.AppModule = AppModule{}

// 确保 AppModule 实现了 porttypes.IBCModule 接口
var _ ibcporttypes.IBCModule = AppModule{}

// IsOnePerModuleType() 和 IsAppModule() 是标记方法，没有具体逻辑
// 只是用来告知 depinject 这是一个标准模块类型

// IsOnePerModuleType 实现 depinject.OnePerModuleType 接口
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule 实现 appmodule.AppModule 接口
func (am AppModule) IsAppModule() {}

// appmodule.Register 用来将模块注册到 cosmos SDK 的模块管理器中
// &modulev1.Module{} 是模块的配置信息，描述模块的基础信息
// appmodule.Provide(ProvideModule) 是模块的依赖注入函数，指定依赖提供者函数 ProvideModule
// ProvideModule 用于实例化模块的核心组件并返回
// 应用通过对 module 的匿名导入，执行该 init 函数，实现模块的注册
func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

// ModuleInputs 模块所需要的输入参数
type ModuleInputs struct {
	// depinject.In 用于标记该结构体为依赖注入的输入参数
	depinject.In

	// Cdc 编解码器，用于序列化和反序列化数据
	Cdc codec.Codec
	// StoreService 存储服务，用于状态管理
	StoreService store.KVStoreService
	// AddressCodec 地址编解码器，用于地址的编解码
	AddressCodec address.Codec

	// Config 模块配置信息，模块的元数据
	Config *modulev1.Module

	// IBC 相关依赖
	IBCKeeperFn        func() *ibckeeper.Keeper                   `optional:"true"`
	CapabilityScopedFn func(string) capabilitykeeper.ScopedKeeper `optional:"true"`
}

// ModuleOutputs 模块的输出参数
type ModuleOutputs struct {
	// depinject.Out 用于标记该结构体为依赖注入的输出参数
	depinject.Out

	// Module 模块实例，供应用框架使用
	Module appmodule.AppModule
	// Keeper 模块的核心逻辑管理器，供其他模块或组件调用
	Keeper keeper.Keeper
}

// ProvideModule 用于依赖注入中实例化模块的核心组件
func ProvideModule(in ModuleInputs) ModuleOutputs {
	// 默认的权限模块地址
	authority := authtypes.NewModuleAddress("gov")
	// 如果配置中有指定权限模块地址，则使用配置中的地址
	// 参见 proto/buzzing/checkers/module/v1/module.proto 对于该字段的定义
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.AddressCodec,
		in.StoreService,
		authority.String(),
		in.IBCKeeperFn,
		in.CapabilityScopedFn,
	)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{Module: m, Keeper: k}
}
