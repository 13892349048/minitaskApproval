package service

import (
	"time"
)

// Task related DTOs
type CreateTaskRequest struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	TaskType      string     `json:"task_type"`
	Priority      string     `json:"priority"`
	ProjectID     string     `json:"project_id"`
	CreatorID     string     `json:"creator_id"`
	ResponsibleID string     `json:"responsible_id"`
	DueDate       *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AssignTaskRequest struct {
	TaskID        string `json:"task_id"`
	ResponsibleID string `json:"responsible_id"`
	AssignedBy    string `json:"assigned_by"`
}

// Permission related DTOs
type CheckPermissionRequest struct {
	UserID   string                 `json:"user_id"`
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context"`
}

type PermissionResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

// Common response types
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
