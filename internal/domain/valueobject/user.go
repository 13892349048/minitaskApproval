package valueobject

import (
	"time"
)

// UserID 用户ID值对象
type UserID string

func (id UserID) String() string {
	return string(id)
}

// UserRole 用户角色
type UserRole string

const (
	UserRoleEmployee    UserRole = "employee"
	UserRoleManager     UserRole = "manager"
	UserRoleDirector    UserRole = "director"
	UserRoleAdmin       UserRole = "admin"
	UserRoleSuperAdmin  UserRole = "super_admin"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// UserSearchCriteria 用户搜索条件
type UserSearchCriteria struct {
	Username     *string     `json:"username"`
	Email        *string     `json:"email"`
	FullName     *string     `json:"full_name"`
	Role         *UserRole   `json:"role"`
	Status       *UserStatus `json:"status"`
	DepartmentID *string     `json:"department_id"`
	ManagerID    *UserID     `json:"manager_id"`
	Limit        int         `json:"limit"`
	Offset       int         `json:"offset"`
	OrderBy      string      `json:"order_by"`
	OrderDir     string      `json:"order_dir"`
}

// TaskStatistics 任务统计信息
type TaskStatistics struct {
	TaskID               TaskID     `json:"task_id"`
	TotalExecutions      int        `json:"total_executions"`
	CompletedExecutions  int        `json:"completed_executions"`
	PendingExecutions    int        `json:"pending_executions"`
	RejectedExecutions   int        `json:"rejected_executions"`
	TotalParticipants    int        `json:"total_participants"`
	ActiveParticipants   int        `json:"active_participants"`
	CompletionRate       float64    `json:"completion_rate"`
	AverageExecutionTime float64    `json:"average_execution_time"`
	LastExecutionAt      time.Time  `json:"last_execution_at"`
	NextExecutionAt      *time.Time `json:"next_execution_at"`
}

// WorkflowStepData 工作流步骤数据
type WorkflowStepData struct {
	StepID  string                 `json:"step_id"`
	Action  string                 `json:"action"`
	ActorID UserID                 `json:"actor_id"`
	Data    map[string]interface{} `json:"data"`
	Comment string                 `json:"comment"`
}

// TaskValidator 任务验证器接口
type TaskValidator interface {
	ValidateTitle(title string) error
	ValidateDescription(description string) error
	ValidateDueDate(dueDate *time.Time) error
	ValidateEstimatedHours(hours int) error
}

// TaskPermissions 任务权限
type TaskPermissions struct {
	CanView    bool `json:"can_view"`
	CanModify  bool `json:"can_modify"`
	CanExecute bool `json:"can_execute"`
	CanApprove bool `json:"can_approve"`
	CanDelete  bool `json:"can_delete"`
	CanAssign  bool `json:"can_assign"`
}
