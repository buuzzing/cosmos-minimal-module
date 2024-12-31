// Package keeper 模块存储管理与业务逻辑
// Keeper 连接链上状态、模块逻辑、和其他模块之间的交互
//
// Keeper 主要作用
//   - 管理模块状态
//     Keeper 提供操作链上存储的接口，例如查询和更新
//     它通过与 cosmos SDK 的 KVStore 交互，管理模块的链上数据
//   - 封装模块逻辑
//     模块的核心业务逻辑通常实现为 Keeper 的方法
//     这些方法对外提供模块的功能，例如处理交易、查询数据等
//   - 模块间交互
//     Keeper 是模块间通信的接口，通过 Keeper 访问其他模块的功能
//   - 依赖注入
//     Keeper 可以持有其他模块的 Keeper，从而访问和调用其逻辑
package keeper

import (
	"fmt"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/buzzing/checkers"
	"github.com/cosmos/cosmos-sdk/codec"

	// collections 是 cosmos SDK 提供的一个存储抽象层，提供了更高级和
	// 类型安全的存储接口，以替代传统的 KVStore 操作
	"cosmossdk.io/collections"
)

type Keeper struct {
	// 数据序列化和反序列化的编解码器
	cdc codec.BinaryCodec
	// 用于处理区块链地址的编码和解码
	// 当模块与账户或模块交互式，来确保地址的正确编码和解析
	addressCodec address.Codec

	// authority 记录允许执行某些高权限操作（例如 MsgUpdateParams）的地址
	// 通常设置为 x/gov 模块账户
	authority string

	// Schema 用于存储和管理模块的状态结构，类似于数据表
	// 模块的存储空间将由 Schema 组织和管理
	// Schema 用来定义模块需要的所有存储项（例如 collections.Item 和 sollections.Map）
	// collections.Schema 是 cosmos SDK 的集合框架，用于简化模块的存储定义和操作
	Schema collections.Schema
	// Params 用于存储模块的配置参数
	// collections.Item 是一种用于存储单一值的类型，本质上是一个 noKey 的 Map
	Params collections.Item[checkers.Params]
	// StoredGames 用于存储游戏数据
	// collections.Map 是一种映射类型，用于存储键值对
	// 它提供了一组方法来操作键值对，例如 Set、Get、Has、Remove 等
	// 因为 collections.Item 是一个 noKey 的 Map，它重写了 Set, Get 等方法，数据操作与 Map 类似
	StoredGames collections.Map[string, checkers.StoredGame]
}

// NewKeeper 创建一个新的 Keeper 实例
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) Keeper {
	// 通过 addressCodec 检查 authority 是否是有效的地址
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	// storeService 类型为 KVStoreService，是 cosmos SDK 提供的 KVStore 的抽象接口
	// 通过 KVStoreService 可以访问和操作 KVStore
	// 每个模块通过 KVStoreService 绑定一个唯一的命名空间，以避免不同模块之间的键冲突
	// NewSchemaBuilder 通过 storeService 来绑定模块的存储命名空间
	// 可以使用 SchemaBuilder 定义存储项，然后将这些存储项添加到 Keeper 的 Schema 中
	sb := collections.NewSchemaBuilder(storeService)

	// 创建一个内含类型为 checkers.Params 的 collections.Item
	// 并通过 SchemaBuilder 将此存储项注册到模块的存储架构中
	//   - sb: SchemaBuilder 实例
	//   - checkers.ParamsKey: 存储项的唯一标识符（键），用于区分模块中的不同存储项
	//   - "params": 存储项的名称，用于标识存储项的用途
	//   - codec.CollValue[checkers.Params](cdc)
	//     codec.CollValue 指定编码和解码规则(它是一个 func)，泛型参数为 checkers.Params
	//     cdc 为编解码器，用于序列化和反序列化数据
	// collections.NewItem 返回一个 collections.Item 实例，它的使用方法为:
	//   - 存储值: params.Set(ctx, checkers.Params{...})
	//     底层通过 cdc 将 checkers.Params 编码为 bytes
	//   - 读取值: params.Get(ctx) 返回 checkers.Params
	//     底层通过 cdc 将 bytes 解码为 checkers.Params
	params := collections.NewItem(sb, checkers.ParamsKey, "params", codec.CollValue[checkers.Params](cdc))
	// 创建一个内涵类型为 string->checkers.StoredGame 的 collections.Map
	// collections.StringKey 指定键类型和编码规则
	// codec.CollValue[checkers.StoredGame](cdc) 指定值类型和编码规则，泛型参数为 checkers.StoredGame
	// collections.NewMap 返回一个 collections.Map 实例，它的使用方法为:
	//   - 存储值: storedGames.Set(ctx, key, checkers.StoredGame{...})
	//     底层通过 cdc 将 checkers.StoredGame 编码为 bytes
	//   - 读取值: storedGames.Get(ctx, key) 返回 checkers.StoredGame
	//     底层通过 cdc 将 bytes 解码为 checkers.StoredGame
	storedGames := collections.NewMap(sb, checkers.StoredGamesKey, "storedGames", collections.StringKey, codec.CollValue[checkers.StoredGame](cdc))

	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params:      params,
		StoredGames: storedGames,
	}

	// 通过 SchemaBuilder 构建模块的存储架构
	// SchemaBuilder 会将所有存储项注册到模块的存储架构中
	// 最后通过 Build 方法将存储架构构建为 Schema
	schema, err := sb.Build()
	if err != nil {
		panic(fmt.Errorf("failed to build schema: %w", err))
	}

	k.Schema = schema

	return k
}

// GetAuthority 返回模块的权限账户地址
func (k *Keeper) GetAuthority() string {
	return k.authority
}
