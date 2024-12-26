// Package checkers 参数
// params.go 定义模块的参数，参数用于控制模块的行为和逻辑
//
// 通常包括
//   - 模块参数的定义，这里位于 types.pb.go 中
//   - 默认值 DefaultParams() 和验证逻辑 Validate()
//
// 模块参数类似于模块的“配置文件”，可以动态调整模块行为（例如修改交易费率等）
// 通过治理提案可以修改参数，以适应链上运行的需求
package checkers

// DefaultParams 返回默认的模块参数
func DefaultParams() Params {
	return Params{}
}

// Validate 对参数进行检查
func (p Params) Validate() error {
	return nil
}
