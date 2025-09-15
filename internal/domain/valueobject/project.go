package valueobject

import (
	"time"
)

// ProjectID 项目ID值对象
type ProjectID string

func (id ProjectID) String() string {
	return string(id)
}

// ProjectType 项目类型
type ProjectType string

const (
	ProjectTypeMaster    ProjectType = "master"
	ProjectTypeSub       ProjectType = "sub"
	ProjectTypeTemporary ProjectType = "temporary"
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusPaused    ProjectStatus = "paused"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

// ProjectRole 项目角色
type ProjectRole string

const (
	ProjectRoleManager   ProjectRole = "manager"
	ProjectRoleMember    ProjectRole = "member"
	ProjectRoleDeveloper ProjectRole = "developer"
	ProjectRoleTester    ProjectRole = "tester"
)

// ProjectMember 项目成员值对象
type ProjectMember struct {
	UserID   UserID      `json:"user_id"`
	Role     ProjectRole `json:"role"`
	JoinedAt time.Time   `json:"joined_at"`
	AddedBy  UserID      `json:"added_by"`
}

// ProjectTaskStatistics 项目任务统计信息
type ProjectTaskStatistics struct {
	ProjectID         ProjectID `json:"project_id"`
	TotalTasks        int       `json:"total_tasks"`
	CompletedTasks    int       `json:"completed_tasks"`
	InProgressTasks   int       `json:"in_progress_tasks"`
	PendingTasks      int       `json:"pending_tasks"`
	OverdueTasks      int       `json:"overdue_tasks"`
	HighPriorityTasks int       `json:"high_priority_tasks"`
	CompletionRate    float64   `json:"completion_rate"`
	AverageTaskTime   float64   `json:"average_task_time"`
}
