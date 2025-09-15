package service

import (
	"context"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

// UserDomainService 用户领域服务接口
// 处理复杂的业务规则和跨聚合的业务逻辑
type UserDomainService interface {
	// 业务规则验证
	ValidateUserCreation(ctx context.Context, email, username string) error
	ValidateManagerAssignment(ctx context.Context, user *aggregate.User, managerID valueobject.UserID) error
	ValidateRoleChange(ctx context.Context, user *aggregate.User, newRole valueobject.UserRole, changedBy valueobject.UserID) error

	// 复杂业务逻辑
	TransferUserDepartment(ctx context.Context, user *aggregate.User, newDepartmentID string, newManagerID valueobject.UserID) error
	DeactivateUserAndTransferTasks(ctx context.Context, user *aggregate.User, deactivatedBy valueobject.UserID) error
}

// 简化的接口定义 - 只保留必要的抽象

// PasswordHasher 密码哈希接口 - 在Infrastructure层实现
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) bool
}

// UserValidator 用户验证器接口 - 在Infrastructure层实现
type UserValidator interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
	ValidateName(name string) error
}
