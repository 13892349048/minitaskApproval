package valueobject

import (
	"time"
)

// TaskID 任务ID值对象
type TaskID string

func (id TaskID) String() string {
	return string(id)
}

// TaskType 任务类型
type TaskType string

const (
	TaskTypeRegular   TaskType = "regular"   // 常规任务
	TaskTypeRecurring TaskType = "recurring" // 重复任务
	TaskTypeTemplate  TaskType = "template"  // 模板任务
	TaskTypeUrgent    TaskType = "urgent"    // 紧急任务
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusDraft           TaskStatus = "draft"            // 草稿
	TaskStatusPendingApproval TaskStatus = "pending_approval" // 待审批
	TaskStatusApproved        TaskStatus = "approved"         // 已审批
	TaskStatusRejected        TaskStatus = "rejected"         // 已拒绝
	TaskStatusInProgress      TaskStatus = "in_progress"      // 进行中
	TaskStatusPaused          TaskStatus = "paused"           // 已暂停
	TaskStatusCompleted       TaskStatus = "completed"        // 已完成
	TaskStatusCancelled       TaskStatus = "cancelled"        // 已取消
)

// TaskPriority 任务优先级
type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"      // 低优先级
	TaskPriorityMedium   TaskPriority = "medium"   // 中优先级
	TaskPriorityHigh     TaskPriority = "high"     // 高优先级
	TaskPriorityCritical TaskPriority = "critical" // 紧急优先级
)

// RecurrenceFrequency 重复频率
type RecurrenceFrequency string

const (
	RecurrenceDaily   RecurrenceFrequency = "daily"   // 每日
	RecurrenceWeekly  RecurrenceFrequency = "weekly"  // 每周
	RecurrenceMonthly RecurrenceFrequency = "monthly" // 每月
	RecurrenceYearly  RecurrenceFrequency = "yearly"  // 每年
)

// TaskExecutionID 任务执行ID
type TaskExecutionID string

func (id TaskExecutionID) String() string {
	return string(id)
}

// ExtensionRequestID 延期请求ID
type ExtensionRequestID string

func (id ExtensionRequestID) String() string {
	return string(id)
}

// ParticipantRole 参与者角色
type ParticipantRole string

const (
	ParticipantRoleExecutor  ParticipantRole = "executor"  // 执行者
	ParticipantRoleReviewer  ParticipantRole = "reviewer"  // 审核者
	ParticipantRoleObserver  ParticipantRole = "observer"  // 观察者
	ParticipantRoleAssistant ParticipantRole = "assistant" // 协助者
)

// TaskParticipant 任务参与者值对象
type TaskParticipant struct {
	UserID  UserID          `json:"user_id"`
	Role    ParticipantRole `json:"role"`
	AddedAt time.Time       `json:"added_at"`
	AddedBy UserID          `json:"added_by"`
}

// TaskSearchCriteria 任务搜索条件
type TaskSearchCriteria struct {
	Title         *string       `json:"title"`
	Description   *string       `json:"description"`
	TaskType      *TaskType     `json:"task_type"`
	Priority      *TaskPriority `json:"priority"`
	Status        *TaskStatus   `json:"status"`
	ProjectID     *ProjectID    `json:"project_id"`
	CreatorID     *UserID       `json:"creator_id"`
	ResponsibleID *UserID       `json:"responsible_id"`
	ParticipantID *UserID       `json:"participant_id"`
	StartDate     *time.Time    `json:"start_date"`
	DueDate       *time.Time    `json:"due_date"`
	CreatedAfter  *time.Time    `json:"created_after"`
	CreatedBefore *time.Time    `json:"created_before"`
	Limit         int           `json:"limit"`
	Offset        int           `json:"offset"`
	OrderBy       string        `json:"order_by"`
	OrderDir      string        `json:"order_dir"`
}

// TaskData 任务数据传输对象（用于持久化和恢复）
type TaskData struct {
	ID             string                `json:"id"`
	Title          string                `json:"title"`
	Description    *string               `json:"description"`
	TaskType       string                `json:"task_type"`
	Priority       string                `json:"priority"`
	Status         string                `json:"status"`
	ProjectID      string                `json:"project_id"`
	CreatorID      string                `json:"creator_id"`
	ResponsibleID  string                `json:"responsible_id"`
	StartDate      *time.Time            `json:"start_date"`
	DueDate        *time.Time            `json:"due_date"`
	CompletedAt    *time.Time            `json:"completed_at"`
	EstimatedHours int                   `json:"estimated_hours"`
	WorkflowID     *string               `json:"workflow_id"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	DeletedAt      *time.Time            `json:"deleted_at"`
	Participants   []TaskParticipantData `json:"participants"`
}

// TaskParticipantData 任务参与者数据传输对象
type TaskParticipantData struct {
	UserID  string    `json:"user_id"`
	AddedAt time.Time `json:"added_at"`
	AddedBy string    `json:"added_by"`
}

