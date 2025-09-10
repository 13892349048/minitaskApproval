package shared

import "context"

// PermissionService 权限服务接口
type PermissionService interface {
	// CanUserPerformAction 检查用户是否可以执行某个操作
	CanUserPerformAction(ctx context.Context, userID, resource, action string, context map[string]interface{}) (bool, error)

	// GetUserPermissions 获取用户的所有权限
	GetUserPermissions(ctx context.Context, userID string) ([]Permission, error)

	// HasRole 检查用户是否具有某个角色
	HasRole(ctx context.Context, userID, roleName string) (bool, error)

	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
}

// Permission 权限结构
type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"` // 资源类型：project, task, user等
	Action      string `json:"action"`   // 操作类型：create, read, update, delete等
	Description string `json:"description"`
}

// Role 角色结构
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
	IsSystem    bool         `json:"is_system"` // 系统角色不可删除
}

// PolicyRule ABAC策略规则
type PolicyRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Effect      PolicyEffect           `json:"effect"`     // allow, deny
	Conditions  map[string]interface{} `json:"conditions"` // 条件表达式
	Priority    int                    `json:"priority"`   // 优先级，数值越大优先级越高
	IsActive    bool                   `json:"is_active"`
}

// PolicyEffect 策略效果
type PolicyEffect string

const (
	PolicyEffectAllow PolicyEffect = "allow"
	PolicyEffectDeny  PolicyEffect = "deny"
)

// EvaluationContext 权限评估上下文
type EvaluationContext struct {
	UserID      string                 `json:"user_id"`
	UserRoles   []string               `json:"user_roles"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	ResourceCtx map[string]interface{} `json:"resource_context"` // 资源上下文，如项目ID、任务ID等
	Environment map[string]interface{} `json:"environment"`      // 环境信息，如时间、IP等
}

// 预定义的系统角色（与业务需求和数据库初始化数据保持一致）
const (
	RoleSuperAdmin     = "super_admin"     // 超级管理员（系统管理）
	RoleAdmin          = "admin"           // 系统管理员
	RoleProjectOwner   = "project_owner"   // 项目所有者（大领导）
	RoleProjectManager = "project_manager" // 项目经理（项目领导）
	RoleTeamLeader     = "team_leader"     // 团队负责人（任务负责人）
	RoleEmployee       = "employee"        // 普通员工
)

// 预定义的资源类型
const (
	ResourceProject = "project"
	ResourceTask    = "task"
	ResourceUser    = "user"
	ResourceFile    = "file"
)

// 预定义的操作类型
const (
	ActionCreate  = "create"
	ActionRead    = "read"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionAssign  = "assign"
	ActionApprove = "approve"
	ActionExecute = "execute"
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
//    - 优先级机制解决权限冲突
