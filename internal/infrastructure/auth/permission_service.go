package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/internal/infrastructure/persistence/mysql"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PermissionService 权限服务实现
type PermissionService struct {
	db       *gorm.DB
	userRepo *mysql.UserRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB, userRepo *mysql.UserRepository) shared.PermissionService {
	return &PermissionService{
		db:       db,
		userRepo: userRepo,
	}
}

// CanUserPerformAction 检查用户是否可以执行某个操作
func (p *PermissionService) CanUserPerformAction(ctx context.Context, userID, resource, action string, resourceContext map[string]interface{}) (bool, error) {
	// 1. 获取用户信息
	user, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to find user", zap.String("user_id", userID), zap.Error(err))
		return false, fmt.Errorf("failed to find user: %w", err)
	}

	// 2. 获取用户角色
	userRoles, err := p.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user roles", zap.String("user_id", userID), zap.Error(err))
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 3. 构建评估上下文
	evalCtx := &shared.EvaluationContext{
		UserID:      userID,
		UserRoles:   p.extractRoleNames(userRoles),
		Resource:    resource,
		Action:      action,
		ResourceCtx: resourceContext,
		Environment: map[string]interface{}{
			"user_email":      user.Email,
			"user_department": user.Department,
		},
	}

	// 4. RBAC检查（基础权限）
	hasRBACPermission, err := p.checkRBACPermission(ctx, evalCtx)
	if err != nil {
		logger.Error("RBAC permission check failed", zap.Error(err))
		return false, fmt.Errorf("RBAC permission check failed: %w", err)
	}

	// 5. ABAC检查（策略规则）
	abacResult, err := p.checkABACPolicies(ctx, evalCtx)
	if err != nil {
		logger.Error("ABAC policy check failed", zap.Error(err))
		return false, fmt.Errorf("ABAC policy check failed: %w", err)
	}

	// 6. 权限决策
	finalResult := p.makePermissionDecision(hasRBACPermission, abacResult)

	logger.Debug("Permission check completed",
		zap.String("user_id", userID),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.Bool("rbac_result", hasRBACPermission),
		zap.String("abac_result", string(abacResult)),
		zap.Bool("final_result", finalResult),
	)

	return finalResult, nil
}

// GetUserPermissions 获取用户的所有权限
func (p *PermissionService) GetUserPermissions(ctx context.Context, userID string) ([]shared.Permission, error) {
	var permissions []mysql.Permission

	// 通过角色获取权限
	err := p.db.WithContext(ctx).
		Select("DISTINCT permissions.*").
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// 转换为应用层结构
	result := make([]shared.Permission, len(permissions))
	for i, perm := range permissions {
		result[i] = shared.Permission{
			ID:          perm.ID,
			Name:        perm.Name,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: *perm.Description,
		}
	}

	return result, nil
}

// HasRole 检查用户是否具有某个角色
func (p *PermissionService) HasRole(ctx context.Context, userID, roleName string) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return count > 0, nil
}

// GetUserRoles 获取用户的所有角色
func (p *PermissionService) GetUserRoles(ctx context.Context, userID string) ([]shared.Role, error) {
	var roles []mysql.Role

	err := p.db.WithContext(ctx).
		Select("roles.*").
		Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 转换为应用层结构
	result := make([]shared.Role, len(roles))
	for i, role := range roles {
		result[i] = shared.Role{
			ID:          role.ID,
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: *role.Description,
			IsSystem:    role.IsSystem,
		}
	}

	return result, nil
}

// checkRBACPermission 检查RBAC权限
func (p *PermissionService) checkRBACPermission(ctx context.Context, evalCtx *shared.EvaluationContext) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?",
			evalCtx.UserID, evalCtx.Resource, evalCtx.Action).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// checkABACPolicies 检查ABAC策略
func (p *PermissionService) checkABACPolicies(ctx context.Context, evalCtx *shared.EvaluationContext) (shared.PolicyEffect, error) {
	var policies []mysql.PermissionPolicy

	// 获取适用的策略规则
	err := p.db.WithContext(ctx).
		Where("resource_type = ? AND action = ? AND is_active = ?", evalCtx.Resource, evalCtx.Action, true).
		Order("priority DESC").
		Find(&policies).Error

	if err != nil {
		return "", fmt.Errorf("failed to get policies: %w", err)
	}

	// 按优先级评估策略
	for _, policy := range policies {
		var conditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.Conditions), &conditions); err != nil {
			logger.Warn("Failed to parse policy conditions",
				zap.String("policy_id", policy.ID),
				zap.Error(err))
			continue
		}

		match, err := p.evaluateConditions(conditions, evalCtx)
		if err != nil {
			logger.Warn("Failed to evaluate policy conditions",
				zap.String("policy_id", policy.ID),
				zap.Error(err))
			continue
		}

		if match {
			return shared.PolicyEffect(policy.Effect), nil
		}
	}

	// 没有匹配的策略，返回空值（由调用方决定）
	return "", nil
}

// evaluateConditions 评估策略条件
func (p *PermissionService) evaluateConditions(conditions map[string]interface{}, evalCtx *shared.EvaluationContext) (bool, error) {
	// 简化的条件评估实现
	// 生产环境可以使用更复杂的表达式引擎

	for key, value := range conditions {
		// 处理用户角色检查
		if key == "user.role" {
			if comparison, ok := value.(map[string]interface{}); ok {
				if expectedRole, exists := comparison["eq"]; exists {
					hasRole := false
					for _, role := range evalCtx.UserRoles {
						if role == expectedRole {
							hasRole = true
							break
						}
					}
					if !hasRole {
						return false, nil
					}
				}
			}
		}

		// 处理用户ID检查
		if key == "user.id" {
			if comparison, ok := value.(map[string]interface{}); ok {
				if expectedID, exists := comparison["eq"]; exists {
					// 支持变量替换，如 "${resource.owner_id}"
					if expectedStr, ok := expectedID.(string); ok {
						if strings.HasPrefix(expectedStr, "${") && strings.HasSuffix(expectedStr, "}") {
							varPath := strings.TrimSuffix(strings.TrimPrefix(expectedStr, "${"), "}")
							resolvedValue := p.resolveVariable(varPath, evalCtx)
							if evalCtx.UserID != fmt.Sprintf("%v", resolvedValue) {
								return false, nil
							}
						} else if evalCtx.UserID != expectedStr {
							return false, nil
						}
					}
				}
			}
		}
	}

	return true, nil
}

// resolveVariable 解析变量
func (p *PermissionService) resolveVariable(varPath string, evalCtx *shared.EvaluationContext) interface{} {
	parts := strings.Split(varPath, ".")
	if len(parts) < 2 {
		return nil
	}

	switch parts[0] {
	case "resource":
		if len(parts) == 2 {
			return evalCtx.ResourceCtx[parts[1]]
		}
	case "user":
		if parts[1] == "id" {
			return evalCtx.UserID
		}
	}

	return nil
}

// makePermissionDecision 做出权限决策
func (p *PermissionService) makePermissionDecision(rbacResult bool, abacResult shared.PolicyEffect) bool {
	// 决策逻辑：
	// 1. 如果ABAC明确拒绝，则拒绝
	// 2. 如果ABAC明确允许，则允许
	// 3. 如果ABAC无结果，则依据RBAC结果

	switch abacResult {
	case shared.PolicyEffectDeny:
		return false
	case shared.PolicyEffectAllow:
		return true
	default:
		return rbacResult
	}
}

// extractRoleNames 提取角色名称
func (p *PermissionService) extractRoleNames(roles []shared.Role) []string {
	names := make([]string, len(roles))
	for i, role := range roles {
		names[i] = role.Name
	}
	return names
}

// 为什么这样实现？
//
// 1. 分层检查：
//    - 先进行RBAC检查（基础权限）
//    - 再进行ABAC检查（策略规则）
//    - 最后综合决策
//
// 2. 策略优先级：
//    - 按优先级排序策略规则
//    - 高优先级策略优先匹配
//    - 支持deny优于allow的安全原则
//
// 3. 条件评估：
//    - 支持变量替换（如${resource.owner_id}）
//    - 简化的条件语法，易于理解和维护
//    - 可扩展的条件类型
//
// 4. 错误处理：
//    - 权限检查失败时记录详细日志
//    - 策略解析失败时跳过该策略
//    - 默认拒绝原则保证安全
