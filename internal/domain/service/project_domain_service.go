package service

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
)

// ProjectDomainService 项目领域服务接口
type ProjectDomainService interface {
	// 项目层级管理
	CanCreateSubProject(ctx context.Context, parentProjectID valueobject.ProjectID, userID valueobject.UserID) (bool, error)
	ValidateProjectHierarchy(ctx context.Context, parentID, childID valueobject.ProjectID) error
	GetProjectHierarchy(ctx context.Context, projectID valueobject.ProjectID) (*ProjectHierarchy, error)

	// 项目权限验证
	CanUserAccessProject(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (bool, error)
	CanUserManageProject(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (bool, error)
	GetUserProjectRole(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (*valueobject.ProjectRole, error)

	// 项目成员管理
	ValidateMemberAddition(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID, role valueobject.ProjectRole) error
	GetProjectMemberStatistics(ctx context.Context, projectID valueobject.ProjectID) (*ProjectMemberStats, error)

	// 项目状态管理
	CanChangeProjectStatus(ctx context.Context, projectID valueobject.ProjectID, newStatus valueobject.ProjectStatus, userID valueobject.UserID) (bool, error)
	ValidateProjectCompletion(ctx context.Context, projectID valueobject.ProjectID) error
}

// ProjectHierarchy 项目层级结构
type ProjectHierarchy struct {
	Project       *aggregate.Project  `json:"project"`
	Parent        *aggregate.Project  `json:"parent,omitempty"`
	Children      []aggregate.Project `json:"children"`
	Depth         int                 `json:"depth"`
	TotalProjects int                 `json:"total_projects"`
}

// ProjectMemberStats 项目成员统计
type ProjectMemberStats struct {
	TotalMembers     int                             `json:"total_members"`
	RoleDistribution map[valueobject.ProjectRole]int `json:"role_distribution"`
	ActiveMembers    int                             `json:"active_members"`
	RecentJoins      []valueobject.ProjectMember     `json:"recent_joins"`
}

// ProjectDomainServiceImpl 项目领域服务实现
type ProjectDomainServiceImpl struct {
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

// NewProjectDomainService 创建项目领域服务
func NewProjectDomainService(
	projectRepo repository.ProjectRepository,
	userRepo repository.UserRepository,
) *ProjectDomainServiceImpl {
	return &ProjectDomainServiceImpl{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CanCreateSubProject 检查是否可以创建子项目
func (s *ProjectDomainServiceImpl) CanCreateSubProject(ctx context.Context, parentProjectID valueobject.ProjectID, userID valueobject.UserID) (bool, error) {
	// 1. 获取父项目
	parentProject, err := s.projectRepo.FindByID(ctx, parentProjectID)
	if err != nil {
		return false, fmt.Errorf("failed to find parent project: %w", err)
	}

	// 2. 检查项目类型
	if parentProject.ProjectType != valueobject.ProjectTypeMaster {
		return false, fmt.Errorf("only master projects can have sub projects")
	}

	// 3. 检查项目状态
	if parentProject.Status != valueobject.ProjectStatusActive {
		return false, fmt.Errorf("parent project must be active")
	}

	// 4. 检查用户权限
	if !parentProject.CanUserAccess(userID) {
		return false, fmt.Errorf("user does not have access to parent project")
	}

	// 5. 检查管理权限
	canManage := parentProject.OwnerID == userID ||
		(parentProject.ManagerID != nil && *parentProject.ManagerID == userID)

	return canManage, nil
}

// ValidateProjectHierarchy 验证项目层级关系
func (s *ProjectDomainServiceImpl) ValidateProjectHierarchy(ctx context.Context, parentID, childID valueobject.ProjectID) error {
	// 1. 检查循环引用
	if parentID == childID {
		return fmt.Errorf("project cannot be parent of itself")
	}

	// 2. 检查父项目存在
	parentProject, err := s.projectRepo.FindByID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("parent project not found: %w", err)
	}

	// 3. 检查子项目存在
	childProject, err := s.projectRepo.FindByID(ctx, childID)
	if err != nil {
		return fmt.Errorf("child project not found: %w", err)
	}

	// 4. 检查层级限制（最多2层）
	if parentProject.ParentID != nil {
		return fmt.Errorf("project hierarchy cannot exceed 2 levels")
	}

	// 5. 检查子项目类型
	if childProject.ProjectType != valueobject.ProjectTypeSub {
		return fmt.Errorf("child project must be of type 'sub'")
	}

	return nil
}

// GetProjectHierarchy 获取项目层级结构
func (s *ProjectDomainServiceImpl) GetProjectHierarchy(ctx context.Context, projectID valueobject.ProjectID) (*ProjectHierarchy, error) {
	// 1. 获取当前项目
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	hierarchy := &ProjectHierarchy{
		Project:       project,
		Depth:         0,
		TotalProjects: 1,
	}

	// 2. 获取父项目
	if project.ParentID != nil {
		parent, err := s.projectRepo.FindByID(ctx, *project.ParentID)
		if err == nil {
			hierarchy.Parent = parent
			hierarchy.Depth = 1
			hierarchy.TotalProjects++
		}
	}

	// 3. 获取子项目
	if len(project.Children) > 0 {
		children, err := s.projectRepo.FindByIDs(ctx, project.Children)
		if err == nil {
			hierarchy.Children = children
			hierarchy.TotalProjects += len(children)
		}
	}

	return hierarchy, nil
}

// CanUserAccessProject 检查用户是否可以访问项目
func (s *ProjectDomainServiceImpl) CanUserAccessProject(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (bool, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to find project: %w", err)
	}

	return project.CanUserAccess(userID), nil
}

// CanUserManageProject 检查用户是否可以管理项目
func (s *ProjectDomainServiceImpl) CanUserManageProject(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (bool, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to find project: %w", err)
	}

	// 所有者和管理者可以管理项目
	canManage := project.OwnerID == userID ||
		(project.ManagerID != nil && *project.ManagerID == userID)

	return canManage, nil
}

// GetUserProjectRole 获取用户在项目中的角色
func (s *ProjectDomainServiceImpl) GetUserProjectRole(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID) (*valueobject.ProjectRole, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	return project.GetMemberRole(userID), nil
}

// ValidateMemberAddition 验证成员添加
func (s *ProjectDomainServiceImpl) ValidateMemberAddition(ctx context.Context, projectID valueobject.ProjectID, userID valueobject.UserID, role valueobject.ProjectRole) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.FindByID(ctx, string(userID))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 2. 检查用户状态
	if !user.IsActive() {
		return fmt.Errorf("cannot add inactive user to project")
	}

	// 3. 检查项目
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// 4. 检查项目状态
	if project.Status == valueobject.ProjectStatusCompleted || project.Status == valueobject.ProjectStatusCancelled {
		return fmt.Errorf("cannot add members to completed or cancelled project")
	}

	// 5. 检查角色有效性
	validRoles := []valueobject.ProjectRole{
		valueobject.ProjectRoleMember,
		valueobject.ProjectRoleDeveloper,
		valueobject.ProjectRoleTester,
	}

	roleValid := false
	for _, validRole := range validRoles {
		if role == validRole {
			roleValid = true
			break
		}
	}

	if !roleValid {
		return fmt.Errorf("invalid project role: %s", role)
	}

	return nil
}

// GetProjectMemberStatistics 获取项目成员统计
func (s *ProjectDomainServiceImpl) GetProjectMemberStatistics(ctx context.Context, projectID valueobject.ProjectID) (*ProjectMemberStats, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	stats := &ProjectMemberStats{
		TotalMembers:     len(project.Members) + 1, // +1 for owner
		RoleDistribution: make(map[valueobject.ProjectRole]int),
		ActiveMembers:    0,
		RecentJoins:      make([]valueobject.ProjectMember, 0),
	}

	// 统计角色分布
	for _, member := range project.Members {
		stats.RoleDistribution[member.Role]++
		stats.ActiveMembers++

		// 最近7天加入的成员
		// if time.Since(member.JoinedAt).Hours() < 168 { // 7 days
		// 	stats.RecentJoins = append(stats.RecentJoins, member)
		// }
	}

	return stats, nil
}

// CanChangeProjectStatus 检查是否可以更改项目状态
func (s *ProjectDomainServiceImpl) CanChangeProjectStatus(ctx context.Context, projectID valueobject.ProjectID, newStatus valueobject.ProjectStatus, userID valueobject.UserID) (bool, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to find project: %w", err)
	}

	// 只有所有者和管理者可以更改状态
	canManage := project.OwnerID == userID ||
		(project.ManagerID != nil && *project.ManagerID == userID)

	if !canManage {
		return false, nil
	}

	// 检查状态转换是否有效
	return s.isValidStatusTransition(project.Status, newStatus), nil
}

// ValidateProjectCompletion 验证项目完成条件
func (s *ProjectDomainServiceImpl) ValidateProjectCompletion(ctx context.Context, projectID valueobject.ProjectID) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	// 检查是否有未完成的任务
	if project.TaskCount > 0 && project.CompletedTasks < project.TaskCount {
		return fmt.Errorf("project has %d pending tasks", project.TaskCount-project.CompletedTasks)
	}

	// 检查子项目状态
	if len(project.Children) > 0 {
		children, err := s.projectRepo.FindByIDs(ctx, project.Children)
		if err != nil {
			return fmt.Errorf("failed to check child projects: %w", err)
		}

		for _, child := range children {
			if child.Status != valueobject.ProjectStatusCompleted && child.Status != valueobject.ProjectStatusCancelled {
				return fmt.Errorf("child project %s is not completed", child.Name)
			}
		}
	}

	return nil
}

// isValidStatusTransition 检查状态转换是否有效
func (s *ProjectDomainServiceImpl) isValidStatusTransition(currentStatus, newStatus valueobject.ProjectStatus) bool {
	// 定义有效的状态转换
	validTransitions := map[valueobject.ProjectStatus][]valueobject.ProjectStatus{
		valueobject.ProjectStatusDraft: {
			valueobject.ProjectStatusActive,
			valueobject.ProjectStatusCancelled,
		},
		valueobject.ProjectStatusActive: {
			valueobject.ProjectStatusPaused,
			valueobject.ProjectStatusCompleted,
			valueobject.ProjectStatusCancelled,
		},
		valueobject.ProjectStatusPaused: {
			valueobject.ProjectStatusActive,
			valueobject.ProjectStatusCancelled,
		},
		// 完成和取消状态不能转换到其他状态
		valueobject.ProjectStatusCompleted: {},
		valueobject.ProjectStatusCancelled: {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowedStatuses {
		if newStatus == allowedStatus {
			return true
		}
	}

	return false
}
