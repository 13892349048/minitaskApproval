package valueobject

import (
	"regexp"
	"time"
)

// 用户扩展值对象和数据传输对象

// Phone 电话号码值对象
type Phone struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

// String 返回完整电话号码
func (p Phone) String() string {
	if p.CountryCode == "" {
		return p.Number
	}
	return p.CountryCode + p.Number
}

// IsValid 验证电话号码格式
func (p Phone) IsValid() bool {
	if p.Number == "" {
		return false
	}
	// 简单的电话号码验证
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return phoneRegex.MatchString(p.String())
}

// Email 邮箱值对象
type Email struct {
	Address string `json:"address"`
}

// String 返回邮箱地址
func (e Email) String() string {
	return e.Address
}

// IsValid 验证邮箱格式
func (e Email) IsValid() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e.Address)
}

// Department 部门值对象
type Department struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id,omitempty"`
}

// UserProfile 用户档案值对象
type UserProfile struct {
	FullName     string      `json:"full_name"`
	Phone        *Phone      `json:"phone,omitempty"`
	Department   *Department `json:"department,omitempty"`
	Position     string      `json:"position,omitempty"`
	ManagerID    *UserID     `json:"manager_id,omitempty"`
	Avatar       string      `json:"avatar,omitempty"`
	Bio          string      `json:"bio,omitempty"`
	Location     string      `json:"location,omitempty"`
	Timezone     string      `json:"timezone,omitempty"`
	Language     string      `json:"language,omitempty"`
	JoinDate     time.Time   `json:"join_date"`
	LastLoginAt  *time.Time  `json:"last_login_at,omitempty"`
}

// UserCredentials 用户凭证值对象
type UserCredentials struct {
	Email        Email     `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Salt         string    `json:"salt,omitempty"`
	LastChanged  time.Time `json:"last_changed"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// UserPreferences 用户偏好设置
type UserPreferences struct {
	Theme             string            `json:"theme"`
	Language          string            `json:"language"`
	Timezone          string            `json:"timezone"`
	NotificationEmail bool              `json:"notification_email"`
	NotificationSMS   bool              `json:"notification_sms"`
	NotificationPush  bool              `json:"notification_push"`
	CustomSettings    map[string]string `json:"custom_settings,omitempty"`
}

// UserSecuritySettings 用户安全设置
type UserSecuritySettings struct {
	TwoFactorEnabled    bool       `json:"two_factor_enabled"`
	TwoFactorSecret     string     `json:"two_factor_secret,omitempty"`
	BackupCodes         []string   `json:"backup_codes,omitempty"`
	LastPasswordChange  time.Time  `json:"last_password_change"`
	PasswordExpiresAt   *time.Time `json:"password_expires_at,omitempty"`
	LoginAttempts       int        `json:"login_attempts"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	SessionTimeout      int        `json:"session_timeout"` // 分钟
}

// 角色常量修正 - 与user.go保持一致
const (
	RoleEmployee = "employee"
	RoleManager  = "manager"
	RoleDirector = "director"
	RoleAdmin    = "admin"
)

// UserData 用户数据传输对象（用于持久化和恢复）
type UserData struct {
	ID               string                `json:"id"`
	Username         string                `json:"username"`
	Email            string                `json:"email"`
	PasswordHash     string                `json:"password_hash"`
	FullName         string                `json:"full_name"`
	Phone            *string               `json:"phone,omitempty"`
	Status           string                `json:"status"`
	Role             string                `json:"role"`
	DepartmentID     *string               `json:"department_id,omitempty"`
	ManagerID        *string               `json:"manager_id,omitempty"`
	Position         *string               `json:"position,omitempty"`
	Avatar           *string               `json:"avatar,omitempty"`
	Bio              *string               `json:"bio,omitempty"`
	Location         *string               `json:"location,omitempty"`
	Timezone         *string               `json:"timezone,omitempty"`
	Language         *string               `json:"language,omitempty"`
	JoinDate         time.Time             `json:"join_date"`
	LastLoginAt      *time.Time            `json:"last_login_at,omitempty"`
	Preferences      *UserPreferences      `json:"preferences,omitempty"`
	SecuritySettings *UserSecuritySettings `json:"security_settings,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	DeletedAt        *time.Time            `json:"deleted_at,omitempty"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []UserSummary       `json:"users"`
	Pagination PaginationResponse  `json:"pagination"`
}

// UserSummary 用户摘要信息
type UserSummary struct {
	ID         string     `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	Phone      *string    `json:"phone,omitempty"`
	Status     string     `json:"status"`
	Role       string     `json:"role"`
	Department *string    `json:"department,omitempty"`
	Position   *string    `json:"position,omitempty"`
	Avatar     *string    `json:"avatar,omitempty"`
	JoinDate   time.Time  `json:"join_date"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
}

// UserDetailResponse 用户详细信息响应
type UserDetailResponse struct {
	ID               string                `json:"id"`
	Username         string                `json:"username"`
	Email            string                `json:"email"`
	FullName         string                `json:"full_name"`
	Phone            *string               `json:"phone,omitempty"`
	Status           string                `json:"status"`
	Roles            []string              `json:"roles"`
	Department       *Department           `json:"department,omitempty"`
	Manager          *UserSummary          `json:"manager,omitempty"`
	Position         *string               `json:"position,omitempty"`
	Avatar           *string               `json:"avatar,omitempty"`
	Bio              *string               `json:"bio,omitempty"`
	Location         *string               `json:"location,omitempty"`
	JoinDate         time.Time             `json:"join_date"`
	LastLogin        *time.Time            `json:"last_login,omitempty"`
	Preferences      *UserPreferences      `json:"preferences,omitempty"`
	SecuritySettings *UserSecuritySettings `json:"security_settings,omitempty"`
	Statistics       *UserStatistics       `json:"statistics,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	TotalTasks        int     `json:"total_tasks"`
	CompletedTasks    int     `json:"completed_tasks"`
	PendingTasks      int     `json:"pending_tasks"`
	OverdueTasks      int     `json:"overdue_tasks"`
	TotalProjects     int     `json:"total_projects"`
	ActiveProjects    int     `json:"active_projects"`
	CompletionRate    float64 `json:"completion_rate"`
	AverageTaskTime   float64 `json:"average_task_time"`
	LoginCount        int     `json:"login_count"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty"`
}

// UserSearchRequest 用户搜索请求
type UserSearchRequest struct {
	BaseSearchCriteria
	Username     *string     `json:"username,omitempty"`
	Email        *string     `json:"email,omitempty"`
	FullName     *string     `json:"full_name,omitempty"`
	Role         *UserRole   `json:"role,omitempty"`
	Status       *UserStatus `json:"status,omitempty"`
	DepartmentID *string     `json:"department_id,omitempty"`
	ManagerID    *UserID     `json:"manager_id,omitempty"`
	Position     *string     `json:"position,omitempty"`
}

// AuthenticationRequest 认证请求
type AuthenticationRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	RememberMe bool   `json:"remember_me"`
	DeviceInfo string `json:"device_info,omitempty"`
}

// AuthenticationResponse 认证响应
type AuthenticationResponse struct {
	User         UserDetailResponse `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	TokenType    string            `json:"token_type"`
}

// PasswordChangeRequest 密码修改请求
type PasswordChangeRequest struct {
	UserID          string `json:"user_id" validate:"required"`
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// ProfileUpdateRequest 档案更新请求
type ProfileUpdateRequest struct {
	UserID      string  `json:"user_id" validate:"required"`
	FullName    *string `json:"full_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Position    *string `json:"position,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Location    *string `json:"location,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

// UserRoleAssignmentRequest 用户角色分配请求
type UserRoleAssignmentRequest struct {
	UserID    string   `json:"user_id" validate:"required"`
	Roles     []string `json:"roles" validate:"required,min=1"`
	AssignedBy string  `json:"assigned_by" validate:"required"`
	Reason    string   `json:"reason,omitempty"`
}
