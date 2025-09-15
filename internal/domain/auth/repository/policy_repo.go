package repository

import (
	"context"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/valueobject"
)

// PolicyRepository 策略仓储接口
type PolicyRepository interface {
	Save(ctx context.Context, policy *aggregate.Policy) error
	FindByID(ctx context.Context, id valueobject.PolicyID) (*aggregate.Policy, error)
	FindByResourceAndAction(ctx context.Context, resource valueobject.ResourceType, action valueobject.ActionType) ([]*aggregate.Policy, error)
	FindAllActive(ctx context.Context) ([]*aggregate.Policy, error)
	Delete(ctx context.Context, id valueobject.PolicyID) error
	CountByResource(ctx context.Context, resource valueobject.ResourceType) (int64, error)
}
