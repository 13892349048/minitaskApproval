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

// userRoleRepository 用户角色关联仓储实现
type userRoleRepository struct {
	*BaseRepository
	event.TransactionManager
}

// NewUserRoleRepository 创建用户角色关联仓储
func NewUserRoleRepository(db *gorm.DB) *userRoleRepository {
	return &userRoleRepository{
		BaseRepository:     NewBaseRepository(db),
		TransactionManager: NewTransactionManager(db),
	}
}

// AssignRole 为用户分配角色
func (r *userRoleRepository) AssignRole(ctx context.Context, userID string, roleID valueobject.RoleID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		// 检查是否已经分配
		tx := r.GetDB(txCtx)
		var count int64
		err := tx.WithContext(ctx).
			Model(&UserRole{}).
			Where("user_id = ? AND role_id = ?", userID, string(roleID)).
			Count(&count).Error

		if err != nil {
			return fmt.Errorf("failed to check existing user role: %w", err)
		}

		if count > 0 {
			return domainerror.NewDomainError(domainerror.ErrRoleAlreadyAssigned, "role already assigned to user")
		}

		// 创建用户角色关联
		userRole := &UserRole{
			UserID: userID,
			RoleID: string(roleID),
		}

		if err := tx.WithContext(ctx).Create(userRole).Error; err != nil {
			return fmt.Errorf("failed to assign role to user: %w", err)
		}

		return nil
	})
}

// RevokeRole 撤销用户角色
func (r *userRoleRepository) RevokeRole(ctx context.Context, userID string, roleID valueobject.RoleID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		result := tx.WithContext(ctx).
			Where("user_id = ? AND role_id = ?", userID, string(roleID)).
			Delete(&UserRole{})

		if result.Error != nil {
			return fmt.Errorf("failed to revoke role from user: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return domainerror.NewDomainError(domainerror.ErrRoleNotAssigned, "role not assigned to user")
		}

		return nil
	})
}

// FindRolesByUser 查找用户的所有角色
func (r *userRoleRepository) FindRolesByUser(ctx context.Context, userID string) ([]*aggregate.Role, error) {
	var models []Role

	err := r.db.WithContext(ctx).
		Select("roles.*").
		Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find roles by user: %w", err)
	}

	roles := make([]*aggregate.Role, len(models))
	for i, model := range models {
		roles[i] = r.modelToRoleAggregate(&model)
	}

	return roles, nil
}

// FindUsersByRole 查找拥有某角色的所有用户
func (r *userRoleRepository) FindUsersByRole(ctx context.Context, roleID valueobject.RoleID) ([]string, error) {
	var userIDs []string

	err := r.db.WithContext(ctx).
		Model(&UserRole{}).
		Select("user_id").
		Where("role_id = ?", string(roleID)).
		Pluck("user_id", &userIDs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	return userIDs, nil
}

// HasRole 检查用户是否拥有某角色
func (r *userRoleRepository) HasRole(ctx context.Context, userID string, roleID valueobject.RoleID) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, string(roleID)).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return count > 0, nil
}

// modelToRoleAggregate 将角色数据模型转换为角色聚合根
func (r *userRoleRepository) modelToRoleAggregate(model *Role) *aggregate.Role {
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
