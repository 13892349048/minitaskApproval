package aggregate

import (
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
)

// Permission 权限聚合根
type Permission struct {
	ID          valueobject.PermissionID
	Name        string
	Resource    valueobject.ResourceType
	Action      valueobject.ActionType
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPermission 创建权限
func NewPermission(id valueobject.PermissionID, name string, resource valueobject.ResourceType, action valueobject.ActionType, description string) *Permission {
	return &Permission{
		ID:          id,
		Name:        name,
		Resource:    resource,
		Action:      action,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateDescription 更新描述
func (p *Permission) UpdateDescription(description string) {
	p.Description = description
	p.UpdatedAt = time.Now()
}

// Matches 检查权限是否匹配资源和操作
func (p *Permission) Matches(resource valueobject.ResourceType, action valueobject.ActionType) bool {
	return p.Resource == resource && p.Action == action
}
