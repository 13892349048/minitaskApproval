package repository

import (
	"context"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/valueobject"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Save(ctx context.Context, role *aggregate.Role) error
	FindByID(ctx context.Context, id valueobject.RoleID) (*aggregate.Role, error)
	FindByName(ctx context.Context, name string) (*aggregate.Role, error)
	FindAll(ctx context.Context) ([]*aggregate.Role, error)
	Delete(ctx context.Context, id valueobject.RoleID) error

	// 角色权限关联
	AddPermissionToRole(ctx context.Context, roleID valueobject.RoleID, permissionID valueobject.PermissionID) error
	RemovePermissionFromRole(ctx context.Context, roleID valueobject.RoleID, permissionID valueobject.PermissionID) error
	FindPermissionsByRole(ctx context.Context, roleID valueobject.RoleID) ([]*aggregate.Permission, error)
}

// UserRoleRepository 用户角色关联仓储接口
type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID string, roleID valueobject.RoleID) error
	RevokeRole(ctx context.Context, userID string, roleID valueobject.RoleID) error
	FindRolesByUser(ctx context.Context, userID string) ([]*aggregate.Role, error)
	FindUsersByRole(ctx context.Context, roleID valueobject.RoleID) ([]string, error)
	HasRole(ctx context.Context, userID string, roleID valueobject.RoleID) (bool, error)
}
