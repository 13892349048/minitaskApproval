package valueobject

import (
	"time"
)

// 审批流程相关值对象

// ApprovalID 审批ID值对象
type ApprovalID string

func (id ApprovalID) String() string {
	return string(id)
}

// WorkflowID 工作流ID值对象
type WorkflowID string

func (id WorkflowID) String() string {
	return string(id)
}

// ApprovalStatus 审批状态
type ApprovalStatus string

const (
	ApprovalStatusPending   ApprovalStatus = "pending"   // 待审批
	ApprovalStatusApproved  ApprovalStatus = "approved"  // 已批准
	ApprovalStatusRejected  ApprovalStatus = "rejected"  // 已拒绝
	ApprovalStatusWithdrawn ApprovalStatus = "withdrawn" // 已撤回
	ApprovalStatusExpired   ApprovalStatus = "expired"   // 已过期
)

// ApprovalType 审批类型
type ApprovalType string

const (
	ApprovalTypeTask        ApprovalType = "task"         // 任务审批
	ApprovalTypeProject     ApprovalType = "project"      // 项目审批
	ApprovalTypeLeave       ApprovalType = "leave"        // 请假审批
	ApprovalTypeExpense     ApprovalType = "expense"      // 费用审批
	ApprovalTypePurchase    ApprovalType = "purchase"     // 采购审批
	ApprovalTypeRecruitment ApprovalType = "recruitment"  // 招聘审批
)

// ApprovalAction 审批动作
type ApprovalAction string

const (
	ApprovalActionSubmit   ApprovalAction = "submit"   // 提交
	ApprovalActionApprove  ApprovalAction = "approve"  // 批准
	ApprovalActionReject   ApprovalAction = "reject"   // 拒绝
	ApprovalActionWithdraw ApprovalAction = "withdraw" // 撤回
	ApprovalActionDelegate ApprovalAction = "delegate" // 委托
	ApprovalActionReturn   ApprovalAction = "return"   // 退回
)

// ApprovalLevel 审批级别
type ApprovalLevel int

const (
	ApprovalLevelDepartment ApprovalLevel = 1 // 部门级别
	ApprovalLevelDivision   ApprovalLevel = 2 // 事业部级别
	ApprovalLevelCompany    ApprovalLevel = 3 // 公司级别
	ApprovalLevelBoard      ApprovalLevel = 4 // 董事会级别
)

// ApprovalPriority 审批优先级
type ApprovalPriority string

const (
	ApprovalPriorityLow      ApprovalPriority = "low"      // 低优先级
	ApprovalPriorityNormal   ApprovalPriority = "normal"   // 普通优先级
	ApprovalPriorityHigh     ApprovalPriority = "high"     // 高优先级
	ApprovalPriorityUrgent   ApprovalPriority = "urgent"   // 紧急优先级
	ApprovalPriorityCritical ApprovalPriority = "critical" // 关键优先级
)

// ApprovalStep 审批步骤值对象
type ApprovalStep struct {
	StepID      string           `json:"step_id"`
	StepName    string           `json:"step_name"`
	Level       ApprovalLevel    `json:"level"`
	ApproverID  UserID           `json:"approver_id"`
	Status      ApprovalStatus   `json:"status"`
	Action      *ApprovalAction  `json:"action,omitempty"`
	Comment     string           `json:"comment,omitempty"`
	ProcessedAt *time.Time       `json:"processed_at,omitempty"`
	DueDate     *time.Time       `json:"due_date,omitempty"`
	IsRequired  bool             `json:"is_required"`
	CanDelegate bool             `json:"can_delegate"`
	DelegatedTo *UserID          `json:"delegated_to,omitempty"`
}

// ApprovalHistory 审批历史记录
type ApprovalHistory struct {
	ID          string          `json:"id"`
	ApprovalID  ApprovalID      `json:"approval_id"`
	StepID      string          `json:"step_id"`
	Action      ApprovalAction  `json:"action"`
	ActorID     UserID          `json:"actor_id"`
	Comment     string          `json:"comment,omitempty"`
	Attachments []string        `json:"attachments,omitempty"`
	ProcessedAt time.Time       `json:"processed_at"`
	IPAddress   string          `json:"ip_address,omitempty"`
	UserAgent   string          `json:"user_agent,omitempty"`
}

// ApprovalRule 审批规则值对象
type ApprovalRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        ApprovalType           `json:"type"`
	Conditions  map[string]interface{} `json:"conditions"`
	Steps       []ApprovalStepRule     `json:"steps"`
	IsActive    bool                   `json:"is_active"`
	CreatedBy   UserID                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ApprovalStepRule 审批步骤规则
type ApprovalStepRule struct {
	StepID       string        `json:"step_id"`
	StepName     string        `json:"step_name"`
	Level        ApprovalLevel `json:"level"`
	ApproverRole string        `json:"approver_role,omitempty"`
	ApproverID   *UserID       `json:"approver_id,omitempty"`
	IsRequired   bool          `json:"is_required"`
	CanDelegate  bool          `json:"can_delegate"`
	TimeoutHours int           `json:"timeout_hours,omitempty"`
	AutoApprove  bool          `json:"auto_approve"`
}

// ApprovalRequest 审批请求值对象
type ApprovalRequest struct {
	ID           ApprovalID       `json:"id"`
	Type         ApprovalType     `json:"type"`
	Title        string           `json:"title"`
	Description  string           `json:"description,omitempty"`
	Priority     ApprovalPriority `json:"priority"`
	RequesterID  UserID           `json:"requester_id"`
	EntityID     string           `json:"entity_id"` // 关联的实体ID（任务、项目等）
	EntityType   string           `json:"entity_type"`
	Status       ApprovalStatus   `json:"status"`
	CurrentStep  *string          `json:"current_step,omitempty"`
	Steps        []ApprovalStep   `json:"steps"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Attachments  []string         `json:"attachments,omitempty"`
	SubmittedAt  time.Time        `json:"submitted_at"`
	CompletedAt  *time.Time       `json:"completed_at,omitempty"`
	DueDate      *time.Time       `json:"due_date,omitempty"`
}

// ApprovalData 审批数据传输对象
type ApprovalData struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	Description  *string                `json:"description,omitempty"`
	Priority     string                 `json:"priority"`
	RequesterID  string                 `json:"requester_id"`
	EntityID     string                 `json:"entity_id"`
	EntityType   string                 `json:"entity_type"`
	Status       string                 `json:"status"`
	CurrentStep  *string                `json:"current_step,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Attachments  []string               `json:"attachments,omitempty"`
	SubmittedAt  time.Time              `json:"submitted_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	DueDate      *time.Time             `json:"due_date,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}

// ApprovalSearchRequest 审批搜索请求
type ApprovalSearchRequest struct {
	BaseSearchCriteria
	Type        *ApprovalType     `json:"type,omitempty"`
	Status      *ApprovalStatus   `json:"status,omitempty"`
	Priority    *ApprovalPriority `json:"priority,omitempty"`
	RequesterID *UserID           `json:"requester_id,omitempty"`
	ApproverID  *UserID           `json:"approver_id,omitempty"`
	EntityType  *string           `json:"entity_type,omitempty"`
	EntityID    *string           `json:"entity_id,omitempty"`
}

// ApprovalSummary 审批摘要信息
type ApprovalSummary struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"`
	Title       string           `json:"title"`
	Priority    string           `json:"priority"`
	Status      string           `json:"status"`
	Requester   UserSummary      `json:"requester"`
	CurrentStep *string          `json:"current_step,omitempty"`
	SubmittedAt time.Time        `json:"submitted_at"`
	DueDate     *time.Time       `json:"due_date,omitempty"`
}

// ApprovalDetailResponse 审批详细信息响应
type ApprovalDetailResponse struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	Description  *string                `json:"description,omitempty"`
	Priority     string                 `json:"priority"`
	Status       string                 `json:"status"`
	Requester    UserSummary            `json:"requester"`
	EntityID     string                 `json:"entity_id"`
	EntityType   string                 `json:"entity_type"`
	CurrentStep  *string                `json:"current_step,omitempty"`
	Steps        []ApprovalStep         `json:"steps"`
	History      []ApprovalHistory      `json:"history"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Attachments  []string               `json:"attachments,omitempty"`
	SubmittedAt  time.Time              `json:"submitted_at"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	DueDate      *time.Time             `json:"due_date,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ApprovalActionRequest 审批动作请求
type ApprovalActionRequest struct {
	ApprovalID  string         `json:"approval_id" validate:"required"`
	StepID      string         `json:"step_id" validate:"required"`
	Action      ApprovalAction `json:"action" validate:"required"`
	Comment     string         `json:"comment,omitempty"`
	Attachments []string       `json:"attachments,omitempty"`
	DelegatedTo *string        `json:"delegated_to,omitempty"`
}

// ApprovalStatistics 审批统计信息
type ApprovalStatistics struct {
	TotalApprovals     int     `json:"total_approvals"`
	PendingApprovals   int     `json:"pending_approvals"`
	ApprovedCount      int     `json:"approved_count"`
	RejectedCount      int     `json:"rejected_count"`
	WithdrawnCount     int     `json:"withdrawn_count"`
	ExpiredCount       int     `json:"expired_count"`
	ApprovalRate       float64 `json:"approval_rate"`
	AverageProcessTime float64 `json:"average_process_time"` // 小时
	OverdueCount       int     `json:"overdue_count"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	EmailEnabled    bool `json:"email_enabled"`
	SMSEnabled      bool `json:"sms_enabled"`
	PushEnabled     bool `json:"push_enabled"`
	ReminderHours   int  `json:"reminder_hours"`   // 提醒间隔小时
	EscalationHours int  `json:"escalation_hours"` // 升级间隔小时
}
