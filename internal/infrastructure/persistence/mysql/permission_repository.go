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

// permissionRepository 权限仓储实现
type permissionRepository struct {
	*BaseRepository
	event.TransactionManager
}

// NewPermissionRepository 创建权限仓储
func NewPermissionRepository(db *gorm.DB) *permissionRepository {
	return &permissionRepository{
		BaseRepository:     NewBaseRepository(db),
		TransactionManager: NewTransactionManager(db),
	}
}

// Save 保存权限
func (r *permissionRepository) Save(ctx context.Context, permission *aggregate.Permission) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		db := r.GetDB(txCtx)
		model := &Permission{
			ID:          string(permission.ID),
			Name:        permission.Name,
			Resource:    string(permission.Resource),
			Action:      string(permission.Action),
			Description: &permission.Description,
			CreatedAt:   permission.CreatedAt,
			UpdatedAt:   permission.UpdatedAt,
		}

		result := db.WithContext(txCtx).Save(model)
		if result.Error != nil {
			return fmt.Errorf("failed to save permission: %w", result.Error)
		}

		return nil
	})
}

// FindByID 根据ID查找权限
func (r *permissionRepository) FindByID(ctx context.Context, id valueobject.PermissionID) (*aggregate.Permission, error) {
	var model Permission
	db := r.GetDB(ctx)

	err := db.WithContext(ctx).Where("id = ?", string(id)).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerror.NewDomainError(domainerror.ErrPermissionNotFound, "permission not found")
		}
		return nil, fmt.Errorf("failed to find permission: %w", err)
	}

	return r.modelToAggregate(&model), nil
}

// FindByResourceAndAction 根据资源和操作查找权限
func (r *permissionRepository) FindByResourceAndAction(ctx context.Context, resource valueobject.ResourceType, action valueobject.ActionType) (*aggregate.Permission, error) {
	var model Permission
	db := r.GetDB(ctx)

	err := db.WithContext(ctx).
		Where("resource = ? AND action = ?", string(resource), string(action)).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerror.NewDomainError(domainerror.ErrPermissionNotFound, "permission not found")
		}
		return nil, fmt.Errorf("failed to find permission: %w", err)
	}

	return r.modelToAggregate(&model), nil
}

// FindAll 查找所有权限
func (r *permissionRepository) FindAll(ctx context.Context) ([]*aggregate.Permission, error) {
	var models []Permission
	db := r.GetDB(ctx)

	err := db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find permissions: %w", err)
	}

	permissions := make([]*aggregate.Permission, len(models))
	for i, model := range models {
		permissions[i] = r.modelToAggregate(&model)
	}

	return permissions, nil
}

// Delete 删除权限
func (r *permissionRepository) Delete(ctx context.Context, id valueobject.PermissionID) error {
	db := r.GetDB(ctx)

	result := db.WithContext(ctx).Delete(&Permission{}, "id = ?", string(id))
	if result.Error != nil {
		return fmt.Errorf("failed to delete permission: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}

	return nil
}

// modelToAggregate 将数据模型转换为聚合根
func (r *permissionRepository) modelToAggregate(model *Permission) *aggregate.Permission {
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
