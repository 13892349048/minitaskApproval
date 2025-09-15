package valueobject

import (
	"time"
)

// 用户业务规则相关值对象

// UserTransferRequest 用户转移请求
type UserTransferRequest struct {
	UserID           UserID    `json:"user_id"`
	FromDepartmentID string    `json:"from_department_id"`
	ToDepartmentID   string    `json:"to_department_id"`
	NewManagerID     UserID    `json:"new_manager_id"`
	TransferReason   string    `json:"transfer_reason"`
	EffectiveDate    time.Time `json:"effective_date"`
	RequestedBy      UserID    `json:"requested_by"`
	ApprovedBy       *UserID   `json:"approved_by,omitempty"`
	RequestedAt      time.Time `json:"requested_at"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
}

// TaskTransferRequest 任务转移请求
type TaskTransferRequest struct {
	FromUserID   UserID    `json:"from_user_id"`
	ToUserID     UserID    `json:"to_user_id"`
	TaskIDs      []TaskID  `json:"task_ids"`
	Reason       string    `json:"reason"`
	TransferType string    `json:"transfer_type"` // "temporary", "permanent", "delegation"
	RequestedBy  UserID    `json:"requested_by"`
	RequestedAt  time.Time `json:"requested_at"`
	DueDate      *time.Time `json:"due_date,omitempty"`
}

// UserDeactivationRequest 用户停用请求
type UserDeactivationRequest struct {
	UserID              UserID                `json:"user_id"`
	DeactivatedBy       UserID                `json:"deactivated_by"`
	Reason              string                `json:"reason"`
	TaskTransferPlan    *TaskTransferPlan     `json:"task_transfer_plan,omitempty"`
	ProjectTransferPlan *ProjectTransferPlan  `json:"project_transfer_plan,omitempty"`
	DataRetentionPlan   *DataRetentionPlan    `json:"data_retention_plan,omitempty"`
	EffectiveDate       time.Time             `json:"effective_date"`
	NotifyContacts      []UserID              `json:"notify_contacts,omitempty"`
}

// TaskTransferPlan 任务转移计划
type TaskTransferPlan struct {
	TransferStrategy string              `json:"transfer_strategy"` // "to_manager", "to_specific_user", "redistribute"
	TargetUserID     *UserID             `json:"target_user_id,omitempty"`
	TaskDistribution []TaskDistribution  `json:"task_distribution,omitempty"`
	AutoTransfer     bool                `json:"auto_transfer"`
	RequireApproval  bool                `json:"require_approval"`
}

// TaskDistribution 任务分配
type TaskDistribution struct {
	TaskID     TaskID `json:"task_id"`
	AssignedTo UserID `json:"assigned_to"`
	Priority   int    `json:"priority"`
	Notes      string `json:"notes,omitempty"`
}

// ProjectTransferPlan 项目转移计划
type ProjectTransferPlan struct {
	TransferStrategy string                `json:"transfer_strategy"`
	TargetUserID     *UserID               `json:"target_user_id,omitempty"`
	ProjectRoles     []ProjectRoleTransfer `json:"project_roles,omitempty"`
}

// ProjectRoleTransfer 项目角色转移
type ProjectRoleTransfer struct {
	ProjectID  ProjectID   `json:"project_id"`
	FromRole   ProjectRole `json:"from_role"`
	ToRole     ProjectRole `json:"to_role"`
	AssignedTo UserID      `json:"assigned_to"`
}

// DataRetentionPlan 数据保留计划
type DataRetentionPlan struct {
	RetentionPeriod   int    `json:"retention_period"` // 天数
	ArchiveData       bool   `json:"archive_data"`
	DeletePersonalData bool   `json:"delete_personal_data"`
	BackupLocation    string `json:"backup_location,omitempty"`
	ResponsibleUser   UserID `json:"responsible_user"`
}

// RoleChangeRequest 角色变更请求
type RoleChangeRequest struct {
	UserID      UserID    `json:"user_id"`
	FromRole    UserRole  `json:"from_role"`
	ToRole      UserRole  `json:"to_role"`
	Reason      string    `json:"reason"`
	ChangedBy   UserID    `json:"changed_by"`
	RequestedAt time.Time `json:"requested_at"`
	EffectiveDate time.Time `json:"effective_date"`
	RequiresApproval bool `json:"requires_approval"`
}

// ManagerAssignmentRequest 管理者分配请求
type ManagerAssignmentRequest struct {
	UserID        UserID    `json:"user_id"`
	FromManagerID *UserID   `json:"from_manager_id,omitempty"`
	ToManagerID   UserID    `json:"to_manager_id"`
	Reason        string    `json:"reason"`
	AssignedBy    UserID    `json:"assigned_by"`
	RequestedAt   time.Time `json:"requested_at"`
	EffectiveDate time.Time `json:"effective_date"`
}

// UserValidationRule 用户验证规则
type UserValidationRule struct {
	RuleID      string                 `json:"rule_id"`
	RuleName    string                 `json:"rule_name"`
	RuleType    string                 `json:"rule_type"` // "email", "username", "role", "department"
	Conditions  map[string]interface{} `json:"conditions"`
	ErrorMessage string                `json:"error_message"`
	IsActive    bool                   `json:"is_active"`
	Priority    int                    `json:"priority"`
}

// BusinessRuleViolation 业务规则违反
type BusinessRuleViolation struct {
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	ViolationType string  `json:"violation_type"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // "error", "warning", "info"
	UserID      UserID    `json:"user_id"`
	Context     map[string]interface{} `json:"context,omitempty"`
	DetectedAt  time.Time `json:"detected_at"`
}

// UserOperationContext 用户操作上下文
type UserOperationContext struct {
	OperatorID    UserID                 `json:"operator_id"`
	OperationType string                 `json:"operation_type"`
	TargetUserID  UserID                 `json:"target_user_id"`
	RequestData   map[string]interface{} `json:"request_data,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

// UserOperationResult 用户操作结果
type UserOperationResult struct {
	Success       bool                    `json:"success"`
	OperationID   string                  `json:"operation_id"`
	Context       UserOperationContext    `json:"context"`
	Result        map[string]interface{}  `json:"result,omitempty"`
	Violations    []BusinessRuleViolation `json:"violations,omitempty"`
	Warnings      []string                `json:"warnings,omitempty"`
	ProcessedAt   time.Time               `json:"processed_at"`
	Duration      time.Duration           `json:"duration"`
}

// DepartmentHierarchy 部门层级结构
type DepartmentHierarchy struct {
	DepartmentID   string                `json:"department_id"`
	ParentID       *string               `json:"parent_id,omitempty"`
	Level          int                   `json:"level"`
	Path           []string              `json:"path"`
	Children       []DepartmentHierarchy `json:"children,omitempty"`
	ManagerID      *UserID               `json:"manager_id,omitempty"`
	MemberCount    int                   `json:"member_count"`
	ActiveMembers  int                   `json:"active_members"`
}

// UserHierarchyPosition 用户层级位置
type UserHierarchyPosition struct {
	UserID         UserID   `json:"user_id"`
	Level          int      `json:"level"`
	DirectReports  []UserID `json:"direct_reports"`
	AllSubordinates []UserID `json:"all_subordinates"`
	ManagerChain   []UserID `json:"manager_chain"`
	DepartmentPath []string `json:"department_path"`
}

// UserCapability 用户能力
type UserCapability struct {
	CapabilityID   string    `json:"capability_id"`
	CapabilityName string    `json:"capability_name"`
	Level          int       `json:"level"` // 1-5 能力等级
	CertifiedAt    time.Time `json:"certified_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CertifiedBy    UserID    `json:"certified_by"`
}

// UserWorkload 用户工作负载
type UserWorkload struct {
	UserID              UserID    `json:"user_id"`
	ActiveTasks         int       `json:"active_tasks"`
	ActiveProjects      int       `json:"active_projects"`
	EstimatedHours      int       `json:"estimated_hours"`
	AvailableHours      int       `json:"available_hours"`
	UtilizationRate     float64   `json:"utilization_rate"`
	OverloadThreshold   float64   `json:"overload_threshold"`
	IsOverloaded        bool      `json:"is_overloaded"`
	LastCalculatedAt    time.Time `json:"last_calculated_at"`
}

// UserPerformanceMetrics 用户绩效指标
type UserPerformanceMetrics struct {
	UserID                UserID    `json:"user_id"`
	Period                string    `json:"period"` // "monthly", "quarterly", "yearly"
	TaskCompletionRate    float64   `json:"task_completion_rate"`
	OnTimeDeliveryRate    float64   `json:"on_time_delivery_rate"`
	QualityScore          float64   `json:"quality_score"`
	CollaborationScore    float64   `json:"collaboration_score"`
	LeadershipScore       float64   `json:"leadership_score"`
	OverallScore          float64   `json:"overall_score"`
	Ranking               int       `json:"ranking"`
	TotalParticipants     int       `json:"total_participants"`
	CalculatedAt          time.Time `json:"calculated_at"`
}
