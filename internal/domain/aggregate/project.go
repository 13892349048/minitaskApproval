package aggregate

import (
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/valueobject"
)

// Project 项目聚合根
type Project struct {
	ID          valueobject.ProjectID
	Name        string
	Description string
	ProjectType valueobject.ProjectType
	Status      valueobject.ProjectStatus

	// 层级关系
	ParentID *valueobject.ProjectID
	Children []valueobject.ProjectID

	// 人员管理
	OwnerID   valueobject.UserID
	ManagerID *valueobject.UserID
	Members   []valueobject.ProjectMember

	// 时间管理
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	// 统计信息
	TaskCount      int
	CompletedTasks int

	// 领域事件
	Events []event.DomainEvent
}

// NewProject 创建新项目
func NewProject(
	id valueobject.ProjectID,
	name, description string,
	projectType valueobject.ProjectType,
	ownerID valueobject.UserID,
) *Project {
	now := time.Now()

	project := &Project{
		ID:          id,
		Name:        name,
		Description: description,
		ProjectType: projectType,
		Status:      valueobject.ProjectStatusDraft,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Events:      make([]event.DomainEvent, 0),
	}

	// 发布项目创建事件
	project.addEvent(event.NewProjectCreatedEvent(id, name, projectType, ownerID))

	return project
}

// UpdateBasicInfo 更新基本信息
func (p *Project) UpdateBasicInfo(name, description string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	oldName := p.Name
	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()

	if oldName != name {
		// 发布项目更新事件
		p.addEvent(event.NewProjectUpdatedEvent(p.ID, oldName, name, valueobject.UserID("")))
	}

	return nil
}

// AssignManager 分配项目管理者
func (p *Project) AssignManager(managerID valueobject.UserID, assignedBy valueobject.UserID) error {
	// 验证权限：只有项目所有者可以分配管理者
	if assignedBy != p.OwnerID {
		return fmt.Errorf("only project owner can assign manager")
	}

	// 管理者不能是所有者本人
	if managerID == p.OwnerID {
		return fmt.Errorf("project owner cannot be manager")
	}

	oldManagerID := p.ManagerID
	p.ManagerID = &managerID
	p.UpdatedAt = time.Now()

	// 如果管理者不在成员列表中，自动添加
	if !p.isMember(managerID) {
		member := valueobject.ProjectMember{
			UserID:   managerID,
			Role:     valueobject.ProjectRoleManager,
			JoinedAt: time.Now(),
			AddedBy:  assignedBy,
		}
		p.Members = append(p.Members, member)
	} else {
		// 更新现有成员的角色
		for i, member := range p.Members {
			if member.UserID == managerID {
				p.Members[i].Role = valueobject.ProjectRoleManager
				break
			}
		}
	}
	// 发布事件
	p.addEvent(event.NewProjectManagerAssignedEvent(p.ID, oldManagerID, &managerID, valueobject.UserID("")))

	return nil
}

// AddMember 添加项目成员
func (p *Project) AddMember(userID valueobject.UserID, role valueobject.ProjectRole, addedBy valueobject.UserID) error {
	// 验证权限：所有者或管理者可以添加成员
	if !p.canManageMembers(addedBy) {
		return fmt.Errorf("insufficient permission to add member")
	}

	// 检查是否已经是成员
	if p.isMember(userID) {
		return fmt.Errorf("user is already a member")
	}

	// 不能添加所有者为普通成员
	if userID == p.OwnerID {
		return fmt.Errorf("project owner cannot be added as member")
	}

	member := valueobject.ProjectMember{
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
		AddedBy:  addedBy,
	}

	p.Members = append(p.Members, member)
	p.UpdatedAt = time.Now()
	// 发布事件
	p.addEvent(event.NewProjectMemberAddedEvent(p.ID, userID, role, valueobject.UserID("")))

	return nil
}

// RemoveMember 移除项目成员
func (p *Project) RemoveMember(userID valueobject.UserID, removedBy valueobject.UserID) error {
	// 验证权限
	if !p.canManageMembers(removedBy) {
		return fmt.Errorf("insufficient permission to remove member")
	}

	// 不能移除所有者
	if userID == p.OwnerID {
		return fmt.Errorf("cannot remove project owner")
	}

	// 不能移除管理者（需要先取消管理者身份）
	if p.ManagerID != nil && userID == *p.ManagerID {
		return fmt.Errorf("cannot remove project manager, unassign manager first")
	}

	// 查找并移除成员
	for i, member := range p.Members {
		if member.UserID == userID {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			p.UpdatedAt = time.Now()

			// 发布事件
			p.addEvent(event.NewProjectMemberRemovedEvent(p.ID, userID, member.Role, valueobject.UserID("")))

			return nil
		}
	}

	return fmt.Errorf("user is not a member of this project")
}

// UpdateMemberRole 更新成员角色
func (p *Project) UpdateMemberRole(userID valueobject.UserID, newRole valueobject.ProjectRole, updatedBy valueobject.UserID) error {
	// 验证权限
	if !p.canManageMembers(updatedBy) {
		return fmt.Errorf("insufficient permission to update member role")
	}

	// 查找并更新成员角色
	for i, member := range p.Members {
		if member.UserID == userID {
			oldRole := member.Role
			p.Members[i].Role = newRole
			p.UpdatedAt = time.Now()

			// 发布成员角色更新事件
			p.addEvent(event.NewProjectMemberRoleUpdatedEvent(
				p.ID,
				userID,
				oldRole,
				newRole,
				updatedBy,
			))

			return nil
		}
	}

	return fmt.Errorf("user is not a member of this project")
}

// CreateSubProject 创建子项目
func (p *Project) CreateSubProject(subProjectID valueobject.ProjectID, name, description string, createdBy valueobject.UserID) (ProjectAggregate, error) {
	// 验证权限
	if !p.canManageProject(createdBy) {
		return nil, fmt.Errorf("insufficient permission to create sub project")
	}

	// 只有主项目可以创建子项目
	if p.ProjectType != valueobject.ProjectTypeMaster {
		return nil, fmt.Errorf("only master project can have sub projects")
	}

	// 项目必须是活跃状态
	if p.Status != valueobject.ProjectStatusActive {
		return nil, fmt.Errorf("parent project must be active to create sub project")
	}

	subProject := NewProject(subProjectID, name, description, valueobject.ProjectTypeSub, createdBy)
	subProject.ParentID = &p.ID

	// 添加到子项目列表
	p.Children = append(p.Children, subProjectID)
	p.UpdatedAt = time.Now()

	// 发布子项目创建事件
	p.addEvent(event.NewSubProjectCreatedEvent(
		p.ID,
		subProjectID,
		name,
		createdBy,
	))

	return subProject, nil
}

// Activate 激活项目
func (p *Project) Activate(activatedBy valueobject.UserID) error {
	if !p.canManageProject(activatedBy) {
		return fmt.Errorf("insufficient permission to activate project")
	}

	if p.Status == valueobject.ProjectStatusActive {
		return fmt.Errorf("project is already active")
	}

	if p.Status == valueobject.ProjectStatusCompleted || p.Status == valueobject.ProjectStatusCancelled {
		return fmt.Errorf("cannot activate completed or cancelled project")
	}

	oldStatus := p.Status
	p.Status = valueobject.ProjectStatusActive
	p.StartDate = time.Now()
	p.UpdatedAt = time.Now()

	p.addEvent(&event.ProjectStatusChangedEvent{
		ProjectID: p.ID,
		OldStatus: oldStatus,
		NewStatus: valueobject.ProjectStatusActive,
		ChangedBy: activatedBy,
	})

	return nil
}

// Pause 暂停项目
func (p *Project) Pause(pausedBy valueobject.UserID, reason string) error {
	if !p.canManageProject(pausedBy) {
		return fmt.Errorf("insufficient permission to pause project")
	}

	if p.Status != valueobject.ProjectStatusActive {
		return fmt.Errorf("only active project can be paused")
	}

	//oldStatus := p.Status
	p.Status = valueobject.ProjectStatusPaused
	p.UpdatedAt = time.Now()
	// 发布事件
	p.addEvent(event.NewProjectStatusChangedEvent(p.ID, p.Status, valueobject.ProjectStatusActive, valueobject.UserID(""), "Project started"))

	return nil
}

// Complete 完成项目
func (p *Project) Complete(completedBy valueobject.UserID) error {
	if !p.canManageProject(completedBy) {
		return fmt.Errorf("insufficient permission to complete project")
	}

	if p.Status == valueobject.ProjectStatusCompleted {
		return fmt.Errorf("project is already completed")
	}

	// 检查是否所有任务都已完成
	if p.TaskCount > 0 && p.CompletedTasks < p.TaskCount {
		return fmt.Errorf("cannot complete project with pending tasks")
	}

	oldStatus := p.Status
	p.Status = valueobject.ProjectStatusCompleted
	now := time.Now()
	p.EndDate = &now
	p.UpdatedAt = now

	p.addEvent(&event.ProjectStatusChangedEvent{
		ProjectID: p.ID,
		OldStatus: oldStatus,
		NewStatus: valueobject.ProjectStatusCompleted,
		ChangedBy: completedBy,
	})

	return nil
}

// Cancel 取消项目
func (p *Project) Cancel(cancelledBy valueobject.UserID, reason string) error {
	if !p.canManageProject(cancelledBy) {
		return fmt.Errorf("insufficient permission to cancel project")
	}

	if p.Status == valueobject.ProjectStatusCompleted || p.Status == valueobject.ProjectStatusCancelled {
		return fmt.Errorf("cannot cancel completed or already cancelled project")
	}

	oldStatus := p.Status
	p.Status = valueobject.ProjectStatusCancelled
	now := time.Now()
	p.EndDate = &now
	p.UpdatedAt = now

	p.addEvent(&event.ProjectStatusChangedEvent{
		ProjectID: p.ID,
		OldStatus: oldStatus,
		NewStatus: valueobject.ProjectStatusCancelled,
		ChangedBy: cancelledBy,
		Reason:    reason,
	})

	return nil
}

// Delete 软删除项目
func (p *Project) Delete(deletedBy valueobject.UserID) error {
	// 只有所有者可以删除项目
	if deletedBy != p.OwnerID {
		return fmt.Errorf("only project owner can delete project")
	}

	// 不能删除有子项目的项目
	if len(p.Children) > 0 {
		return fmt.Errorf("cannot delete project with sub projects")
	}

	// 不能删除有未完成任务的项目
	if p.TaskCount > p.CompletedTasks {
		return fmt.Errorf("cannot delete project with pending tasks")
	}

	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
	// 发布删除事件
	p.addEvent(event.NewProjectDeletedEvent(p.ID, valueobject.UserID("")))

	return nil
}

// UpdateTaskStatistics 更新任务统计
func (p *Project) UpdateTaskStatistics(totalTasks, completedTasks int) {
	p.TaskCount = totalTasks
	p.CompletedTasks = completedTasks
	p.UpdatedAt = time.Now()
}

// CanUserAccess 检查用户是否可以访问项目
func (p *Project) CanUserAccess(userID valueobject.UserID) bool {
	// 所有者和管理者可以访问
	if userID == p.OwnerID {
		return true
	}
	if p.ManagerID != nil && userID == *p.ManagerID {
		return true
	}

	// 成员可以访问
	return p.isMember(userID)
}

// GetMemberRole 获取成员角色
func (p *Project) GetMemberRole(userID valueobject.UserID) *valueobject.ProjectRole {
	if userID == p.OwnerID {
		role := valueobject.ProjectRoleManager
		return &role
	}

	if p.ManagerID != nil && userID == *p.ManagerID {
		role := valueobject.ProjectRoleManager
		return &role
	}

	for _, member := range p.Members {
		if member.UserID == userID {
			return &member.Role
		}
	}

	return nil
}

// GetMemberIDs 获取所有成员ID列表
func (p *Project) GetMemberIDs() []string {
	memberIDs := make([]string, 0, len(p.Members)+2)

	// 添加所有者
	memberIDs = append(memberIDs, string(p.OwnerID))

	// 添加管理者
	if p.ManagerID != nil {
		memberIDs = append(memberIDs, string(*p.ManagerID))
	}

	// 添加普通成员
	for _, member := range p.Members {
		memberIDs = append(memberIDs, string(member.UserID))
	}

	return memberIDs
}

// 私有方法

// isMember 检查是否是项目成员
func (p *Project) isMember(userID valueobject.UserID) bool {
	for _, member := range p.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

// canManageMembers 检查是否可以管理成员
func (p *Project) canManageMembers(userID valueobject.UserID) bool {
	// 所有者可以管理成员
	if userID == p.OwnerID {
		return true
	}
	// 管理者可以管理成员
	if p.ManagerID != nil && userID == *p.ManagerID {
		return true
	}
	return false
}

// canManageProject 检查是否可以管理项目
func (p *Project) canManageProject(userID valueobject.UserID) bool {
	return p.canManageMembers(userID)
}

// addEvent 添加领域事件
func (p *Project) addEvent(event event.DomainEvent) {
	p.Events = append(p.Events, event)
}

// ClearEvents 清空事件
func (p *Project) ClearEvents() {
	p.Events = make([]event.DomainEvent, 0)
}
