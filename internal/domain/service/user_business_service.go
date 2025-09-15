package service

import (
	"context"
	"github.com/taskflow/internal/domain/valueobject"
)

// UserBusinessService 用户业务服务接口
// 处理复杂的用户业务逻辑和跨聚合操作
type UserBusinessService interface {
	// 用户转移相关
	ValidateUserTransfer(ctx context.Context, request *valueobject.UserTransferRequest) (*valueobject.UserOperationResult, error)
	ExecuteUserTransfer(ctx context.Context, request *valueobject.UserTransferRequest) error
	
	// 任务转移相关
	ValidateTaskTransfer(ctx context.Context, request *valueobject.TaskTransferRequest) (*valueobject.UserOperationResult, error)
	ExecuteTaskTransfer(ctx context.Context, request *valueobject.TaskTransferRequest) error
	
	// 用户停用相关
	ValidateUserDeactivation(ctx context.Context, request *valueobject.UserDeactivationRequest) (*valueobject.UserOperationResult, error)
	ExecuteUserDeactivation(ctx context.Context, request *valueobject.UserDeactivationRequest) error
	
	// 角色变更相关
	ValidateRoleChange(ctx context.Context, request *valueobject.RoleChangeRequest) (*valueobject.UserOperationResult, error)
	ExecuteRoleChange(ctx context.Context, request *valueobject.RoleChangeRequest) error
	
	// 管理者分配相关
	ValidateManagerAssignment(ctx context.Context, request *valueobject.ManagerAssignmentRequest) (*valueobject.UserOperationResult, error)
	ExecuteManagerAssignment(ctx context.Context, request *valueobject.ManagerAssignmentRequest) error
	
	// 工作负载管理
	CalculateUserWorkload(ctx context.Context, userID valueobject.UserID) (*valueobject.UserWorkload, error)
	CheckWorkloadBalance(ctx context.Context, departmentID string) ([]valueobject.UserWorkload, error)
	
	// 绩效评估
	CalculateUserPerformance(ctx context.Context, userID valueobject.UserID, period string) (*valueobject.UserPerformanceMetrics, error)
	
	// 层级管理
	GetUserHierarchy(ctx context.Context, userID valueobject.UserID) (*valueobject.UserHierarchyPosition, error)
	GetDepartmentHierarchy(ctx context.Context, departmentID string) (*valueobject.DepartmentHierarchy, error)
	
	// 业务规则验证
	ValidateBusinessRules(ctx context.Context, context *valueobject.UserOperationContext) ([]valueobject.BusinessRuleViolation, error)
}

// UserValidationService 用户验证服务接口
type UserValidationService interface {
	// 基础验证
	ValidateUserCreation(ctx context.Context, email, username string) error
	ValidateEmailFormat(email string) error
	ValidateUsernameFormat(username string) error
	ValidatePasswordStrength(password string) error
	
	// 业务规则验证
	ValidateUserRole(ctx context.Context, userID valueobject.UserID, role valueobject.UserRole) error
	ValidateDepartmentAssignment(ctx context.Context, userID valueobject.UserID, departmentID string) error
	ValidateManagerHierarchy(ctx context.Context, userID, managerID valueobject.UserID) error
	
	// 权限验证
	ValidateOperationPermission(ctx context.Context, operatorID valueobject.UserID, operation string, targetUserID valueobject.UserID) error
	
	// 数据完整性验证
	ValidateUserDataIntegrity(ctx context.Context, userID valueobject.UserID) ([]valueobject.BusinessRuleViolation, error)
}

// UserAnalyticsService 用户分析服务接口
type UserAnalyticsService interface {
	// 用户统计
	GetUserStatistics(ctx context.Context, userID valueobject.UserID) (*valueobject.UserStatistics, error)
	GetDepartmentStatistics(ctx context.Context, departmentID string) (*DepartmentStatistics, error)
	
	// 趋势分析
	GetUserPerformanceTrend(ctx context.Context, userID valueobject.UserID, months int) ([]valueobject.UserPerformanceMetrics, error)
	GetWorkloadTrend(ctx context.Context, userID valueobject.UserID, days int) ([]valueobject.UserWorkload, error)
	
	// 比较分析
	CompareUserPerformance(ctx context.Context, userIDs []valueobject.UserID, period string) ([]UserPerformanceComparison, error)
	
	// 预测分析
	PredictUserWorkload(ctx context.Context, userID valueobject.UserID, days int) (*WorkloadPrediction, error)
	IdentifyRiskUsers(ctx context.Context, departmentID string) ([]UserRiskAssessment, error)
}

// 辅助数据结构

// DepartmentStatistics 部门统计
type DepartmentStatistics struct {
	DepartmentID        string  `json:"department_id"`
	TotalUsers          int     `json:"total_users"`
	ActiveUsers         int     `json:"active_users"`
	AveragePerformance  float64 `json:"average_performance"`
	AverageWorkload     float64 `json:"average_workload"`
	TotalTasks          int     `json:"total_tasks"`
	CompletedTasks      int     `json:"completed_tasks"`
	OverdueTasksRate    float64 `json:"overdue_tasks_rate"`
	TurnoverRate        float64 `json:"turnover_rate"`
}

// UserPerformanceComparison 用户绩效比较
type UserPerformanceComparison struct {
	UserID           valueobject.UserID                `json:"user_id"`
	Metrics          valueobject.UserPerformanceMetrics `json:"metrics"`
	RelativeRanking  int                               `json:"relative_ranking"`
	PerformanceGap   float64                           `json:"performance_gap"`
	ImprovementAreas []string                          `json:"improvement_areas"`
}

// WorkloadPrediction 工作负载预测
type WorkloadPrediction struct {
	UserID              valueobject.UserID `json:"user_id"`
	PredictedWorkload   float64           `json:"predicted_workload"`
	ConfidenceLevel     float64           `json:"confidence_level"`
	RiskLevel           string            `json:"risk_level"`
	RecommendedActions  []string          `json:"recommended_actions"`
}

// UserRiskAssessment 用户风险评估
type UserRiskAssessment struct {
	UserID          valueobject.UserID `json:"user_id"`
	RiskScore       float64           `json:"risk_score"`
	RiskCategory    string            `json:"risk_category"`
	RiskFactors     []string          `json:"risk_factors"`
	Recommendations []string          `json:"recommendations"`
	Priority        string            `json:"priority"`
}
