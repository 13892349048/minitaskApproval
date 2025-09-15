package dto

import (
	"time"
	"github.com/taskflow/internal/domain/valueobject"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title         string    `json:"title" validate:"required"`
	Description   *string   `json:"description"`
	TaskType      string    `json:"task_type" validate:"required"`
	Priority      string    `json:"priority" validate:"required"`
	ProjectID     string    `json:"project_id" validate:"required"`
	CreatorID     string    `json:"creator_id" validate:"required"`
	ResponsibleID string    `json:"responsible_id" validate:"required"`
	DueDate       *time.Time `json:"due_date"`
	EstimatedHours int      `json:"estimated_hours"`
}

// CreateTaskResponse 创建任务响应
type CreateTaskResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	TaskType      string    `json:"task_type"`
	Priority      string    `json:"priority"`
	Status        string    `json:"status"`
	ProjectID     string    `json:"project_id"`
	CreatorID     string    `json:"creator_id"`
	ResponsibleID string    `json:"responsible_id"`
	DueDate       *time.Time `json:"due_date"`
	EstimatedHours int      `json:"estimated_hours"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	ID            string     `json:"id"`
	Title         *string    `json:"title"`
	Description   *string    `json:"description"`
	Priority      *string    `json:"priority"`
	DueDate       *time.Time `json:"due_date"`
	EstimatedHours *int      `json:"estimated_hours"`
}

// UpdateTaskResponse 更新任务响应
type UpdateTaskResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	TaskType      string    `json:"task_type"`
	Priority      string    `json:"priority"`
	Status        string    `json:"status"`
	ProjectID     string    `json:"project_id"`
	CreatorID     string    `json:"creator_id"`
	ResponsibleID string    `json:"responsible_id"`
	DueDate       *time.Time `json:"due_date"`
	EstimatedHours int      `json:"estimated_hours"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID            string                `json:"id"`
	Title         string                `json:"title"`
	Description   *string               `json:"description"`
	TaskType      string                `json:"task_type"`
	Priority      string                `json:"priority"`
	Status        string                `json:"status"`
	ProjectID     string                `json:"project_id"`
	CreatorID     string                `json:"creator_id"`
	ResponsibleID string                `json:"responsible_id"`
	DueDate       *time.Time            `json:"due_date"`
	EstimatedHours int                  `json:"estimated_hours"`
	ActualHours   float64               `json:"actual_hours"`
	Participants  []TaskParticipantDTO  `json:"participants"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// TaskParticipantDTO 任务参与者DTO
type TaskParticipantDTO struct {
	UserID  string    `json:"user_id"`
	Role    string    `json:"role"`
	AddedAt time.Time `json:"added_at"`
	AddedBy string    `json:"added_by"`
}

// TaskSearchCriteria 任务搜索条件
type TaskSearchCriteria struct {
	Title         *string                      `json:"title"`
	Description   *string                      `json:"description"`
	TaskType      *valueobject.TaskType        `json:"task_type"`
	Priority      *valueobject.TaskPriority    `json:"priority"`
	Status        *valueobject.TaskStatus      `json:"status"`
	ProjectID     *valueobject.ProjectID       `json:"project_id"`
	CreatorID     *valueobject.UserID          `json:"creator_id"`
	ResponsibleID *valueobject.UserID          `json:"responsible_id"`
	ParticipantID *valueobject.UserID          `json:"participant_id"`
	StartDate     *time.Time                   `json:"start_date"`
	DueDate       *time.Time                   `json:"due_date"`
	CreatedAfter  *time.Time                   `json:"created_after"`
	CreatedBefore *time.Time                   `json:"created_before"`
}

// ListTasksRequest 任务列表请求
type ListTasksRequest struct {
	Criteria TaskSearchCriteria `json:"criteria"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// ListTasksResponse 任务列表响应
type ListTasksResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	TaskID        string `json:"task_id"`
	ResponsibleID string `json:"responsible_id" validate:"required"`
	AssignedBy    string `json:"assigned_by" validate:"required"`
}

// UpdateTaskStatusRequest 更新任务状态请求
type UpdateTaskStatusRequest struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status" validate:"required"`
	UpdatedBy string `json:"updated_by" validate:"required"`
	Comment   string `json:"comment"`
}

// AddTaskParticipantRequest 添加任务参与者请求
type AddTaskParticipantRequest struct {
	TaskID        string `json:"task_id"`
	ParticipantID string `json:"participant_id" validate:"required"`
	Role          string `json:"role"`
	AddedBy       string `json:"added_by" validate:"required"`
}

// RemoveTaskParticipantRequest 移除任务参与者请求
type RemoveTaskParticipantRequest struct {
	TaskID        string `json:"task_id"`
	ParticipantID string `json:"participant_id"`
	RemovedBy     string `json:"removed_by" validate:"required"`
}

// TaskStatisticsResponse 任务统计响应
type TaskStatisticsResponse struct {
	TotalTasks      int                        `json:"total_tasks"`
	TasksByStatus   map[string]int             `json:"tasks_by_status"`
	TasksByPriority map[string]int             `json:"tasks_by_priority"`
	TasksByType     map[string]int             `json:"tasks_by_type"`
	OverdueTasks    int                        `json:"overdue_tasks"`
	CompletionRate  float64                    `json:"completion_rate"`
	AverageHours    float64                    `json:"average_hours"`
}

// ProjectTaskStatisticsResponse 项目任务统计响应
type ProjectTaskStatisticsResponse struct {
	ProjectID       string                     `json:"project_id"`
	TotalTasks      int                        `json:"total_tasks"`
	TasksByStatus   map[string]int             `json:"tasks_by_status"`
	TasksByPriority map[string]int             `json:"tasks_by_priority"`
	TasksByType     map[string]int             `json:"tasks_by_type"`
	OverdueTasks    int                        `json:"overdue_tasks"`
	CompletionRate  float64                    `json:"completion_rate"`
	AverageHours    float64                    `json:"average_hours"`
}
