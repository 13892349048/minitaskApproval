package aggregate

import (
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
)

// Policy ABAC策略聚合根
type Policy struct {
	ID          valueobject.PolicyID
	Name        string
	Description string
	Resource    valueobject.ResourceType
	Action      valueobject.ActionType
	Effect      valueobject.PolicyEffect
	Conditions  valueobject.PolicyConditions
	Priority    int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPolicy 创建策略
func NewPolicy(id valueobject.PolicyID, name, description string, resource valueobject.ResourceType, action valueobject.ActionType, effect valueobject.PolicyEffect, conditions valueobject.PolicyConditions, priority int) *Policy {
	return &Policy{
		ID:          id,
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
		Effect:      effect,
		Conditions:  conditions,
		Priority:    priority,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdatePolicy 更新策略
func (p *Policy) UpdatePolicy(name, description string, effect valueobject.PolicyEffect, conditions valueobject.PolicyConditions, priority int) {
	p.Name = name
	p.Description = description
	p.Effect = effect
	p.Conditions = conditions
	p.Priority = priority
	p.UpdatedAt = time.Now()
}

// Activate 激活策略
func (p *Policy) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate 停用策略
func (p *Policy) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// Matches 检查策略是否匹配资源和操作
func (p *Policy) Matches(resource valueobject.ResourceType, action valueobject.ActionType) bool {
	return p.IsActive && p.Resource == resource && p.Action == action
}
