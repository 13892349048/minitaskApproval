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

// policyRepository 策略仓储实现
type policyRepository struct {
	*BaseRepository
	event.TransactionManager
}

// NewPolicyRepository 创建策略仓储
func NewPolicyRepository(db *gorm.DB) *policyRepository {
	return &policyRepository{
		BaseRepository:     NewBaseRepository(db),
		TransactionManager: NewTransactionManager(db),
	}
}

// Save 保存策略
func (r *policyRepository) Save(ctx context.Context, policy *aggregate.Policy) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		conditionsJSON, err := policy.Conditions.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to serialize policy conditions: %w", err)
		}

		model := &PermissionPolicy{
			ID:           string(policy.ID),
			Name:         policy.Name,
			Description:  &policy.Description,
			ResourceType: string(policy.Resource),
			Action:       string(policy.Action),
			Effect:       string(policy.Effect),
			Conditions:   conditionsJSON,
			Priority:     policy.Priority,
			IsActive:     policy.IsActive,
			CreatedAt:    policy.CreatedAt,
			UpdatedAt:    policy.UpdatedAt,
		}

		if err := tx.WithContext(txCtx).Save(model).Error; err != nil {
			return fmt.Errorf("failed to save policy: %w", err)
		}

		return nil
	})
}

// FindByID 根据ID查找策略
func (r *policyRepository) FindByID(ctx context.Context, id valueobject.PolicyID) (*aggregate.Policy, error) {
	var model PermissionPolicy

	err := r.db.WithContext(ctx).Where("id = ?", string(id)).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerror.NewDomainError(domainerror.ErrPolicyNotFound, "policy not found")
		}
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}

	return r.modelToAggregate(&model)
}

// FindByResourceAndAction 根据资源和操作查找策略
func (r *policyRepository) FindByResourceAndAction(ctx context.Context, resource valueobject.ResourceType, action valueobject.ActionType) ([]*aggregate.Policy, error) {
	var models []PermissionPolicy

	err := r.db.WithContext(ctx).
		Where("resource_type = ? AND action = ?", string(resource), string(action)).
		Order("priority DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find policies: %w", err)
	}

	policies := make([]*aggregate.Policy, 0, len(models))
	for _, model := range models {
		policy, err := r.modelToAggregate(&model)
		if err != nil {
			// 记录错误但继续处理其他策略
			continue
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

// FindAllActive 查找所有活跃策略
func (r *policyRepository) FindAllActive(ctx context.Context) ([]*aggregate.Policy, error) {
	var models []PermissionPolicy

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("priority DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find active policies: %w", err)
	}

	policies := make([]*aggregate.Policy, 0, len(models))
	for _, model := range models {
		policy, err := r.modelToAggregate(&model)
		if err != nil {
			// 记录错误但继续处理其他策略
			continue
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

// Delete 删除策略
func (r *policyRepository) Delete(ctx context.Context, id valueobject.PolicyID) error {
	return r.WithTransaction(ctx, func(txCtx context.Context) error {
		tx := r.GetDB(txCtx)
		result := tx.WithContext(txCtx).Delete(&PermissionPolicy{}, "id = ?", string(id))
		if result.Error != nil {
			return fmt.Errorf("failed to delete policy: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return domainerror.NewDomainError(domainerror.ErrPolicyNotFound, "policy not found")
		}

		return nil
	})
}

// CountByResource 统计资源的策略数量
func (r *policyRepository) CountByResource(ctx context.Context, resource valueobject.ResourceType) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&PermissionPolicy{}).
		Where("resource_type = ?", string(resource)).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count policies by resource: %w", err)
	}

	return count, nil
}

// modelToAggregate 将数据模型转换为策略聚合根
func (r *policyRepository) modelToAggregate(model *PermissionPolicy) (*aggregate.Policy, error) {
	var conditions valueobject.PolicyConditions
	if err := conditions.FromJSON(model.Conditions); err != nil {
		return nil, fmt.Errorf("failed to deserialize policy conditions: %w", err)
	}

	policy := aggregate.NewPolicy(
		valueobject.PolicyID(model.ID),
		model.Name,
		*model.Description,
		valueobject.ResourceType(model.ResourceType),
		valueobject.ActionType(model.Action),
		valueobject.PolicyEffect(model.Effect),
		conditions,
		model.Priority,
	)

	if !model.IsActive {
		policy.Deactivate()
	}

	return policy, nil
}
