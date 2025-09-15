package dto

import (
	"time"
	"github.com/taskflow/internal/domain/valueobject"
)

// 用户应用层数据传输对象

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	FullName    string `json:"full_name" validate:"required,min=2,max=100"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,phone"`
	Position    string `json:"position,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	ManagerID   string `json:"manager_id,omitempty"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID          string  `json:"id" validate:"required"`
	FullName    *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,phone"`
	Position    *string `json:"position,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Location    *string `json:"location,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID           string                        `json:"id"`
	Username     string                        `json:"username"`
	Email        string                        `json:"email"`
	FullName     string                        `json:"full_name"`
	Phone        *string                       `json:"phone,omitempty"`
	Status       string                        `json:"status"`
	Roles        []string                      `json:"roles"`
	Department   *valueobject.Department       `json:"department,omitempty"`
	Manager      *UserSummary                  `json:"manager,omitempty"`
	Position     *string                       `json:"position,omitempty"`
	Avatar       *string                       `json:"avatar,omitempty"`
	Bio          *string                       `json:"bio,omitempty"`
	Location     *string                       `json:"location,omitempty"`
	JoinDate     time.Time                     `json:"join_date"`
	LastLogin    *time.Time                    `json:"last_login,omitempty"`
	Preferences  *valueobject.UserPreferences  `json:"preferences,omitempty"`
	Statistics   *valueobject.UserStatistics   `json:"statistics,omitempty"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// UserSummary 用户摘要
type UserSummary struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar,omitempty"`
	Position *string `json:"position,omitempty"`
	Status   string  `json:"status"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	valueobject.BaseSearchCriteria
	Username     *string                   `json:"username,omitempty"`
	Email        *string                   `json:"email,omitempty"`
	FullName     *string                   `json:"full_name,omitempty"`
	Role         *valueobject.UserRole     `json:"role,omitempty"`
	Status       *valueobject.UserStatus   `json:"status,omitempty"`
	DepartmentID *string                   `json:"department_id,omitempty"`
	ManagerID    *valueobject.UserID       `json:"manager_id,omitempty"`
	Position     *string                   `json:"position,omitempty"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []UserSummary                   `json:"users"`
	Pagination valueobject.PaginationResponse  `json:"pagination"`
}

// AuthenticationRequest 认证请求
type AuthenticationRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
	DeviceInfo string `json:"device_info,omitempty"`
}

// AuthenticationResponse 认证响应
type AuthenticationResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	TokenType    string       `json:"token_type"`
}

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	UserID          string `json:"user_id" validate:"required"`
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// UserRoleAssignmentRequest 用户角色分配请求
type UserRoleAssignmentRequest struct {
	UserID     string   `json:"user_id" validate:"required"`
	Roles      []string `json:"roles" validate:"required,min=1"`
	AssignedBy string   `json:"assigned_by" validate:"required"`
	Reason     string   `json:"reason,omitempty"`
}

// UserStatusChangeRequest 用户状态变更请求
type UserStatusChangeRequest struct {
	UserID    string                    `json:"user_id" validate:"required"`
	Status    valueobject.UserStatus    `json:"status" validate:"required"`
	ChangedBy string                    `json:"changed_by" validate:"required"`
	Reason    string                    `json:"reason,omitempty"`
}

// UserPreferencesUpdateRequest 用户偏好设置更新请求
type UserPreferencesUpdateRequest struct {
	UserID                string            `json:"user_id" validate:"required"`
	Theme                 *string           `json:"theme,omitempty"`
	Language              *string           `json:"language,omitempty"`
	Timezone              *string           `json:"timezone,omitempty"`
	NotificationEmail     *bool             `json:"notification_email,omitempty"`
	NotificationSMS       *bool             `json:"notification_sms,omitempty"`
	NotificationPush      *bool             `json:"notification_push,omitempty"`
	CustomSettings        map[string]string `json:"custom_settings,omitempty"`
}

// UserStatisticsResponse 用户统计响应
type UserStatisticsResponse struct {
	UserID              string     `json:"user_id"`
	TotalTasks          int        `json:"total_tasks"`
	CompletedTasks      int        `json:"completed_tasks"`
	PendingTasks        int        `json:"pending_tasks"`
	OverdueTasks        int        `json:"overdue_tasks"`
	TotalProjects       int        `json:"total_projects"`
	ActiveProjects      int        `json:"active_projects"`
	CompletionRate      float64    `json:"completion_rate"`
	AverageTaskTime     float64    `json:"average_task_time"`
	LoginCount          int        `json:"login_count"`
	LastActivityAt      *time.Time `json:"last_activity_at,omitempty"`
	ProductivityScore   float64    `json:"productivity_score"`
	CollaborationScore  float64    `json:"collaboration_score"`
}
