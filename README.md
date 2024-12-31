# Checkers Minimal

西洋跳棋模块编写学习

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
