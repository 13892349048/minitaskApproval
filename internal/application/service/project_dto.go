package service

import (
	"time"

	"github.com/taskflow/internal/domain/valueobject"
)

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	ProjectType string `json:"project_type" binding:"required,oneof=master sub"`
	OwnerID     string `json:"owner_id" binding:"required"`
	ParentID    string `json:"parent_id,omitempty"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID          string                        `json:"id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	ProjectType string                        `json:"project_type"`
	Status      string                        `json:"status"`
	OwnerID     string                        `json:"owner_id"`
	ManagerID   *string                       `json:"manager_id,omitempty"`
	ParentID    *string                       `json:"parent_id,omitempty"`
	Members     []ProjectMemberResponse       `json:"members"`
	Children    []string                      `json:"children"`
	StartDate   time.Time                     `json:"start_date"`
	EndDate     *time.Time                    `json:"end_date,omitempty"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	Statistics  *ProjectStatisticsResponse    `json:"statistics,omitempty"`
}

// ProjectMemberResponse 项目成员响应
type ProjectMemberResponse struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	AddedBy  string    `json:"added_by"`
}

// ProjectStatisticsResponse 项目统计响应
type ProjectStatisticsResponse struct {
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	PendingTasks   int `json:"pending_tasks"`
	TotalMembers   int `json:"total_members"`
	Progress       int `json:"progress_percentage"`
}

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=member developer tester"`
}

// UpdateMemberRoleRequest 更新成员角色请求
type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=member developer tester"`
}

// AssignManagerRequest 分配管理者请求
type AssignManagerRequest struct {
	ManagerID string `json:"manager_id" binding:"required"`
}

// ChangeStatusRequest 更改状态请求
type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft active paused completed cancelled"`
	Reason string `json:"reason,omitempty"`
}

// ProjectListRequest 项目列表请求
type ProjectListRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Status     string `form:"status,omitempty" binding:"omitempty,oneof=draft active paused completed cancelled"`
	Type       string `form:"type,omitempty" binding:"omitempty,oneof=master sub"`
	OwnerID    string `form:"owner_id,omitempty"`
	ManagerID  string `form:"manager_id,omitempty"`
	Search     string `form:"search,omitempty"`
	SortBy     string `form:"sort_by,default=created_at" binding:"omitempty,oneof=name created_at updated_at status"`
	SortOrder  string `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// CreateSubProjectRequest 创建子项目请求
type CreateSubProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ProjectHierarchyResponse 项目层级响应
type ProjectHierarchyResponse struct {
	Project       *ProjectResponse   `json:"project"`
	Parent        *ProjectResponse   `json:"parent,omitempty"`
	Children      []ProjectResponse  `json:"children"`
	Depth         int                `json:"depth"`
	TotalProjects int                `json:"total_projects"`
}

// 转换函数

// ToProjectMemberResponse 转换项目成员响应
func ToProjectMemberResponse(member valueobject.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		UserID:   string(member.UserID),
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
		AddedBy:  string(member.AddedBy),
	}
}

// ToProjectStatisticsResponse 转换项目统计响应
func ToProjectStatisticsResponse(totalTasks, completedTasks, totalMembers int) *ProjectStatisticsResponse {
	pendingTasks := totalTasks - completedTasks
	progress := 0
	if totalTasks > 0 {
		progress = (completedTasks * 100) / totalTasks
	}

	return &ProjectStatisticsResponse{
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
		PendingTasks:   pendingTasks,
		TotalMembers:   totalMembers,
		Progress:       progress,
	}
}
