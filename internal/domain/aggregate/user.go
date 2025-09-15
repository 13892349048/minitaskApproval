package aggregate

import (
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/valueobject"
)

// User 用户聚合根
type User struct {
	ID           valueobject.UserID     `json:"id"`
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	FullName     string                 `json:"full_name"`
	PasswordHash string                 `json:"-"` // 密码哈希，不序列化
	Role         valueobject.UserRole   `json:"role"`
	Status       valueobject.UserStatus `json:"status"`
	DepartmentID *string                `json:"department_id,omitempty"`
	ManagerID    *valueobject.UserID    `json:"manager_id,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`

	// 领域事件
	events []event.DomainEvent
}

// NewUser 创建新用户 - Domain层工厂方法
func NewUser(
	id valueobject.UserID,
	username, email, fullName, passwordHash string,
	role valueobject.UserRole,
) *User {
	now := time.Now()

	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		FullName:     fullName,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       valueobject.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// 业务方法 - 字段现在是公共的，不需要getter方法

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(fullName, email string) error {
	if fullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	u.FullName = fullName
	u.Email = email
	u.UpdatedAt = time.Now()

	return nil
}

// ChangeRole 更改用户角色
func (u *User) ChangeRole(newRole valueobject.UserRole) {
	u.Role = newRole
	u.UpdatedAt = time.Now()
}

// AssignToDepartment 分配到部门
func (u *User) AssignToDepartment(departmentID string) {
	u.DepartmentID = &departmentID
	u.UpdatedAt = time.Now()
}

// AssignManager 分配经理
func (u *User) AssignManager(managerID valueobject.UserID) {
	u.ManagerID = &managerID
	u.UpdatedAt = time.Now()
}

// Activate 激活用户
func (u *User) Activate() {
	u.Status = valueobject.UserStatusActive
	u.UpdatedAt = time.Now()
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.Status = valueobject.UserStatusInactive
	u.UpdatedAt = time.Now()
}

// Suspend 暂停用户
func (u *User) Suspend() {
	u.Status = valueobject.UserStatusSuspended
	u.UpdatedAt = time.Now()
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == valueobject.UserStatusActive && u.DeletedAt == nil
}

// SoftDelete 软删除用户
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.Status = valueobject.UserStatusInactive
	u.UpdatedAt = now
}

// ChangePassword 更改密码 - 需要外部密码哈希服务
func (u *User) ChangePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now()
}

// HasPermission 检查用户权限 - 基础权限检查
func (u *User) HasPermission(permission string) bool {
	// 基于角色的基础权限检查
	switch u.Role {
	case valueobject.UserRoleAdmin:
		return true // 管理员拥有所有权限
	case valueobject.UserRoleDirector:
		return permission != "system:admin" // 总监除了系统管理外都有权限
	case valueobject.UserRoleManager:
		return permission == "user:read" || permission == "user:update" || permission == "task:manage"
	case valueobject.UserRoleEmployee:
		return permission == "user:read" || permission == "task:read"
	default:
		return false
	}
}

// CanManageUser 检查是否可以管理指定用户
func (u *User) CanManageUser(targetUser *User) bool {
	if !u.IsActive() {
		return false
	}

	// 管理员可以管理所有用户
	if u.Role == valueobject.UserRoleAdmin {
		return true
	}

	// 总监可以管理经理和员工
	if u.Role == valueobject.UserRoleDirector {
		return targetUser.Role == valueobject.UserRoleManager || targetUser.Role == valueobject.UserRoleEmployee
	}

	// 经理可以管理直接下属
	if u.Role == valueobject.UserRoleManager && targetUser.ManagerID != nil {
		return *targetUser.ManagerID == u.ID
	}

	return false
}

// AddEvent 添加领域事件
func (u *User) AddEvent(event event.DomainEvent) {
	u.events = append(u.events, event)
}

// GetEvents 获取领域事件
func (u *User) GetEvents() []event.DomainEvent {
	return u.events
}

// ClearEvents 清空领域事件
func (u *User) ClearEvents() {
	u.events = nil
}
