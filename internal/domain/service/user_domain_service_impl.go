package service

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
)

// UserDomainServiceImpl 用户领域服务实现
type UserDomainServiceImpl struct {
	userRepo repository.UserRepository
	taskRepo repository.TaskRepository
	// TODO: 添加项目仓储接口
	projectRepo repository.ProjectRepository
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	projectRepo repository.ProjectRepository,
) UserDomainService {
	return &UserDomainServiceImpl{
		userRepo:    userRepo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

// ValidateUserCreation 验证用户创建
func (s *UserDomainServiceImpl) ValidateUserCreation(ctx context.Context, email, username string) error {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("email already exists: %s", email)
	}

	// TODO: 添加用户名唯一性检查
	// 当前UserRepository接口中没有FindByUsername方法
	// 需要在repository接口中添加此方法

	return nil
}

// ValidateManagerAssignment 验证管理者分配
func (s *UserDomainServiceImpl) ValidateManagerAssignment(ctx context.Context, user *aggregate.User, managerID valueobject.UserID) error {
	// 不能将自己设为管理者
	if user.ID == managerID {
		return fmt.Errorf("user cannot be their own manager")
	}

	// 检查管理者是否存在
	manager, err := s.userRepo.FindByID(ctx, string(managerID))
	if err != nil {
		return fmt.Errorf("manager not found: %w", err)
	}

	// 检查管理者是否激活
	if manager.Status != valueobject.UserStatusActive {
		return fmt.Errorf("manager must be active")
	}

	return nil
}

// ValidateRoleChange 验证角色变更
func (s *UserDomainServiceImpl) ValidateRoleChange(ctx context.Context, user *aggregate.User, newRole valueobject.UserRole, changedBy valueobject.UserID) error {
	// 检查操作者权限
	operator, err := s.userRepo.FindByID(ctx, string(changedBy))
	if err != nil {
		return fmt.Errorf("operator not found: %w", err)
	}

	// 只有管理员可以变更角色
	if operator.Role != valueobject.UserRoleAdmin && operator.Role != valueobject.UserRoleSuperAdmin {
		return fmt.Errorf("insufficient permission to change user role")
	}

	// 超级管理员角色只能由超级管理员分配
	if newRole == valueobject.UserRoleSuperAdmin && operator.Role != valueobject.UserRoleSuperAdmin {
		return fmt.Errorf("only super admin can assign super admin role")
	}

	return nil
}

// TransferUserDepartment 转移用户部门
func (s *UserDomainServiceImpl) TransferUserDepartment(ctx context.Context, user *aggregate.User, newDepartmentID string, newManagerID valueobject.UserID) error {
	// 验证新管理者
	if err := s.ValidateManagerAssignment(ctx, user, newManagerID); err != nil {
		return fmt.Errorf("invalid manager assignment: %w", err)
	}

	// TODO: 验证部门是否存在
	// 需要添加部门仓储接口来验证部门存在性
	if newDepartmentID != "" {
		// 暂时跳过部门验证，需要实现部门仓储
	}

	// TODO: 检查是否会造成循环管理关系
	// 需要实现层级验证逻辑

	return nil
}

// DeactivateUserAndTransferTasks 停用用户并转移任务
func (s *UserDomainServiceImpl) DeactivateUserAndTransferTasks(ctx context.Context, user *aggregate.User, deactivatedBy valueobject.UserID) error {
	// 检查操作者权限
	operator, err := s.userRepo.FindByID(ctx, string(deactivatedBy))
	if err != nil {
		return fmt.Errorf("operator not found: %w", err)
	}

	// 只有管理员可以停用用户
	if operator.Role != valueobject.UserRoleAdmin && operator.Role != valueobject.UserRoleSuperAdmin {
		return fmt.Errorf("insufficient permission to deactivate user")
	}

	// 不能停用超级管理员
	if user.Role == valueobject.UserRoleSuperAdmin {
		return fmt.Errorf("cannot deactivate super admin")
	}

	// 获取用户的所有活跃任务
	activeTasks, err := s.taskRepo.FindByResponsible(ctx, valueobject.UserID(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user tasks: %w", err)
	}

	// 如果有活跃任务，需要先转移
	if len(activeTasks) > 0 {
		return fmt.Errorf("user has %d active tasks that must be transferred first", len(activeTasks))
	}

	// TODO: 获取用户管理的项目
	// 需要添加项目仓储接口来检查用户管理的项目
	// managedProjects, err := s.projectRepo.FindByManager(ctx, valueobject.UserID(user.ID))
	// if err != nil {
	//     return fmt.Errorf("failed to get managed projects: %w", err)
	// }
	// if len(managedProjects) > 0 {
	//     return fmt.Errorf("user manages %d projects that must be transferred first", len(managedProjects))
	// }

	return nil
}
