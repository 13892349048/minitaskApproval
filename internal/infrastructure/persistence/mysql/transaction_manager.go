package mysql

import (
	"context"

	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TransactionManager GORM事务管理器实现
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) shared.TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction 在事务中执行业务逻辑
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 使用GORM的Transaction方法，它会自动处理开启/提交/回滚
	return tm.db.Transaction(func(tx *gorm.DB) error {
		// 将事务实例放入上下文，供Repository使用
		txCtx := context.WithValue(ctx, shared.TransactionKey, tx)

		// 执行业务逻辑
		if err := fn(txCtx); err != nil {
			// 记录事务回滚日志（用于调试）
			logger.Error("Transaction rolled back",
				zap.Error(err),
				zap.String("operation", "WithTransaction"))
			return err // GORM会自动回滚
		}

		// 记录事务提交日志（用于调试）
		logger.Debug("Transaction committed successfully")
		return nil // GORM会自动提交
	})
}

// WithTransactionResult 在事务中执行业务逻辑并返回结果
func (tm *TransactionManager) WithTransactionResult(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var result interface{}
	var resultErr error

	// 使用GORM事务
	err := tm.db.Transaction(func(tx *gorm.DB) error {
		// 将事务实例放入上下文
		txCtx := context.WithValue(ctx, shared.TransactionKey, tx)

		// 执行业务逻辑并获取结果
		result, resultErr = fn(txCtx)
		return resultErr // 如果有错误，GORM会自动回滚
	})

	if err != nil {
		logger.Error("Transaction with result rolled back",
			zap.Error(err),
			zap.String("operation", "WithTransactionResult"))
		return nil, err
	}

	logger.Debug("Transaction with result committed successfully")
	return result, nil
}

// 为什么这样实现？
//
// 1. 依赖GORM的Transaction方法：
//    - GORM已经处理了所有事务细节（开启/提交/回滚）
//    - 我们不需要重复造轮子
//    - 如果fn返回错误，GORM自动回滚；否则自动提交
//
// 2. 通过Context传递事务：
//    - Repository可以通过context.Value获取事务实例
//    - 符合Go的context使用惯例
//    - 不需要修改Repository接口
//
// 3. 日志记录：
//    - Debug级别记录成功提交（开发时有用）
//    - Error级别记录回滚（生产环境重要）
//    - 包含操作类型，便于调试
//
// 4. 泛型支持：
//    - WithTransactionResult支持返回任意类型
//    - 避免了interface{}和类型转换
//    - 编译时类型检查
