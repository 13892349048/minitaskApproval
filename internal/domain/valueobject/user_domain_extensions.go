package valueobject

import (
	"time"
)

// 用户领域服务缺失的值对象扩展

// 超级管理员角色已在user.go中定义，此处移除重复定义

// DepartmentID 部门ID值对象
type DepartmentID string

func (id DepartmentID) String() string {
	return string(id)
}

// DepartmentInfo 部门信息值对象
type DepartmentInfo struct {
	ID        DepartmentID  `json:"id"`
	Name      string        `json:"name"`
	Code      string        `json:"code"`
	ParentID  *DepartmentID `json:"parent_id,omitempty"`
	ManagerID *UserID       `json:"manager_id,omitempty"`
	Level     int           `json:"level"`
	Path      string        `json:"path"`
	IsActive  bool          `json:"is_active"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// UserHierarchyInfo 用户层级信息
type UserHierarchyInfo struct {
	UserID        UserID       `json:"user_id"`
	ManagerID     *UserID      `json:"manager_id,omitempty"`
	DepartmentID  DepartmentID `json:"department_id"`
	Level         int          `json:"level"`
	DirectReports []UserID     `json:"direct_reports"`
	AllReports    []UserID     `json:"all_reports"`
	ManagerChain  []UserID     `json:"manager_chain"`
	IsLeaf        bool         `json:"is_leaf"`
	ReportCount   int          `json:"report_count"`
}

// UserValidationContext 用户验证上下文
type UserValidationContext struct {
	OperatorID      UserID                 `json:"operator_id"`
	TargetUserID    UserID                 `json:"target_user_id"`
	Operation       string                 `json:"operation"`
	RequestData     map[string]interface{} `json:"request_data,omitempty"`
	ValidationRules []string               `json:"validation_rules,omitempty"`
	Timestamp       time.Time              `json:"timestamp"`
}

// UserValidationResult 用户验证结果
type UserValidationResult struct {
	IsValid     bool                  `json:"is_valid"`
	Context     UserValidationContext `json:"context"`
	Violations  []ValidationViolation `json:"violations,omitempty"`
	Warnings    []ValidationWarning   `json:"warnings,omitempty"`
	ValidatedAt time.Time             `json:"validated_at"`
}

// ValidationViolation 验证违规
type ValidationViolation struct {
	RuleID   string                 `json:"rule_id"`
	RuleName string                 `json:"rule_name"`
	Field    string                 `json:"field"`
	Value    interface{}            `json:"value"`
	Message  string                 `json:"message"`
	Severity ValidationSeverity     `json:"severity"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// ValidationWarning 验证警告
type ValidationWarning struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Field      string                 `json:"field,omitempty"`
	Suggestion string                 `json:"suggestion,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// ValidationSeverity 验证严重程度
type ValidationSeverity string

const (
	ValidationSeverityError   ValidationSeverity = "error"
	ValidationSeverityWarning ValidationSeverity = "warning"
	ValidationSeverityInfo    ValidationSeverity = "info"
)

// UserOperationLog 用户操作日志
type UserOperationLog struct {
	ID            string                 `json:"id"`
	OperatorID    UserID                 `json:"operator_id"`
	TargetUserID  UserID                 `json:"target_user_id"`
	Operation     string                 `json:"operation"`
	OperationType UserOperationType      `json:"operation_type"`
	BeforeState   map[string]interface{} `json:"before_state,omitempty"`
	AfterState    map[string]interface{} `json:"after_state,omitempty"`
	Changes       []FieldChange          `json:"changes,omitempty"`
	Reason        string                 `json:"reason,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	Success       bool                   `json:"success"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	Duration      time.Duration          `json:"duration"`
	CreatedAt     time.Time              `json:"created_at"`
}

// UserOperationType 用户操作类型
type UserOperationType string

const (
	UserOperationCreate             UserOperationType = "create"
	UserOperationUpdate             UserOperationType = "update"
	UserOperationDelete             UserOperationType = "delete"
	UserOperationActivate           UserOperationType = "activate"
	UserOperationDeactivate         UserOperationType = "deactivate"
	UserOperationSuspend            UserOperationType = "suspend"
	UserOperationRoleChange         UserOperationType = "role_change"
	UserOperationManagerAssign      UserOperationType = "manager_assign"
	UserOperationDepartmentTransfer UserOperationType = "department_transfer"
	UserOperationPasswordReset      UserOperationType = "password_reset"
	UserOperationPermissionGrant    UserOperationType = "permission_grant"
	UserOperationPermissionRevoke   UserOperationType = "permission_revoke"
)

// FieldChange 字段变更记录
type FieldChange struct {
	FieldName  string      `json:"field_name"`
	OldValue   interface{} `json:"old_value"`
	NewValue   interface{} `json:"new_value"`
	ChangeType string      `json:"change_type"` // "create", "update", "delete"
}

// CircularManagershipError 循环管理关系错误
type CircularManagershipError struct {
	UserID    UserID   `json:"user_id"`
	ManagerID UserID   `json:"manager_id"`
	Chain     []UserID `json:"chain"`
	Message   string   `json:"message"`
}

func (e CircularManagershipError) Error() string {
	return e.Message
}

// DepartmentValidationError 部门验证错误
type DepartmentValidationError struct {
	DepartmentID DepartmentID `json:"department_id"`
	ErrorType    string       `json:"error_type"` // "not_found", "inactive", "invalid"
	Message      string       `json:"message"`
}

func (e DepartmentValidationError) Error() string {
	return e.Message
}

// UserBusinessRule 用户业务规则
type UserBusinessRule struct {
	RuleID     string                 `json:"rule_id"`
	RuleName   string                 `json:"rule_name"`
	RuleType   UserBusinessRuleType   `json:"rule_type"`
	Conditions map[string]interface{} `json:"conditions"`
	Actions    []string               `json:"actions"`
	Priority   int                    `json:"priority"`
	IsActive   bool                   `json:"is_active"`
	CreatedBy  UserID                 `json:"created_by"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// UserBusinessRuleType 用户业务规则类型
type UserBusinessRuleType string

const (
	UserRuleTypeValidation    UserBusinessRuleType = "validation"
	UserRuleTypeAuthorization UserBusinessRuleType = "authorization"
	UserRuleTypeWorkflow      UserBusinessRuleType = "workflow"
	UserRuleTypeNotification  UserBusinessRuleType = "notification"
	UserRuleTypeAudit         UserBusinessRuleType = "audit"
)

// UserPermissionMatrix 用户权限矩阵
type UserPermissionMatrix struct {
	UserID      UserID                   `json:"user_id"`
	Role        UserRole                 `json:"role"`
	Permissions map[string]PermissionSet `json:"permissions"`
	Constraints map[string]interface{}   `json:"constraints,omitempty"`
	ValidFrom   time.Time                `json:"valid_from"`
	ValidUntil  *time.Time               `json:"valid_until,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// PermissionSet 权限集合
type PermissionSet struct {
	Resource   string   `json:"resource"`
	Actions    []string `json:"actions"`
	Conditions []string `json:"conditions,omitempty"`
	Scope      string   `json:"scope"` // "own", "department", "company", "all"
}
