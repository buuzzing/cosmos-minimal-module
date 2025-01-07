# Checkers Minimal

cosmos 模块编写学习指南

> 参考文档: [SDK v0.50 Native](https://tutorials.cosmos.network/hands-on-exercise/0-native/)
> 
> 时间: 2024/12/31

## F.A.Q.

### 链、应用、模块 module、keeper 有何关系？

**链**

链是基于 Tendermint 共识引擎运行的分布式账本系统

cosmos 链类似于状态机，包括状态 state、交易 transaction 和共识 consensus

链是由应用构建的，是一组模块功能的载体。一个应用部署后会运行在 Tendermint 上，成为一条链

多个节点运行同一个应用构建的链，通过 Tendermint 达成一致，形成分布式网络；应用的功能通过链上的节点对外提供服务

**应用**

应用是基于 cosmos SDK 开发的具体区块链项目

应用通过注册 (depinject) 和调用模块 (module) 来实现其功能

一个应用通常包含多个模块，例如账户管理 `x/auth`、交易 `x/bank` 和其他用于实现应用逻辑的自定义模块等

**模块**

模块 module 是通过 cosmos SDK 构建的功能单元，也是应用的功能单元

每个模块专注于实现一个特性功能，例如账户管理 (`x/auth`)、治理 (`x/gov`) 或其他自定义逻辑等；一个应用通常包含多个模块，一个模块也可以被多个应用重用

模块的逻辑可以独立开发、测试和更新

模块通过实现标准接口 (appmodule.AppModule) 以与应用集成，模块通过持有其他模块的 keeper 或通过 depinject 与其他模块交互

模块会提供与外部系统的接口，包括 gRPC, CLI 或 REST 接口

**keeper**

每个模块都有自己的 keeper，负责:

+ 模块的状态管理: 直接与区块链状态交互，存储管理
+ 业务逻辑实现: 处理模块的业务逻辑
+ 对外接口: 提供接口以供本模块和其他模块使用
+ 外部依赖管理: 持有其他模块的 keeper 以使用其他模块的功能

module 封装了 keeper，以对接到应用；module 持有 keeper 实例，通过调用 keeper 提供的接口来完成核心功能

keeper 专注于模块的逻辑和存储操作，不需要考虑框架；module 专注于框架集成，而不需要考虑存储和具体的业务逻辑

### 定义一个链上状态和对应操作流程

1. 在 keeper/keeper.go 中定义该字段，*通常使用 collections*

   1.1 在 keys.go 中定义存储它的 key prefix

   1.2 在 keeper/keeper.go 的 NewKeeper() 中初始化它

2. 定义该字段的创世纪状态和验证逻辑

   2.1 在 genesis.go 的 NewGenesisState() 函数中定义字段的初始值

   2.2 在 genesis.go 的 Validate() 函数中对创世纪状态中的该字段进行合法性验证

3. 定义该字段在创世状态的导出导出逻辑

   3.1 在 types.proto 中的 GenesisState 消息中定义该字段

   3.2 在 keeper/genesis.go 的 InitGenesis() 函数中添加该字段的导入逻辑

   3.3 在 keeper/genesis.go 的 ExportGenesis() 函数中添加该字段的导出逻辑

4. 添加 message，实现对字段的增删改逻辑

   4.1 在 tx.proto 中定义增删改的 rpc 接口和对应参数

   4.2 在 keeper/msg_server 定义 message 的 handler，即实现 proto 文件中定义的 rpc 接口

   4.3 在 codec.go 的 RegisterInterfaces() 函数中注册 proto 生成的消息类型，注册 rpc 服务描述

   4.4 在 module/module.go 的 RegisterServices() 函数中将 MsgServer 的实现注册到 proto 文件

5. 添加 query，实现对字段的查询逻辑

   5.1 在 query.proto 中定义查询的 rpc 接口和对应参数

   5.2 在 keeper/query_server 定义 query 的 handler，即实现 proto 文件中定义的 rpc 接口

   5.3 在 module/module.go 的 RegisterServices() 函数中将 QueryServer 的实现注册到 proto 文件

   5.4 在 module/module.go 的 RegsterGRPCGatewayRoutes() 函数中将 QueryClient 注册为 HTTP

6. 配置 CLI

   6.1 在 module/autocli.go 的 AutoCLIOptions() 函数中的 Tx 部分添加期望为 message 生成的 CLI 命令

   6.2 在 module/autocli.go 的 AutoCLIOptions() 函数中的 Query 部分添加期望为 query 生成的 CLI 命令