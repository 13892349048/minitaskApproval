package mysql

import (
	"context"

	"github.com/taskflow/internal/domain/shared"
	"gorm.io/gorm"
)

// BaseRepository 基础仓储，提供事务支持
// 所有具体的Repository都应该嵌入这个结构体
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB 从上下文获取数据库连接（自动支持事务）
func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	// 尝试从上下文获取事务实例
	if tx, ok := ctx.Value(shared.TransactionKey).(*gorm.DB); ok {
		// 如果在事务中，使用事务连接
		return tx
	}
	// 如果不在事务中，使用普通连接
	return r.db
}

// 为什么这样设计？
//
// 1. 自动事务检测：
//    - Repository不需要知道自己是否在事务中
//    - GetDB()自动返回正确的数据库连接
//    - 事务中用事务连接，非事务中用普通连接
//
// 2. 零侵入性：
//    - Repository的业务代码不需要修改
//    - 只需要用GetDB(ctx)替换直接使用r.db
//    - 事务管理完全透明
//
// 3. 继承模式：
//    - 所有Repository都嵌入BaseRepository
//    - 自动获得事务支持能力
//    - 代码复用，减少重复
