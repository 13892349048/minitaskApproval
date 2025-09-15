package aggregate

import (
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/valueobject"
)

// ProjectAggregate 项目聚合根接口
type ProjectAggregate interface {

	// 业务行为方法
	UpdateBasicInfo(name, description string) error
	AssignManager(managerID valueobject.UserID, assignedBy valueobject.UserID) error
	AddMember(userID valueobject.UserID, role valueobject.ProjectRole, addedBy valueobject.UserID) error
	RemoveMember(userID valueobject.UserID, removedBy valueobject.UserID) error
	UpdateMemberRole(userID valueobject.UserID, newRole valueobject.ProjectRole, updatedBy valueobject.UserID) error
	CreateSubProject(subProjectID valueobject.ProjectID, name, description string, createdBy valueobject.UserID) (ProjectAggregate, error)

	// 状态管理
	Activate(activatedBy valueobject.UserID) error
	Pause(pausedBy valueobject.UserID, reason string) error
	Complete(completedBy valueobject.UserID) error
	Cancel(cancelledBy valueobject.UserID, reason string) error
	Delete(deletedBy valueobject.UserID) error

	// 统计和权限
	UpdateTaskStatistics(totalTasks, completedTasks int)
	CanUserAccess(userID valueobject.UserID) bool
	GetMemberRole(userID valueobject.UserID) *valueobject.ProjectRole
	GetMemberIDs() []string

	// 事件管理
	ClearEvents()
}

// ProjectSearchCriteria 项目搜索条件
type ProjectSearchCriteria struct {
	Name        *string
	Description *string
	ProjectType *valueobject.ProjectType
	Status      *valueobject.ProjectStatus
	OwnerID     *valueobject.UserID
	ManagerID   *valueobject.UserID
	MemberID    *valueobject.UserID
	ParentID    *valueobject.ProjectID
	StartDate   *time.Time
	EndDate     *time.Time
	Limit       int
	Offset      int
	OrderBy     string
	OrderDir    string
}

// ProjectStatistics 项目统计信息
type ProjectStatistics struct {
	ProjectID       valueobject.ProjectID `json:"project_id"`
	TotalTasks      int                   `json:"total_tasks"`
	CompletedTasks  int                   `json:"completed_tasks"`
	InProgressTasks int                   `json:"in_progress_tasks"`
	PendingTasks    int                   `json:"pending_tasks"`
	OverdueTasks    int                   `json:"overdue_tasks"`
	TotalMembers    int                   `json:"total_members"`
	ActiveMembers   int                   `json:"active_members"`
	CompletionRate  float64               `json:"completion_rate"`
	AverageTaskTime float64               `json:"average_task_time"`
	LastActivityAt  time.Time             `json:"last_activity_at"`
}

// ProjectFactory 项目工厂 - Go风格：返回具体类型
type ProjectFactory struct {
	// 可以注入依赖，如ID生成器、验证器等
}

// NewProjectFactory 创建项目工厂
func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{}
}

// CreateProject 创建新项目 - 返回具体类型
func (f *ProjectFactory) CreateProject(
	id valueobject.ProjectID,
	name, description string,
	projectType valueobject.ProjectType,
	ownerID valueobject.UserID,
) *Project {
	return NewProject(id, name, description, projectType, ownerID)
}

// CreateSubProject 创建子项目
func (f *ProjectFactory) CreateSubProject(
	parent *Project,
	id valueobject.ProjectID,
	name, description string,
	createdBy valueobject.UserID,
) (*Project, error) {
	subProject, err := parent.CreateSubProject(id, name, description, createdBy)
	if err != nil {
		return nil, err
	}
	// 类型断言，因为我们知道返回的是 *Project
	return subProject.(*Project), nil
}

// RestoreProject 从数据恢复项目
func (f *ProjectFactory) RestoreProject(data ProjectData) *Project {
	// 实现从持久化数据恢复项目逻辑
	project := &Project{
		ID:             valueobject.ProjectID(data.ID),
		Name:           data.Name,
		Description:    data.Description,
		ProjectType:    valueobject.ProjectType(data.Type),
		Status:         valueobject.ProjectStatus(data.Status),
		OwnerID:        valueobject.UserID(data.OwnerID),
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
		DeletedAt:      data.DeletedAt,
		TaskCount:      data.TaskCount,
		CompletedTasks: data.CompletedTasks,
		Events:         make([]event.DomainEvent, 0),
	}

	if data.ParentID != nil {
		parentID := valueobject.ProjectID(*data.ParentID)
		project.ParentID = &parentID
	}

	if data.ManagerID != nil {
		managerID := valueobject.UserID(*data.ManagerID)
		project.ManagerID = &managerID
	}

	// 恢复成员列表
	for _, memberData := range data.Members {
		member := valueobject.ProjectMember{
			UserID:   valueobject.UserID(memberData.UserID),
			Role:     valueobject.ProjectRole(memberData.Role),
			JoinedAt: memberData.JoinedAt,
			AddedBy:  valueobject.UserID(memberData.AddedBy),
		}
		project.Members = append(project.Members, member)
	}

	// 恢复子项目列表
	for _, childID := range data.Children {
		project.Children = append(project.Children, valueobject.ProjectID(childID))
	}

	return project
}

// ProjectData 项目数据传输对象（用于持久化和恢复）
type ProjectData struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	Type           string              `json:"type"`
	Status         string              `json:"status"`
	ParentID       *string             `json:"parent_id"`
	OwnerID        string              `json:"owner_id"`
	ManagerID      *string             `json:"manager_id"`
	StartDate      time.Time           `json:"start_date"`
	EndDate        *time.Time          `json:"end_date"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      *time.Time          `json:"deleted_at"`
	Members        []ProjectMemberData `json:"members"`
	Children       []string            `json:"children"`
	TaskCount      int                 `json:"task_count"`
	CompletedTasks int                 `json:"completed_tasks"`
}

// ProjectMemberData 项目成员数据传输对象
type ProjectMemberData struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	AddedBy  string    `json:"added_by"`
}
