package shared

import "context"

// TransactionManager 事务管理器接口
// 设计理念：简单实用，一个接口搞定所有事务需求
type TransactionManager interface {
	// WithTransaction 在事务中执行业务逻辑
	// 如果fn返回error，自动回滚；否则自动提交
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	// WithTransactionResult 在事务中执行业务逻辑并返回结果
	// 注意：接口方法不能有类型参数，所以使用interface{}
	WithTransactionResult(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
}

// TransactionContextKey 事务上下文键
// 用于在上下文中传递事务实例
type TransactionContextKey string

const TransactionKey TransactionContextKey = "transaction"

// 为什么这样设计？
// 1. 接口简单：只有2个方法，容易理解和使用
// 2. 自动管理：开发者不需要手动开启/提交/回滚事务
// 3. 类型安全：使用泛型支持返回值，避免类型转换
// 4. 上下文传递：通过context传递事务，符合Go的惯例
