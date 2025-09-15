package mysql

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/domainerror"
	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/taskflow/internal/domain/event"
	"gorm.io/gorm"
)

// roleRepository 角色仓储实现
type roleRepository struct {
	*BaseRepository
	event.TransactionManager
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		BaseRepository:     NewBaseRepository(db),
		TransactionManager: NewTransactionManager(db),
	}
}

// Save 保存角色
func (r *roleRepository) Save(ctx context.Context, role *aggregate.Role) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		model := &Role{
			ID:          string(role.ID),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: &role.Description,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
		}

		if err := tx.WithContext(ctx).Save(model).Error; err != nil {
			return fmt.Errorf("failed to save role: %w", err)
		}

		return nil
	})
}

// FindByID 根据ID查找角色
func (r *roleRepository) FindByID(ctx context.Context, id valueobject.RoleID) (*aggregate.Role, error) {
	var model Role

	err := r.db.WithContext(ctx).Where("id = ?", string(id)).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerror.NewDomainError(domainerror.ErrRoleNotFound, "role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	return r.modelToAggregate(&model), nil
}

// FindByName 根据名称查找角色
func (r *roleRepository) FindByName(ctx context.Context, name string) (*aggregate.Role, error) {
	var model Role

	err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerror.NewDomainError(domainerror.ErrRoleNotFound, "role not found")
		}
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	return r.modelToAggregate(&model), nil
}

// FindAll 查找所有角色
func (r *roleRepository) FindAll(ctx context.Context) ([]*aggregate.Role, error) {
	var models []Role

	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}

	roles := make([]*aggregate.Role, len(models))
	for i, model := range models {
		roles[i] = r.modelToAggregate(&model)
	}

	return roles, nil
}

// Delete 删除角色
func (r *roleRepository) Delete(ctx context.Context, id valueobject.RoleID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		result := tx.WithContext(txCtx).Delete(&Role{}, "id = ?", string(id))
		if result.Error != nil {
			return fmt.Errorf("failed to delete role: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return domainerror.NewDomainError(domainerror.ErrRoleNotFound, "role not found")
		}

		return nil
	})
}

// AddPermissionToRole 为角色添加权限
func (r *roleRepository) AddPermissionToRole(ctx context.Context, roleID valueobject.RoleID, permissionID valueobject.PermissionID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		// 检查关联是否已存在
		tx := r.GetDB(txCtx)
		var count int64
		err := tx.WithContext(txCtx).
			Table("role_permissions").
			Where("role_id = ? AND permission_id = ?", string(roleID), string(permissionID)).
			Count(&count).Error

		if err != nil {
			return fmt.Errorf("failed to check existing role permission: %w", err)
		}

		if count > 0 {
			return nil // 已存在，不需要重复添加
		}

		// 创建关联
		rolePermission := &RolePermission{
			RoleID:       string(roleID),
			PermissionID: string(permissionID),
		}

		if err := tx.WithContext(txCtx).Create(rolePermission).Error; err != nil {
			return fmt.Errorf("failed to add permission to role: %w", err)
		}

		return nil
	})
}

// RemovePermissionFromRole 从角色移除权限
func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID valueobject.RoleID, permissionID valueobject.PermissionID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		result := tx.WithContext(txCtx).
			Where("role_id = ? AND permission_id = ?", string(roleID), string(permissionID)).
			Delete(&RolePermission{})

		if result.Error != nil {
			return fmt.Errorf("failed to remove permission from role: %w", result.Error)
		}

		return nil
	})
}

// FindPermissionsByRole 查找角色的所有权限
func (r *roleRepository) FindPermissionsByRole(ctx context.Context, roleID valueobject.RoleID) ([]*aggregate.Permission, error) {
	var models []Permission

	err := r.db.WithContext(ctx).
		Select("permissions.*").
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", string(roleID)).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find permissions by role: %w", err)
	}

	permissions := make([]*aggregate.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.modelToPermissionAggregate(&model)
	}

	return permissions, nil
}

// modelToAggregate 将数据模型转换为角色聚合根
func (r *roleRepository) modelToAggregate(model *Role) *aggregate.Role {
	description := ""
	if model.Description != nil {
		description = *model.Description
	}

	return aggregate.NewRole(
		valueobject.RoleID(model.ID),
		model.Name,
		model.DisplayName,
		description,
		model.IsSystem,
	)
}

// modelToPermissionAggregate 将权限数据模型转换为权限聚合根
func (r *roleRepository) modelToPermissionAggregate(model *Permission) *aggregate.Permission {
	description := ""
	if model.Description != nil {
		description = *model.Description
	}

	return aggregate.NewPermission(
		valueobject.PermissionID(model.ID),
		model.Name,
		valueobject.ResourceType(model.Resource),
		valueobject.ActionType(model.Action),
		description,
	)
}
