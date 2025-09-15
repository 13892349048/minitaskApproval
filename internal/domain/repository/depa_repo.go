package repository

import (
	"context"

	"github.com/taskflow/internal/domain/valueobject"
)

// DepartmentRepository 部门仓储接口
type DepartmentRepository interface {
	FindByID(ctx context.Context, id valueobject.DepartmentID) (*valueobject.DepartmentInfo, error)
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*valueobject.DepartmentInfo, error)
	IsActive(ctx context.Context, id valueobject.DepartmentID) (bool, error)
}
