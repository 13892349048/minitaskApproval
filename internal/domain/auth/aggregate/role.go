package aggregate

import (
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/taskflow/pkg/errors"
)

// Role 角色聚合根
type Role struct {
	ID          valueobject.RoleID
	Name        string
	DisplayName string
	Description string
	Permissions []valueobject.PermissionID
	IsSystem    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewRole 创建角色
func NewRole(id valueobject.RoleID, name, displayName, description string, isSystem bool) *Role {
	return &Role{
		ID:          id,
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Permissions: make([]valueobject.PermissionID, 0),
		IsSystem:    isSystem,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// AddPermission 添加权限
func (r *Role) AddPermission(permissionID valueobject.PermissionID) error {
	// 检查权限是否已存在
	for _, pid := range r.Permissions {
		if pid == permissionID {
			return errors.NewDomainError("权限已存在", "PERMISSION_ALREADY_EXISTS")
		}
	}

	r.Permissions = append(r.Permissions, permissionID)
	r.UpdatedAt = time.Now()
	return nil
}

// RemovePermission 移除权限
func (r *Role) RemovePermission(permissionID valueobject.PermissionID) error {
	if r.IsSystem {
		return errors.NewDomainError("系统角色不能修改权限", "SYSTEM_ROLE_IMMUTABLE")
	}

	for i, pid := range r.Permissions {
		if pid == permissionID {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.NewDomainError("权限不存在", "PERMISSION_NOT_FOUND")
}

// HasPermission 检查是否拥有权限
func (r *Role) HasPermission(permissionID valueobject.PermissionID) bool {
	for _, pid := range r.Permissions {
		if pid == permissionID {
			return true
		}
	}
	return false
}

// UpdateInfo 更新角色信息
func (r *Role) UpdateInfo(displayName, description string) error {
	if r.IsSystem {
		return errors.NewDomainError("系统角色不能修改", "SYSTEM_ROLE_IMMUTABLE")
	}

	r.DisplayName = displayName
	r.Description = description
	r.UpdatedAt = time.Now()
	return nil
}
