// Package module 模块与应用集成的核心功能
// 主要定义模块的功能和接口，并实现与 cosmos SDK 框架的标准交互
// 使得模块可以在应用中注册、初始化、运行和导出自己的状态
package module

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/buzzing/checkers"
	"github.com/buzzing/checkers/keeper"
)

var (
	// AppModuleBasic 提供基本模块功能
	// 包括模块名称、Amino 编解码器注册、接口注册和 gRPC 注册
	_ module.AppModuleBasic = AppModule{}

	// HasGenesis 定义 Genesis 相关逻辑
	// 包括对默认 Genesis 状态的生成、验证、初始化和导出
	_ module.HasGenesis = AppModule{}

	// AppModule 不包含功能，模块通过实现 AppModule 接口来被其他模块注入
	_ appmodule.AppModule = AppModule{}
	// 包含 BeginBlock 方法
	_ appmodule.HasBeginBlocker = AppModule{}
	// 包含 EndBlock 方法
	_ appmodule.HasEndBlocker = AppModule{}
)

// ConsensusVersion 定义当前模块的共识版本
const ConsensusVersion = 1

type AppModule struct {
	cdc codec.Codec
	// AppModule 通过 keeper 实现对模块功能的访问
	keeper keeper.Keeper
}

// NewAppModule 创建一个新的 AppModule 实例
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func NewAppModuleBasic(m AppModule) module.AppModuleBasic {
	return module.CoreAppModuleBasicAdaptor(m.Name(), m)
}

// Name 返回模块的名称
func (AppModule) Name() string {
	return checkers.ModuleName
}

// RegisterLegacyAminoCodec 注册模块的 Amino 编解码器
// Amino 是 cosmos SDK 的早期序列化框架，现已被 Protobuf 取代
// 在现代模块中留空即可
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterGRPCGatewayRoutes 注册模块的 gRPC 网关路由，暴露 REST 接口，为 HTTP 客户端提供支持
// clientCtx 为客户端上下文，用于构造查询客户端
// mux 为 gRPC Gateway 的路由器，用于注册 REST 路由服务
// gRPC Gateway 可以将 HTTP 请求转换为 gRPC 请求
// 通过将 gRPC 服务注册到 mux 中，使得 gRPC 服务可以通过 HTTP 访问
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	// 对比于 MsgServer，消息的作用是修改状态，通过交易的方式广播到链上
	// 消息的入口是通过交易广播的，而不是直接通过 HTTP 或 gRPC 接口访问
	// QueryServer 查询是只读操作，支持 HTTP 和 gRPC 接口访问
	// 因此消息不需要 gRPC Gateway，只需要注册 QueryServer

	// checkers.RegisterQueryHandlerClient 将 gRPC 的查询接口转换为 REST API
	if err := checkers.RegisterQueryHandlerClient(context.Background(), mux, checkers.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces 注册模块的接口
func (AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	checkers.RegisterInterfaces(registry)
}

// ConsensusVersion 返回模块的共识版本
func (AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// RegisterServices 注册模块的 gRPC 服务，仅为 gRPC 服务提供支持
// 例如消息处理 MsgServer 和查询服务 QueryServer
// 通过模块的 Keeper 实现服务的具体逻辑
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// proto 中定义的服务接口 service Msg 生成方法 RegisterMsgServer
	// 用于负责将实现绑定到 gRPC 服务注册器
	// 在 keeper 包中实现了具体的 MsgServer，参见 keeper/msg_server.go
	// keeper.NewMsgServerImpl(am.keeper) 返回实现了 proto 中定义的 MsgServer 接口的对象
	// 通过 RegisterMsgServer 将其绑定到 gRPC 服务注册器
	checkers.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	// 同理在 keeper 包中实现了具体的 QueryServer，参见 keeper/query_server.go
	// keeper.NewQueryServerImpl(am.keeper) 返回实现了 proto 中定义的 QueryServer 接口的对象
	// 通过 RegisterQueryServer 将其绑定到 gRPC 服务注册器
	checkers.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// DefaultGenesis 返回默认的创世状态，并进行序列化
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(checkers.NewGenesisState())
}

// ValidateGenesis 验证创世状态是否有效，包括解析二进制数据并验证字段
func (AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data checkers.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", checkers.ModuleName, err)
	}

	return data.Validate()
}

// InitGenesis 初始化模块的创世状态
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState checkers.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(fmt.Errorf("failed to initialize %s genesis state: %w", checkers.ModuleName, err))
	}
}

// ExportGenesis 导出模块的创世状态
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to export %s genesis state: %w", checkers.ModuleName, err))
	}

	return cdc.MustMarshalJSON(gs)
}

func (am AppModule) BeginBlock(goCtx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	height := sdkCtx.BlockHeight()
	record := fmt.Sprintf("%s by BeginBlocker at %d", checkers.GetFormatDates(), height)
	err := am.keeper.RecordList.Set(goCtx, record)
	if err != nil {
		return err
	}
	return nil
}

func (am AppModule) EndBlock(goCtx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	height := sdkCtx.BlockHeight()
	record := fmt.Sprintf("%s by EndBlocker at %d", checkers.GetFormatDates(), height)
	err := am.keeper.RecordList.Set(goCtx, record)
	if err != nil {
		return err
	}
	return nil
}
