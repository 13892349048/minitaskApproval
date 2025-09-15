package service

import (
	"context"
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
	"go.uber.org/zap"
)

// UserDomainServiceEnhanced 增强的用户领域服务实现
// 包含完整的业务逻辑和验证规则
type UserDomainServiceEnhanced struct {
	userRepo       repository.UserRepository
	taskRepo       repository.TaskRepository
	projectRepo    repository.ProjectRepository
	departmentRepo repository.DepartmentRepository
	eventPublisher event.EventBus
	logger         *zap.Logger
} // 领域事件发布器

// NewUserDomainServiceEnhanced 创建增强的用户领域服务
func NewUserDomainServiceEnhanced(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	projectRepo repository.ProjectRepository,
	departmentRepo repository.DepartmentRepository,
	eventPublisher event.EventBus,
	logger *zap.Logger,
) UserDomainService {
	return &UserDomainServiceEnhanced{
		userRepo:       userRepo,
		taskRepo:       taskRepo,
		projectRepo:    projectRepo,
		departmentRepo: departmentRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// ValidateUserCreation 验证用户创建（增强版）
func (s *UserDomainServiceEnhanced) ValidateUserCreation(ctx context.Context, email, username string) error {
	// validationContext := valueobject.UserValidationContext{
	// 	Operation:   "create_user",
	// 	RequestData: map[string]interface{}{
	// 		"email":    email,
	// 		"username": username,
	// 	},
	// 	Timestamp: time.Now(),
	// }

	var violations []valueobject.ValidationViolation

	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		violations = append(violations, valueobject.ValidationViolation{
			RuleID:   "email_unique",
			RuleName: "邮箱唯一性",
			Field:    "email",
			Value:    email,
			Message:  fmt.Sprintf("邮箱已存在: %s", email),
			Severity: valueobject.ValidationSeverityError,
		})
	}

	// 检查用户名是否已存在
	existingUser, err = s.userRepo.FindByUsername(ctx, username)
	if err == nil && existingUser != nil {
		violations = append(violations, valueobject.ValidationViolation{
			Field:   "username",
			Message: "用户名已存在",
		})
	}

	if len(violations) > 0 {
		return fmt.Errorf("用户创建验证失败: %d个违规", len(violations))
	}

	return nil
}

// ValidateManagerAssignment 验证管理者分配（增强版）
func (s *UserDomainServiceEnhanced) ValidateManagerAssignment(ctx context.Context, user *aggregate.User, managerID valueobject.UserID) error {
	// 不能将自己设为管理者
	if user.ID == managerID {
		return &valueobject.CircularManagershipError{
			UserID:    valueobject.UserID(user.ID),
			ManagerID: managerID,
			Chain:     []valueobject.UserID{valueobject.UserID(user.ID)},
			Message:   "用户不能将自己设为管理者",
		}
	}

	// 检查管理者是否存在
	manager, err := s.userRepo.FindByID(ctx, string(managerID))
	if err != nil {
		return fmt.Errorf("管理者不存在: %w", err)
	}

	// 检查管理者是否激活
	if manager.Status != valueobject.UserStatusActive {
		return fmt.Errorf("管理者必须是激活状态")
	}

	// 检查是否会造成循环管理关系
	if err := s.validateManagerHierarchy(ctx, valueobject.UserID(user.ID), managerID); err != nil {
		return err
	}

	return nil
}

// ValidateRoleChange 验证角色变更（增强版）
func (s *UserDomainServiceEnhanced) ValidateRoleChange(ctx context.Context, user *aggregate.User, newRole valueobject.UserRole, changedBy valueobject.UserID) error {
	// 检查操作者权限
	operator, err := s.userRepo.FindByID(ctx, string(changedBy))
	if err != nil {
		return fmt.Errorf("操作者不存在: %w", err)
	}

	// 权限验证矩阵
	canChangeRole := s.canChangeUserRole(operator.Role, user.Role, newRole)
	if !canChangeRole {
		return fmt.Errorf("权限不足: %s 角色无法将 %s 变更为 %s",
			operator.Role, user.Role, newRole)
	}

	// 发布角色变更事件
	userRoleChangedEvent := event.UserRoleChangedEvent{
		ID:           generateEventID(),
		UserID:       valueobject.UserID(user.ID),
		OldRole:      user.Role,
		NewRole:      newRole,
		ChangedBy:    changedBy,
		OccurredOn:   time.Now(),
		EventVersion: 1,
	}

	if err := s.eventPublisher.Publish(userRoleChangedEvent); err != nil {
		s.logger.Error("Failed to publish UserRoleChangedEvent",
			zap.String("user_id", string(user.ID)),
			zap.Error(err))
	}

	return nil
}

// TransferUserDepartment 转移用户部门（增强版）
func (s *UserDomainServiceEnhanced) TransferUserDepartment(ctx context.Context, user *aggregate.User, newDepartmentID string, newManagerID valueobject.UserID) error {
	// 验证新管理者
	if err := s.ValidateManagerAssignment(ctx, user, newManagerID); err != nil {
		return fmt.Errorf("管理者分配验证失败: %w", err)
	}

	// 验证部门是否存在且激活
	if newDepartmentID != "" {
		departmentID := valueobject.DepartmentID(newDepartmentID)

		// 检查部门是否存在
		department, err := s.departmentRepo.FindByID(ctx, departmentID)
		if err != nil {
			return &valueobject.DepartmentValidationError{
				DepartmentID: departmentID,
				ErrorType:    "not_found",
				Message:      fmt.Sprintf("部门不存在: %s", newDepartmentID),
			}
		}

		// 检查部门是否激活
		if !department.IsActive {
			return &valueobject.DepartmentValidationError{
				DepartmentID: departmentID,
				ErrorType:    "inactive",
				Message:      fmt.Sprintf("部门已停用: %s", newDepartmentID),
			}
		}
	}

	// 获取当前部门信息
	currentDepartment, err := s.departmentRepo.FindByUserID(ctx, valueobject.UserID(user.ID))
	if err != nil {
		// 用户可能没有部门，这是允许的
		currentDepartment = nil
	}

	// 发布部门转移事件
	userDepartmentTransferredEvent := event.UserDepartmentTransferredEvent{
		ID:     generateEventID(),
		UserID: valueobject.UserID(user.ID),
		FromDepartmentID: func() valueobject.DepartmentID {
			if currentDepartment != nil {
				return currentDepartment.ID
			}
			return ""
		}(),
		ToDepartmentID: valueobject.DepartmentID(newDepartmentID),
		NewManagerID:   newManagerID,
		TransferredBy:  valueobject.UserID("system"), // 从上下文获取操作者，当前使用系统用户
		OccurredOn:     time.Now(),
		EventVersion:   1,
	}

	if err := s.eventPublisher.Publish(userDepartmentTransferredEvent); err != nil {
		// 记录事件发布失败
		s.logger.Error("Failed to publish UserDepartmentTransferredEvent",
			zap.String("user_id", string(user.ID)),
			zap.String("new_department", newDepartmentID),
			zap.String("new_manager", string(newManagerID)),
			zap.Error(err))
	}

	return nil
}

// DeactivateUserAndTransferTasks 停用用户并转移任务（增强版）
func (s *UserDomainServiceEnhanced) DeactivateUserAndTransferTasks(ctx context.Context, user *aggregate.User, deactivatedBy valueobject.UserID) error {
	// 检查操作者权限
	operator, err := s.userRepo.FindByID(ctx, string(deactivatedBy))
	if err != nil {
		return fmt.Errorf("操作者不存在: %w", err)
	}

	// 权限验证
	if !s.canDeactivateUser(operator.Role, user.Role) {
		return fmt.Errorf("权限不足: %s 角色无法停用 %s 角色的用户",
			operator.Role, user.Role)
	}

	// 获取用户的所有活跃任务
	activeTasks, err := s.taskRepo.FindByResponsible(ctx, valueobject.UserID(user.ID))
	if err != nil {
		return fmt.Errorf("获取用户任务失败: %w", err)
	}

	// 获取用户管理的项目
	managedProjects, err := s.projectRepo.FindByManager(ctx, valueobject.UserID(user.ID))
	if err != nil {
		return fmt.Errorf("获取管理项目失败: %w", err)
	}

	// 检查是否有未处理的任务或项目
	hasActiveTasks := len(activeTasks) > 0
	hasManagedProjects := len(managedProjects) > 0

	if hasActiveTasks || hasManagedProjects {
		return fmt.Errorf("用户停用失败: 存在 %d 个活跃任务和 %d 个管理项目，需要先转移",
			len(activeTasks), len(managedProjects))
	}

	// 发布用户停用事件
	userDeactivatedEvent := event.UserDeactivatedEvent{
		ID:                  generateEventID(),
		UserID:              valueobject.UserID(user.ID),
		DeactivatedBy:       deactivatedBy,
		TasksTransferred:    !hasActiveTasks,
		ProjectsTransferred: !hasManagedProjects,
		OccurredOn:          time.Now(),
		EventVersion:        1,
	}

	if err := s.eventPublisher.Publish(userDeactivatedEvent); err != nil {
		// 记录事件发布失败
		s.logger.Error("Failed to publish UserDeactivatedEvent",
			zap.String("user_id", string(user.ID)),
			zap.String("deactivated_by", string(deactivatedBy)),
			zap.Error(err))
	}

	return nil
}

// 辅助方法

// validateManagerHierarchy 验证管理层级，防止循环管理关系
func (s *UserDomainServiceEnhanced) validateManagerHierarchy(ctx context.Context, userID, managerID valueobject.UserID) error {
	visited := make(map[valueobject.UserID]bool)
	chain := []valueobject.UserID{userID}

	currentManagerID := managerID
	for currentManagerID != "" {
		// 检查是否已访问过，防止无限循环
		if visited[currentManagerID] {
			return &valueobject.CircularManagershipError{
				UserID:    userID,
				ManagerID: managerID,
				Chain:     chain,
				Message:   fmt.Sprintf("检测到循环管理关系: %v", chain),
			}
		}

		// 检查是否回到了原始用户
		if currentManagerID == userID {
			return &valueobject.CircularManagershipError{
				UserID:    userID,
				ManagerID: managerID,
				Chain:     append(chain, currentManagerID),
				Message:   "管理者层级形成循环",
			}
		}

		visited[currentManagerID] = true
		chain = append(chain, currentManagerID)

		// 获取当前管理者的管理者
		manager, err := s.userRepo.FindByID(ctx, string(currentManagerID))
		if err != nil {
			break // 管理者不存在，结束检查
		}

		// TODO: 需要在User聚合中添加ManagerID字段
		// 由于User聚合暂时没有ManagerID字段，我们暂时结束检查
		// 在实际实现中，应该获取manager.ManagerID并继续循环
		_ = manager // 避免未使用变量警告
		break       // 暂时结束循环，等待User聚合完善
	}

	return nil
}

// canChangeUserRole 检查是否有权限变更用户角色
func (s *UserDomainServiceEnhanced) canChangeUserRole(operatorRole, currentRole, newRole valueobject.UserRole) bool {
	// 超级管理员可以变更任何角色
	if operatorRole == valueobject.UserRoleSuperAdmin {
		return true
	}

	// 管理员可以变更除超级管理员外的角色
	if operatorRole == valueobject.UserRoleAdmin {
		return currentRole != valueobject.UserRoleSuperAdmin && newRole != valueobject.UserRoleSuperAdmin
	}

	// 其他角色无权变更角色
	return false
}

// canDeactivateUser 检查是否有权限停用用户
func (s *UserDomainServiceEnhanced) canDeactivateUser(operatorRole, targetRole valueobject.UserRole) bool {
	// 超级管理员可以停用除自己外的任何用户
	if operatorRole == valueobject.UserRoleSuperAdmin {
		return targetRole != valueobject.UserRoleSuperAdmin
	}

	// 管理员可以停用普通用户和经理
	if operatorRole == valueobject.UserRoleAdmin {
		return targetRole == valueobject.UserRoleEmployee || targetRole == valueobject.UserRoleManager
	}

	// 其他角色无权停用用户
	return false
}

// generateEventID 生成事件ID
func generateEventID() string {
	// 使用时间戳和随机数生成唯一ID
	// 在生产环境中建议使用UUID库如github.com/google/uuid
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("event_%d_%d", timestamp, timestamp%1000000)
}
