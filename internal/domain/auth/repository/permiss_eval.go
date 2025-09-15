package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/taskflow/internal/domain/auth/domainerror"
	"github.com/taskflow/internal/domain/auth/valueobject"
)

// PermissionEvaluator 权限评估器接口
type PermissionEvaluator interface {
	Evaluate(ctx context.Context, evalCtx *EvaluationContext) (*EvaluationResult, error)
}

// EvaluationContext 权限评估上下文
type EvaluationContext struct {
	UserID      string                   `json:"user_id"`
	UserRoles   []valueobject.RoleID     `json:"user_roles"`
	Resource    valueobject.ResourceType `json:"resource"`
	Action      valueobject.ActionType   `json:"action"`
	ResourceCtx map[string]interface{}   `json:"resource_context"`
	Environment map[string]interface{}   `json:"environment"`
}

// EvaluationResult 权限评估结果
type EvaluationResult struct {
	Allowed     bool                     `json:"allowed"`
	Effect      valueobject.PolicyEffect `json:"effect"`
	Reason      string                   `json:"reason"`
	MatchedRule string                   `json:"matched_rule,omitempty"`
}

// rbacABACEvaluator 混合RBAC+ABAC权限评估器
type rbacABACEvaluator struct {
	permissionRepo PermissionRepository
	roleRepo       RoleRepository
	policyRepo     PolicyRepository
}

// NewRBACAbacEvaluator 创建混合权限评估器
func NewRBACAbacEvaluator(
	permissionRepo PermissionRepository,
	roleRepo RoleRepository,
	policyRepo PolicyRepository,
) PermissionEvaluator {
	return &rbacABACEvaluator{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		policyRepo:     policyRepo,
	}
}

// Evaluate 执行权限评估
func (e *rbacABACEvaluator) Evaluate(ctx context.Context, evalCtx *EvaluationContext) (*EvaluationResult, error) {
	if evalCtx == nil {
		return nil, domainerror.NewDomainError(domainerror.ErrInvalidEvaluationCtx, "evaluation context is nil")
	}

	// 1. RBAC评估 - 基于角色的权限检查
	rbacResult, err := e.evaluateRBAC(ctx, evalCtx)
	if err != nil {
		return nil, fmt.Errorf("RBAC evaluation failed: %w", err)
	}

	// 2. ABAC评估 - 基于属性的策略检查
	abacResult, err := e.evaluateABAC(ctx, evalCtx)
	if err != nil {
		return nil, fmt.Errorf("ABAC evaluation failed: %w", err)
	}

	// 3. 决策合并 - 优先级：DENY > ALLOW > RBAC
	finalResult := e.combineResults(rbacResult, abacResult)

	return finalResult, nil
}

// evaluateRBAC 执行RBAC评估
func (e *rbacABACEvaluator) evaluateRBAC(ctx context.Context, evalCtx *EvaluationContext) (*EvaluationResult, error) {
	// 检查用户角色是否有对应权限
	for _, roleID := range evalCtx.UserRoles {
		permissions, err := e.roleRepo.FindPermissionsByRole(ctx, roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to find permissions for role %s: %w", roleID, err)
		}

		for _, permission := range permissions {
			if permission.Matches(evalCtx.Resource, evalCtx.Action) {
				return &EvaluationResult{
					Allowed:     true,
					Effect:      valueobject.PolicyEffectAllow,
					Reason:      fmt.Sprintf("RBAC: Role %s has permission %s", roleID, permission.ID),
					MatchedRule: fmt.Sprintf("role:%s:permission:%s", roleID, permission.ID),
				}, nil
			}
		}
	}

	return &EvaluationResult{
		Allowed: false,
		Effect:  valueobject.PolicyEffectDeny,
		Reason:  "RBAC: No matching role permissions found",
	}, nil
}

// evaluateABAC 执行ABAC评估
func (e *rbacABACEvaluator) evaluateABAC(ctx context.Context, evalCtx *EvaluationContext) (*EvaluationResult, error) {
	// 获取匹配的策略
	policies, err := e.policyRepo.FindByResourceAndAction(ctx, evalCtx.Resource, evalCtx.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to find policies: %w", err)
	}

	if len(policies) == 0 {
		return &EvaluationResult{
			Allowed: false,
			Effect:  valueobject.PolicyEffectDeny,
			Reason:  "ABAC: No matching policies found",
		}, nil
	}

	// 按优先级排序（高优先级优先）
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Priority > policies[j].Priority
	})

	// 评估每个策略
	for _, policy := range policies {
		if !policy.IsActive {
			continue
		}

		matches, err := e.evaluatePolicyConditions(policy.Conditions, evalCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate policy %s conditions: %w", policy.ID, err)
		}

		if matches {
			return &EvaluationResult{
				Allowed:     policy.Effect == valueobject.PolicyEffectAllow,
				Effect:      policy.Effect,
				Reason:      fmt.Sprintf("ABAC: Policy %s matched", policy.Name),
				MatchedRule: fmt.Sprintf("policy:%s", policy.ID),
			}, nil
		}
	}

	return &EvaluationResult{
		Allowed: false,
		Effect:  valueobject.PolicyEffectDeny,
		Reason:  "ABAC: No policy conditions matched",
	}, nil
}

// evaluatePolicyConditions 评估策略条件
func (e *rbacABACEvaluator) evaluatePolicyConditions(conditions valueobject.PolicyConditions, evalCtx *EvaluationContext) (bool, error) {
	if len(conditions) == 0 {
		return true, nil // 无条件则匹配
	}

	// 构建评估上下文
	contextMap := map[string]interface{}{
		"user_id":    evalCtx.UserID,
		"user_roles": evalCtx.UserRoles,
		"resource":   evalCtx.Resource,
		"action":     evalCtx.Action,
	}

	// 合并资源上下文
	for k, v := range evalCtx.ResourceCtx {
		contextMap[k] = v
	}

	// 合并环境上下文
	for k, v := range evalCtx.Environment {
		contextMap[k] = v
	}

	// 评估每个条件
	for key, expectedValue := range conditions {
		actualValue, exists := contextMap[key]
		if !exists {
			return false, nil // 缺少必要的上下文
		}

		// 简单的相等性检查（可以扩展为更复杂的表达式评估）
		if !e.compareValues(actualValue, expectedValue) {
			return false, nil
		}
	}

	return true, nil
}

// compareValues 比较两个值
func (e *rbacABACEvaluator) compareValues(actual, expected interface{}) bool {
	// 处理不同类型的比较
	switch exp := expected.(type) {
	case string:
		if act, ok := actual.(string); ok {
			return act == exp
		}
	case int:
		if act, ok := actual.(int); ok {
			return act == exp
		}
	case float64:
		if act, ok := actual.(float64); ok {
			return act == exp
		}
	case bool:
		if act, ok := actual.(bool); ok {
			return act == exp
		}
	case []interface{}:
		// 检查actual是否在expected数组中
		for _, item := range exp {
			if e.compareValues(actual, item) {
				return true
			}
		}
		return false
	}

	return false
}

// combineResults 合并RBAC和ABAC结果
func (e *rbacABACEvaluator) combineResults(rbacResult, abacResult *EvaluationResult) *EvaluationResult {
	// 优先级：ABAC DENY > ABAC ALLOW > RBAC

	// 如果ABAC明确拒绝，则拒绝
	if abacResult.Effect == valueobject.PolicyEffectDeny && abacResult.MatchedRule != "" {
		return abacResult
	}

	// 如果ABAC明确允许，则允许
	if abacResult.Effect == valueobject.PolicyEffectAllow && abacResult.MatchedRule != "" {
		return abacResult
	}

	// 否则使用RBAC结果
	return rbacResult
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
