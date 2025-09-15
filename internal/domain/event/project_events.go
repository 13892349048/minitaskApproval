package event

import "github.com/taskflow/internal/domain/valueobject"

// ProjectCreatedEvent 项目创建事件
type ProjectCreatedEvent struct {
	*BaseEvent
	ProjectID   valueobject.ProjectID   `json:"project_id"`
	Name        string                  `json:"name"`
	ProjectType valueobject.ProjectType `json:"project_type"`
	OwnerID     valueobject.UserID      `json:"owner_id"`
}

// NewProjectCreatedEvent 创建项目创建事件
func NewProjectCreatedEvent(projectID valueobject.ProjectID, name string, projectType valueobject.ProjectType, ownerID valueobject.UserID) *ProjectCreatedEvent {
	return &ProjectCreatedEvent{
		BaseEvent:   NewBaseEvent("project.created", string(projectID), "project"),
		ProjectID:   projectID,
		Name:        name,
		ProjectType: projectType,
		OwnerID:     ownerID,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectCreatedEvent) EventData() interface{} {
	return e
}

// ProjectUpdatedEvent 项目更新事件
type ProjectUpdatedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID `json:"project_id"`
	OldName   string                `json:"old_name"`
	NewName   string                `json:"new_name"`
	UpdatedBy valueobject.UserID    `json:"updated_by"`
}

// NewProjectUpdatedEvent 创建项目更新事件
func NewProjectUpdatedEvent(projectID valueobject.ProjectID, oldName, newName string, updatedBy valueobject.UserID) *ProjectUpdatedEvent {
	return &ProjectUpdatedEvent{
		BaseEvent: NewBaseEvent("project.updated", string(projectID), "project"),
		ProjectID: projectID,
		OldName:   oldName,
		NewName:   newName,
		UpdatedBy: updatedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectUpdatedEvent) EventData() interface{} {
	return e
}

// ProjectManagerAssignedEvent 项目管理者分配事件
type ProjectManagerAssignedEvent struct {
	*BaseEvent
	ProjectID    valueobject.ProjectID `json:"project_id"`
	OldManagerID *valueobject.UserID   `json:"old_manager_id"`
	NewManagerID *valueobject.UserID   `json:"new_manager_id"`
	AssignedBy   valueobject.UserID    `json:"assigned_by"`
}

// NewProjectManagerAssignedEvent 创建项目管理者分配事件
func NewProjectManagerAssignedEvent(projectID valueobject.ProjectID, oldManagerID, newManagerID *valueobject.UserID, assignedBy valueobject.UserID) *ProjectManagerAssignedEvent {
	return &ProjectManagerAssignedEvent{
		BaseEvent:    NewBaseEvent("project.manager_assigned", string(projectID), "project"),
		ProjectID:    projectID,
		OldManagerID: oldManagerID,
		NewManagerID: newManagerID,
		AssignedBy:   assignedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectManagerAssignedEvent) EventData() interface{} {
	return e
}

// ProjectMemberAddedEvent 项目成员添加事件
type ProjectMemberAddedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID   `json:"project_id"`
	UserID    valueobject.UserID      `json:"user_id"`
	Role      valueobject.ProjectRole `json:"role"`
	AddedBy   valueobject.UserID      `json:"added_by"`
}

// NewProjectMemberAddedEvent 创建项目成员添加事件
func NewProjectMemberAddedEvent(projectID valueobject.ProjectID, userID valueobject.UserID, role valueobject.ProjectRole, addedBy valueobject.UserID) *ProjectMemberAddedEvent {
	return &ProjectMemberAddedEvent{
		BaseEvent: NewBaseEvent("project.member_added", string(projectID), "project"),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		AddedBy:   addedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectMemberAddedEvent) EventData() interface{} {
	return e
}

// ProjectMemberRemovedEvent 项目成员移除事件
type ProjectMemberRemovedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID   `json:"project_id"`
	UserID    valueobject.UserID      `json:"user_id"`
	Role      valueobject.ProjectRole `json:"role"`
	RemovedBy valueobject.UserID      `json:"removed_by"`
}

// NewProjectMemberRemovedEvent 创建项目成员移除事件
func NewProjectMemberRemovedEvent(projectID valueobject.ProjectID, userID valueobject.UserID, role valueobject.ProjectRole, removedBy valueobject.UserID) *ProjectMemberRemovedEvent {
	return &ProjectMemberRemovedEvent{
		BaseEvent: NewBaseEvent("project.member_removed", string(projectID), "project"),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		RemovedBy: removedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectMemberRemovedEvent) EventData() interface{} {
	return e
}

// ProjectStatusChangedEvent 项目状态变更事件
type ProjectStatusChangedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID     `json:"project_id"`
	OldStatus valueobject.ProjectStatus `json:"old_status"`
	NewStatus valueobject.ProjectStatus `json:"new_status"`
	ChangedBy valueobject.UserID        `json:"changed_by"`
	Reason    string                    `json:"reason,omitempty"`
}

// NewProjectStatusChangedEvent 创建项目状态变更事件
func NewProjectStatusChangedEvent(projectID valueobject.ProjectID, oldStatus, newStatus valueobject.ProjectStatus, changedBy valueobject.UserID, reason string) *ProjectStatusChangedEvent {
	return &ProjectStatusChangedEvent{
		BaseEvent: NewBaseEvent("project.status_changed", string(projectID), "project"),
		ProjectID: projectID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		ChangedBy: changedBy,
		Reason:    reason,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectStatusChangedEvent) EventData() interface{} {
	return e
}

// ProjectDeletedEvent 项目删除事件
type ProjectDeletedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID `json:"project_id"`
	DeletedBy valueobject.UserID    `json:"deleted_by"`
}

// NewProjectDeletedEvent 创建项目删除事件
func NewProjectDeletedEvent(projectID valueobject.ProjectID, deletedBy valueobject.UserID) *ProjectDeletedEvent {
	return &ProjectDeletedEvent{
		BaseEvent: NewBaseEvent("project.deleted", string(projectID), "project"),
		ProjectID: projectID,
		DeletedBy: deletedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectDeletedEvent) EventData() interface{} {
	return e
}

// ProjectMemberRoleUpdatedEvent 项目成员角色更新事件
type ProjectMemberRoleUpdatedEvent struct {
	*BaseEvent
	ProjectID valueobject.ProjectID   `json:"project_id"`
	UserID    valueobject.UserID      `json:"user_id"`
	OldRole   valueobject.ProjectRole `json:"old_role"`
	NewRole   valueobject.ProjectRole `json:"new_role"`
	UpdatedBy valueobject.UserID      `json:"updated_by"`
}

// NewProjectMemberRoleUpdatedEvent 创建项目成员角色更新事件
func NewProjectMemberRoleUpdatedEvent(projectID valueobject.ProjectID, userID valueobject.UserID, oldRole, newRole valueobject.ProjectRole, updatedBy valueobject.UserID) *ProjectMemberRoleUpdatedEvent {
	return &ProjectMemberRoleUpdatedEvent{
		BaseEvent: NewBaseEvent("project.member_role_updated", string(projectID), "project"),
		ProjectID: projectID,
		UserID:    userID,
		OldRole:   oldRole,
		NewRole:   newRole,
		UpdatedBy: updatedBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *ProjectMemberRoleUpdatedEvent) EventData() interface{} {
	return e
}

// SubProjectCreatedEvent 子项目创建事件
type SubProjectCreatedEvent struct {
	*BaseEvent
	ParentProjectID valueobject.ProjectID `json:"parent_project_id"`
	SubProjectID    valueobject.ProjectID `json:"sub_project_id"`
	Name            string                `json:"name"`
	CreatedBy       valueobject.UserID    `json:"created_by"`
}

// NewSubProjectCreatedEvent 创建子项目创建事件
func NewSubProjectCreatedEvent(parentProjectID, subProjectID valueobject.ProjectID, name string, createdBy valueobject.UserID) *SubProjectCreatedEvent {
	return &SubProjectCreatedEvent{
		BaseEvent:       NewBaseEvent("project.sub_project_created", string(parentProjectID), "project"),
		ParentProjectID: parentProjectID,
		SubProjectID:    subProjectID,
		Name:            name,
		CreatedBy:       createdBy,
	}
}

// EventData 实现 DomainEvent 接口
func (e *SubProjectCreatedEvent) EventData() interface{} {
	return e
}

// 确保所有事件都实现了 DomainEvent 接口
var _ DomainEvent = (*ProjectCreatedEvent)(nil)
var _ DomainEvent = (*ProjectUpdatedEvent)(nil)
var _ DomainEvent = (*ProjectManagerAssignedEvent)(nil)
var _ DomainEvent = (*ProjectMemberAddedEvent)(nil)
var _ DomainEvent = (*ProjectMemberRemovedEvent)(nil)
var _ DomainEvent = (*ProjectStatusChangedEvent)(nil)
var _ DomainEvent = (*ProjectDeletedEvent)(nil)
var _ DomainEvent = (*ProjectMemberRoleUpdatedEvent)(nil)
var _ DomainEvent = (*SubProjectCreatedEvent)(nil)
