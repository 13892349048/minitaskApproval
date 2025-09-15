package service

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/domainerror"
	"github.com/taskflow/internal/domain/auth/repository"
	"github.com/taskflow/internal/domain/auth/valueobject"
)

// PermissionDomainService 权限领域服务接口
type PermissionDomainService interface {
	// 权限检查
	CanUserPerformAction(ctx context.Context, userID string, resource valueobject.ResourceType, action valueobject.ActionType, resourceCtx map[string]interface{}) (bool, error)

	// 角色管理
	AssignRoleToUser(ctx context.Context, userID string, roleID valueobject.RoleID) error
	RevokeRoleFromUser(ctx context.Context, userID string, roleID valueobject.RoleID) error

	// 权限查询
	GetUserPermissions(ctx context.Context, userID string) ([]*aggregate.Permission, error)
	GetUserRoles(ctx context.Context, userID string) ([]*aggregate.Role, error)
}

// permissionDomainService 权限领域服务实现
type permissionDomainService struct {
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
	policyRepo     repository.PolicyRepository
	userRoleRepo   repository.UserRoleRepository
	evaluator      repository.PermissionEvaluator
	txManager      TransactionManager
}

// NewPermissionDomainService 创建权限领域服务
func NewPermissionDomainService(
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RoleRepository,
	policyRepo repository.PolicyRepository,
	userRoleRepo repository.UserRoleRepository,
	evaluator repository.PermissionEvaluator,
	txManager TransactionManager,
) PermissionDomainService {
	return &permissionDomainService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		policyRepo:     policyRepo,
		userRoleRepo:   userRoleRepo,
		evaluator:      evaluator,
		txManager:      txManager,
	}
}

// CanUserPerformAction 检查用户是否可以执行特定操作
func (s *permissionDomainService) CanUserPerformAction(
	ctx context.Context,
	userID string,
	resource valueobject.ResourceType,
	action valueobject.ActionType,
	resourceCtx map[string]interface{},
) (bool, error) {
	// 1. 获取用户角色
	userRoles, err := s.userRoleRepo.FindRolesByUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 2. 构建评估上下文
	roleIDs := make([]valueobject.RoleID, len(userRoles))
	for i, role := range userRoles {
		roleIDs[i] = role.ID
	}

	evalCtx := &repository.EvaluationContext{
		UserID:      userID,
		UserRoles:   roleIDs,
		Resource:    resource,
		Action:      action,
		ResourceCtx: resourceCtx,
		Environment: make(map[string]interface{}),
	}

	// 3. 执行权限评估
	result, err := s.evaluator.Evaluate(ctx, evalCtx)
	if err != nil {
		return false, fmt.Errorf("permission evaluation failed: %w", err)
	}

	return result.Allowed, nil
}

// AssignRoleToUser 为用户分配角色
func (s *permissionDomainService) AssignRoleToUser(ctx context.Context, userID string, roleID valueobject.RoleID) error {
	// 1. 验证角色存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// 2. 检查是否已经分配
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if hasRole {
		return domainerror.NewDomainError(domainerror.ErrRoleAlreadyAssigned, "user already has this role")
	}

	// 3. 检查系统角色分配限制
	if role.IsSystem {
		return domainerror.NewDomainError(domainerror.ErrSystemRoleImmutable, "cannot assign system role to user")
	}

	// 4. 分配角色
	if err := s.userRoleRepo.AssignRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RevokeRoleFromUser 撤销用户角色
func (s *permissionDomainService) RevokeRoleFromUser(ctx context.Context, userID string, roleID valueobject.RoleID) error {
	// 1. 验证角色存在
	role, err := s.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// 2. 检查是否已分配
	hasRole, err := s.userRoleRepo.HasRole(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if !hasRole {
		return domainerror.NewDomainError(domainerror.ErrRoleNotAssigned, "user does not have this role")
	}

	// 3. 检查系统角色撤销限制
	if role.IsSystem {
		return domainerror.NewDomainError(domainerror.ErrSystemRoleImmutable, "cannot revoke system role from user")
	}

	// 4. 撤销角色
	if err := s.userRoleRepo.RevokeRole(ctx, userID, roleID); err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	return nil
}

// GetUserPermissions 获取用户所有权限
func (s *permissionDomainService) GetUserPermissions(ctx context.Context, userID string) ([]*aggregate.Permission, error) {
	// 1. 获取用户角色
	userRoles, err := s.userRoleRepo.FindRolesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 2. 收集所有权限
	permissionMap := make(map[valueobject.PermissionID]*aggregate.Permission)

	for _, role := range userRoles {
		rolePermissions, err := s.roleRepo.FindPermissionsByRole(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}

		for _, permission := range rolePermissions {
			permissionMap[permission.ID] = permission
		}
	}

	// 3. 转换为切片
	permissions := make([]*aggregate.Permission, 0, len(permissionMap))
	for _, permission := range permissionMap {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetUserRoles 获取用户角色
func (s *permissionDomainService) GetUserRoles(ctx context.Context, userID string) ([]*aggregate.Role, error) {
	return s.userRoleRepo.FindRolesByUser(ctx, userID)
}
