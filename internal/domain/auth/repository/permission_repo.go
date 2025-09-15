package repository

import (
	"context"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/valueobject"
)

// 为什么这样设计？
//
// 1. RBAC基础：
//    - Role：角色定义，包含权限集合
//    - Permission：具体的权限，定义资源和操作
//    - 用户通过角色获得基础权限
//
// 2. ABAC扩展：
//    - PolicyRule：基于属性的策略规则
//    - EvaluationContext：权限评估的上下文信息
//    - 支持复杂的条件判断（如：只能管理自己的项目）
//
// 3. 混合模式优势：
//    - RBAC处理基础权限（简单、高效）
//    - ABAC处理复杂权限（灵活、精确）

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	Save(ctx context.Context, permission *aggregate.Permission) error
	FindByID(ctx context.Context, id valueobject.PermissionID) (*aggregate.Permission, error)
	FindByResourceAndAction(ctx context.Context, resource valueobject.ResourceType, action valueobject.ActionType) (*aggregate.Permission, error)
	FindAll(ctx context.Context) ([]*aggregate.Permission, error)
	Delete(ctx context.Context, id valueobject.PermissionID) error
}
