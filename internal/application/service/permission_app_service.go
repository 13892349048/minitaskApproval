package service

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/auth/aggregate"
	"github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// PermissionAppService 权限应用服务
type PermissionAppService struct {
	domainService service.PermissionDomainService
	txManager     service.TransactionManager
}

// NewPermissionAppService 创建权限应用服务
func NewPermissionAppService(
	domainService service.PermissionDomainService,
	txManager service.TransactionManager,
) *PermissionAppService {
	return &PermissionAppService{
		domainService: domainService,
		txManager:     txManager,
	}
}

// CanUserPerformAction 检查用户是否可以执行某个操作
func (s *PermissionAppService) CanUserPerformAction(
	ctx context.Context,
	userID, resource, action string,
	resourceContext map[string]interface{},
) (bool, error) {
	// 转换为域类型
	resourceType := valueobject.ResourceType(resource)
	actionType := valueobject.ActionType(action)

	// 调用领域服务
	allowed, err := s.domainService.CanUserPerformAction(ctx, userID, resourceType, actionType, resourceContext)
	if err != nil {
		logger.Error("Permission check failed",
			zap.String("user_id", userID),
			zap.String("resource", resource),
			zap.String("action", action),
			zap.Error(err))
		return false, fmt.Errorf("permission check failed: %w", err)
	}

	logger.Debug("Permission check completed",
		zap.String("user_id", userID),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.Bool("allowed", allowed))

	return allowed, nil
}

// GetUserPermissions 获取用户的所有权限
func (s *PermissionAppService) GetUserPermissions(ctx context.Context, userID string) ([]aggregate.Permission, error) {
	domainPermissions, err := s.domainService.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user permissions",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// 转换为应用层DTO
	permissions := make([]aggregate.Permission, len(domainPermissions))
	for i, perm := range domainPermissions {
		permissions[i] = aggregate.Permission{
			ID:          valueobject.PermissionID(perm.ID),
			Name:        perm.Name,
			Resource:    valueobject.ResourceType(perm.Resource),
			Action:      valueobject.ActionType(perm.Action),
			Description: perm.Description,
		}
	}

	return permissions, nil
}

// HasRole 检查用户是否具有某个角色
func (s *PermissionAppService) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
	userRoles, err := s.domainService.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user roles",
			zap.String("user_id", userID),
			zap.Error(err))
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 检查是否有匹配的角色
	for _, role := range userRoles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

// GetUserRoles 获取用户的所有角色
func (s *PermissionAppService) GetUserRoles(ctx context.Context, userID string) ([]aggregate.Role, error) {
	domainRoles, err := s.domainService.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user roles",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 转换为应用层DTO
	roles := make([]aggregate.Role, len(domainRoles))
	for i, role := range domainRoles {
		roles[i] = aggregate.Role{
			ID:          valueobject.RoleID(role.ID),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
			IsSystem:    role.IsSystem,
		}
	}

	return roles, nil
}

// AssignRoleToUser 为用户分配角色
func (s *PermissionAppService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		err := s.domainService.AssignRoleToUser(txCtx, userID, valueobject.RoleID(roleID))
		if err != nil {
			logger.Error("Failed to assign role to user",
				zap.String("user_id", userID),
				zap.String("role_id", roleID),
				zap.Error(err))
			return fmt.Errorf("failed to assign role: %w", err)
		}

		logger.Info("Role assigned to user",
			zap.String("user_id", userID),
			zap.String("role_id", roleID))

		return nil
	})
}

// RevokeRoleFromUser 撤销用户角色
func (s *PermissionAppService) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		err := s.domainService.RevokeRoleFromUser(txCtx, userID, valueobject.RoleID(roleID))
		if err != nil {
			logger.Error("Failed to revoke role from user",
				zap.String("user_id", userID),
				zap.String("role_id", roleID),
				zap.Error(err))
			return fmt.Errorf("failed to revoke role: %w", err)
		}

		logger.Info("Role revoked from user",
			zap.String("user_id", userID),
			zap.String("role_id", roleID))

		return nil
	})
}
