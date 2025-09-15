package repository

import (
	"context"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

// UserRepository 用户仓储接口 - 定义在Domain层
type UserRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, user *aggregate.User) error
	FindByID(ctx context.Context, id string) (*aggregate.User, error)
	FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
	FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *aggregate.User) error

	// 查询方法
	FindByRole(ctx context.Context, role valueobject.UserRole) ([]*aggregate.User, error)
	FindByStatus(ctx context.Context, status valueobject.UserStatus) ([]*aggregate.User, error)
	FindByDepartment(ctx context.Context, departmentID string) ([]*aggregate.User, error)
	FindByManager(ctx context.Context, managerID valueobject.UserID) ([]*aggregate.User, error)

	// 复杂查询
	SearchUsers(ctx context.Context, criteria valueobject.UserSearchCriteria) ([]*aggregate.User, int, error)
	FindUsersByRole(ctx context.Context, roleName string, limit, offset int) ([]*aggregate.User, int, error)

	// 统计查询
	CountByStatus(ctx context.Context, status valueobject.UserStatus) (int, error)
	CountByDepartment(ctx context.Context, department string) (int, error)
}
