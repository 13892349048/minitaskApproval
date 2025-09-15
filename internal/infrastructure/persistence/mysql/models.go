package mysql

import (
	"time"

	"gorm.io/gorm"
)

// ================================================
// 用户相关模型
// ================================================

// Role 角色模型
type Role struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"type:varchar(100);not null" json:"display_name"`
	Description *string   `gorm:"type:text" json:"description"`
	IsSystem    bool      `gorm:"default:false" json:"is_system"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	Users       []UserModel  `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission 权限模型
type Permission struct {
	ID          string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Resource    string    `gorm:"type:varchar(50);not null" json:"resource"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// UserRole 用户角色关联模型
type UserRole struct {
	UserID     string    `gorm:"type:varchar(36);primaryKey" json:"user_id"`
	RoleID     string    `gorm:"type:varchar(36);primaryKey" json:"role_id"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	AssignedBy *string   `gorm:"type:varchar(36)" json:"assigned_by"`

	// 关联关系
	User     UserModel  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role     Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Assigner *UserModel `gorm:"foreignKey:AssignedBy" json:"assigner,omitempty"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	RoleID       string `gorm:"type:varchar(36);primaryKey" json:"role_id"`
	PermissionID string `gorm:"type:varchar(36);primaryKey" json:"permission_id"`

	// 关联关系
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// PermissionPolicy ABAC权限策略模型
type PermissionPolicy struct {
	ID           string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`
	Description  *string        `gorm:"type:text" json:"description"`
	ResourceType string         `gorm:"type:varchar(50);not null" json:"resource_type"`
	Action       string         `gorm:"type:varchar(50);not null" json:"action"`
	Effect       string         `gorm:"type:enum('allow','deny');not null" json:"effect"`
	Conditions   string         `gorm:"type:json;not null" json:"conditions"`
	Priority     int            `gorm:"default:0" json:"priority"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ================================================
// 项目相关模型
// ================================================

// Project 项目模型
type Project struct {
	ID              string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name            string         `gorm:"type:varchar(200);not null" json:"name"`
	Description     *string        `gorm:"type:text" json:"description"`
	ProjectType     string         `gorm:"type:enum('master','sub','temporary');not null" json:"project_type"`
	ParentProjectID *string        `gorm:"type:varchar(36)" json:"parent_project_id"`
	OwnerID         string         `gorm:"type:varchar(36);not null" json:"owner_id"`
	ManagerID       *string        `gorm:"type:varchar(36)" json:"manager_id"`
	Status          string         `gorm:"type:enum('draft','active','paused','completed','cancelled');default:'draft'" json:"status"`
	StartDate       *time.Time     `gorm:"type:date" json:"start_date"`
	EndDate         *time.Time     `gorm:"type:date" json:"end_date"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	ParentProject *Project        `gorm:"foreignKey:ParentProjectID" json:"parent_project,omitempty"`
	ChildProjects []Project       `gorm:"foreignKey:ParentProjectID" json:"child_projects,omitempty"`
	Owner         UserModel       `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Manager       *UserModel      `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Members       []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Tasks         []Task          `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// ProjectMember 项目成员模型
type ProjectMember struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ProjectID string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_project_user" json:"project_id"`
	UserID    string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_project_user" json:"user_id"`
	Role      string    `gorm:"type:enum('manager','member');not null" json:"role"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
	AddedBy   *string   `gorm:"type:varchar(36)" json:"added_by"`

	// 关联关系
	Project Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	User    UserModel  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Adder   *UserModel `gorm:"foreignKey:AddedBy" json:"adder,omitempty"`
}

// ================================================
// 任务相关模型
// ================================================

// Task 任务模型
type Task struct {
	ID             string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title          string         `gorm:"type:varchar(300);not null" json:"title"`
	Description    *string        `gorm:"type:text" json:"description"`
	TaskType       string         `gorm:"type:enum('single_execution','recurring');not null" json:"task_type"`
	Priority       string         `gorm:"type:enum('low','normal','high','urgent');default:'normal'" json:"priority"`
	ProjectID      string         `gorm:"type:varchar(36);not null" json:"project_id"`
	CreatorID      string         `gorm:"type:varchar(36);not null" json:"creator_id"`
	ResponsibleID  string         `gorm:"type:varchar(36);not null" json:"responsible_id"`
	Status         string         `gorm:"type:enum('draft','pending_approval','approved','in_progress','pending_final_review','completed','rejected','cancelled','paused');default:'draft'" json:"status"`
	StartDate      *time.Time     `gorm:"type:timestamp" json:"start_date"`
	DueDate        *time.Time     `gorm:"type:timestamp" json:"due_date"`
	CompletedAt    *time.Time     `gorm:"type:timestamp" json:"completed_at"`
	EstimatedHours int            `gorm:"default:0" json:"estimated_hours"`
	WorkflowID     *string        `gorm:"type:varchar(36)" json:"workflow_id"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Project          Project            `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Creator          UserModel          `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Responsible      UserModel          `gorm:"foreignKey:ResponsibleID" json:"responsible,omitempty"`
	Participants     []TaskParticipant  `gorm:"foreignKey:TaskID" json:"participants,omitempty"`
	Executions       []TaskExecution    `gorm:"foreignKey:TaskID" json:"executions,omitempty"`
	RecurrenceRule   *RecurrenceRule    `gorm:"foreignKey:TaskID" json:"recurrence_rule,omitempty"`
	Approvals        []ApprovalRecord   `gorm:"foreignKey:TaskID" json:"approvals,omitempty"`
	Extensions       []ExtensionRequest `gorm:"foreignKey:TaskID" json:"extensions,omitempty"`
	FileAssociations []FileAssociation  `gorm:"foreignKey:ResourceID;foreignKey:ResourceType" json:"file_associations,omitempty"`
}

// TaskParticipant 任务参与人员模型
type TaskParticipant struct {
	ID      string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID  string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_task_user" json:"task_id"`
	UserID  string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_task_user" json:"user_id"`
	AddedAt time.Time `gorm:"autoCreateTime" json:"added_at"`
	AddedBy string    `gorm:"type:varchar(36);not null" json:"added_by"`

	// 关联关系
	Task  Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User  UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Adder UserModel `gorm:"foreignKey:AddedBy" json:"adder,omitempty"`
}

// RecurrenceRule 重复任务规则模型
type RecurrenceRule struct {
	ID            string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID        string     `gorm:"type:varchar(36);not null;uniqueIndex" json:"task_id"`
	Frequency     string     `gorm:"type:enum('daily','weekly','monthly');not null" json:"frequency"`
	IntervalValue int        `gorm:"default:1" json:"interval_value"`
	EndDate       *time.Time `gorm:"type:timestamp" json:"end_date"`
	MaxExecutions *int       `gorm:"type:int" json:"max_executions"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	Task Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

// TaskExecution 任务执行记录模型
type TaskExecution struct {
	ID            string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID        string     `gorm:"type:varchar(36);not null" json:"task_id"`
	ExecutionDate time.Time  `gorm:"type:timestamp;not null" json:"execution_date"`
	Status        string     `gorm:"type:enum('pending','in_progress','pending_review','pending_final_review','completed','rejected','cancelled');default:'pending'" json:"status"`
	StartedAt     *time.Time `gorm:"type:timestamp" json:"started_at"`
	SubmittedAt   *time.Time `gorm:"type:timestamp" json:"submitted_at"`
	CompletedAt   *time.Time `gorm:"type:timestamp" json:"completed_at"`
	Result        *string    `gorm:"type:text" json:"result"`

	// 关联关系
	Task                   Task                    `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	ParticipantCompletions []ParticipantCompletion `gorm:"foreignKey:ExecutionID" json:"participant_completions,omitempty"`
}

// ParticipantCompletion 参与人员完成记录模型
type ParticipantCompletion struct {
	ID            string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	ExecutionID   string     `gorm:"type:varchar(36);not null;uniqueIndex:idx_execution_participant" json:"execution_id"`
	ParticipantID string     `gorm:"type:varchar(36);not null;uniqueIndex:idx_execution_participant" json:"participant_id"`
	WorkResult    *string    `gorm:"type:text" json:"work_result"`
	Status        string     `gorm:"type:enum('pending','submitted','approved','rejected');default:'pending'" json:"status"`
	SubmittedAt   *time.Time `gorm:"type:timestamp" json:"submitted_at"`
	ReviewedAt    *time.Time `gorm:"type:timestamp" json:"reviewed_at"`
	ReviewerID    *string    `gorm:"type:varchar(36)" json:"reviewer_id"`
	ReviewComment *string    `gorm:"type:text" json:"review_comment"`

	// 关联关系
	Execution   TaskExecution `gorm:"foreignKey:ExecutionID" json:"execution,omitempty"`
	Participant UserModel     `gorm:"foreignKey:ParticipantID" json:"participant,omitempty"`
	Reviewer    *UserModel    `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// ================================================
// 审批和延期相关模型
// ================================================

// ApprovalRecord 审批记录模型
type ApprovalRecord struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID       string    `gorm:"type:varchar(36);not null" json:"task_id"`
	ExecutionID  *string   `gorm:"type:varchar(36)" json:"execution_id"`
	ApproverID   string    `gorm:"type:varchar(36);not null" json:"approver_id"`
	ApprovalType string    `gorm:"type:enum('task_creation','task_completion','extension_request');not null" json:"approval_type"`
	Action       string    `gorm:"type:enum('approve','reject');not null" json:"action"`
	Comment      *string   `gorm:"type:text" json:"comment"`
	ApprovedAt   time.Time `gorm:"autoCreateTime" json:"approved_at"`

	// 关联关系
	Task      Task           `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Execution *TaskExecution `gorm:"foreignKey:ExecutionID" json:"execution,omitempty"`
	Approver  UserModel      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

// ExtensionRequest 延期申请模型
type ExtensionRequest struct {
	ID               string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	TaskID           string     `gorm:"type:varchar(36);not null" json:"task_id"`
	RequesterID      string     `gorm:"type:varchar(36);not null" json:"requester_id"`
	OriginalDueDate  time.Time  `gorm:"type:timestamp;not null" json:"original_due_date"`
	RequestedDueDate time.Time  `gorm:"type:timestamp;not null" json:"requested_due_date"`
	Reason           string     `gorm:"type:text;not null" json:"reason"`
	Status           string     `gorm:"type:enum('pending','approved','rejected');default:'pending'" json:"status"`
	RequestedAt      time.Time  `gorm:"autoCreateTime" json:"requested_at"`
	ReviewedAt       *time.Time `gorm:"type:timestamp" json:"reviewed_at"`
	ReviewerID       *string    `gorm:"type:varchar(36)" json:"reviewer_id"`
	ReviewComment    *string    `gorm:"type:text" json:"review_comment"`

	// 关联关系
	Task      Task       `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Requester UserModel  `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	Reviewer  *UserModel `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// ================================================
// 事件和日志相关模型
// ================================================

// DomainEvent 领域事件模型
type DomainEvent struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	EventType     string    `gorm:"type:varchar(100);not null" json:"event_type"`
	AggregateID   string    `gorm:"type:varchar(36);not null" json:"aggregate_id"`
	AggregateType string    `gorm:"type:varchar(50);not null" json:"aggregate_type"`
	EventData     string    `gorm:"type:json;not null" json:"event_data"`
	EventVersion  int       `gorm:"default:1" json:"event_version"`
	OccurredAt    time.Time `gorm:"autoCreateTime" json:"occurred_at"`
	UserID        *string   `gorm:"type:varchar(36)" json:"user_id"`

	// 关联关系
	User *UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID         *string   `gorm:"type:varchar(36)" json:"user_id"`
	Operation      string    `gorm:"type:varchar(100);not null" json:"operation"`
	ResourceType   string    `gorm:"type:varchar(50);not null" json:"resource_type"`
	ResourceID     string    `gorm:"type:varchar(36);not null" json:"resource_id"`
	IPAddress      *string   `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent      *string   `gorm:"type:text" json:"user_agent"`
	RequestData    *string   `gorm:"type:json" json:"request_data"`
	ResponseStatus *int      `gorm:"type:int" json:"response_status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	User *UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ================================================
// 文件相关模型
// ================================================

// File 文件模型
type File struct {
	ID           string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null" json:"filename"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FileType     string         `gorm:"type:varchar(50);not null" json:"file_type"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	FilePath     string         `gorm:"type:varchar(500);not null" json:"file_path"`
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	MD5Hash      string         `gorm:"type:varchar(32);not null" json:"md5_hash"`
	UploaderID   string         `gorm:"type:varchar(36);not null" json:"uploader_id"`
	UploadStatus string         `gorm:"type:enum('uploading','completed','failed');default:'uploading'" json:"upload_status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Uploader     UserModel         `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
	Associations []FileAssociation `gorm:"foreignKey:FileID" json:"associations,omitempty"`
}

// FileAssociation 文件关联模型
type FileAssociation struct {
	ID              string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	FileID          string    `gorm:"type:varchar(36);not null" json:"file_id"`
	ResourceType    string    `gorm:"type:varchar(50);not null" json:"resource_type"`
	ResourceID      string    `gorm:"type:varchar(36);not null" json:"resource_id"`
	AssociationType string    `gorm:"type:enum('attachment','avatar','document');not null" json:"association_type"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联关系
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

// ================================================
// 表名映射
// ================================================

func (Role) TableName() string                  { return "roles" }
func (Permission) TableName() string            { return "permissions" }
func (UserRole) TableName() string              { return "user_roles" }
func (RolePermission) TableName() string        { return "role_permissions" }
func (PermissionPolicy) TableName() string      { return "permission_policies" }
func (Project) TableName() string               { return "projects" }
func (ProjectMember) TableName() string         { return "project_members" }
func (Task) TableName() string                  { return "tasks" }
func (TaskParticipant) TableName() string       { return "task_participants" }
func (RecurrenceRule) TableName() string        { return "recurrence_rules" }
func (TaskExecution) TableName() string         { return "task_executions" }
func (ParticipantCompletion) TableName() string { return "participant_completions" }
func (ApprovalRecord) TableName() string        { return "approval_records" }
func (ExtensionRequest) TableName() string      { return "extension_requests" }
func (DomainEvent) TableName() string           { return "domain_events" }
func (OperationLog) TableName() string          { return "operation_logs" }
func (File) TableName() string                  { return "files" }
func (FileAssociation) TableName() string       { return "file_associations" }

// ================================================
// 模型切片类型定义（用于批量操作）
// ================================================

type Users []UserModel
type Roles []Role
type Permissions []Permission
type Projects []Project
type Tasks []Task
type TaskExecutions []TaskExecution
