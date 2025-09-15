package service

import (
	"context"
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/aggregate"
	authService "github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/service"
	"github.com/taskflow/internal/domain/valueobject"
)

// ProjectAppService 项目应用服务
type ProjectAppService struct {
	projectDomainService service.ProjectDomainService
	transactionMgr       authService.TransactionManager
	projectRepo          repository.ProjectRepository
}

// NewProjectAppService 创建项目应用服务
func NewProjectAppService(
	projectDomainService service.ProjectDomainService,
	transactionMgr authService.TransactionManager,
	projectRepo repository.ProjectRepository,
) *ProjectAppService {
	return &ProjectAppService{
		projectDomainService: projectDomainService,
		transactionMgr:       transactionMgr,
		projectRepo:          projectRepo,
	}
}

// CreateProject 创建项目（需要事务）
func (s *ProjectAppService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectResponse, error) {
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 创建项目聚合
		project := aggregate.NewProject(
			valueobject.ProjectID(req.ID),
			req.Name,
			req.Description,
			valueobject.ProjectType(req.ProjectType),
			valueobject.UserID(req.OwnerID),
		)

		// 2. 保存项目
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return nil, fmt.Errorf("保存项目失败: %w", err)
		}

		// 3. 返回结果
		return &ProjectResponse{
			ID:          string(project.ID),
			Name:        project.Name,
			Description: project.Description,
			ProjectType: string(project.ProjectType),
			Status:      string(project.Status),
			OwnerID:     string(project.OwnerID),
			StartDate:   project.StartDate,
			EndDate:     project.EndDate,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	if projectResponse, ok := result.(*ProjectResponse); ok {
		return projectResponse, nil
	}

	return nil, fmt.Errorf("unexpected result type")
}

// GetProject 获取项目（不需要事务）
func (s *ProjectAppService) GetProject(ctx context.Context, id string) (*ProjectResponse, error) {
	project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(id))
	if err != nil {
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}

	var managerID *string
	if project.ManagerID != nil {
		managerIDStr := string(*project.ManagerID)
		managerID = &managerIDStr
	}

	return &ProjectResponse{
		ID:          string(project.ID),
		Name:        project.Name,
		Description: project.Description,
		ProjectType: string(project.ProjectType),
		Status:      string(project.Status),
		OwnerID:     string(project.OwnerID),
		ManagerID:   managerID,
		StartDate:   project.StartDate,
		EndDate:     project.EndDate,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}, nil
}

// UpdateProject 更新项目（需要事务）
func (s *ProjectAppService) UpdateProject(ctx context.Context, req *UpdateProjectRequest) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(req.ID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 2. 更新项目信息
		if err := project.UpdateBasicInfo(req.Name, req.Description); err != nil {
			return fmt.Errorf("更新项目信息失败: %w", err)
		}

		// 3. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// AssignManager 分配项目管理者（需要事务）
func (s *ProjectAppService) AssignManager(ctx context.Context, projectID, managerID, assignedBy string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 2. 分配管理者
		if err := project.AssignManager(
			valueobject.UserID(managerID),
			valueobject.UserID(assignedBy),
		); err != nil {
			return fmt.Errorf("分配管理者失败: %w", err)
		}

		// 3. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// AddMember 添加项目成员（需要事务）
func (s *ProjectAppService) AddMember(ctx context.Context, projectID, userID, addedBy string, role string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 领域服务验证
		if err := s.projectDomainService.ValidateMemberAddition(
			ctx,
			valueobject.ProjectID(projectID),
			valueobject.UserID(userID),
			valueobject.ProjectRole(role),
		); err != nil {
			return fmt.Errorf("成员添加验证失败: %w", err)
		}

		// 2. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 3. 添加成员
		if err := project.AddMember(
			valueobject.UserID(userID),
			valueobject.ProjectRole(role),
			valueobject.UserID(addedBy),
		); err != nil {
			return fmt.Errorf("添加成员失败: %w", err)
		}

		// 4. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// RemoveMember 移除项目成员（需要事务）
func (s *ProjectAppService) RemoveMember(ctx context.Context, projectID, userID, removedBy string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 2. 移除成员
		if err := project.RemoveMember(
			valueobject.UserID(userID),
			valueobject.UserID(removedBy),
		); err != nil {
			return fmt.Errorf("移除成员失败: %w", err)
		}

		// 3. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// UpdateMemberRole 更新成员角色（需要事务）
func (s *ProjectAppService) UpdateMemberRole(ctx context.Context, projectID, userID, updatedBy string, newRole string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 2. 更新成员角色
		if err := project.UpdateMemberRole(
			valueobject.UserID(userID),
			valueobject.ProjectRole(newRole),
			valueobject.UserID(updatedBy),
		); err != nil {
			return fmt.Errorf("更新成员角色失败: %w", err)
		}

		// 3. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// ChangeStatus 更改项目状态（需要事务）
func (s *ProjectAppService) ChangeStatus(ctx context.Context, projectID, userID string, newStatus string, reason string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 验证状态更改权限
		canChange, err := s.projectDomainService.CanChangeProjectStatus(
			ctx,
			valueobject.ProjectID(projectID),
			valueobject.ProjectStatus(newStatus),
			valueobject.UserID(userID),
		)
		if err != nil {
			return fmt.Errorf("状态更改验证失败: %w", err)
		}
		if !canChange {
			return fmt.Errorf("用户无权限更改项目状态")
		}

		// 2. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 3. 根据状态执行相应操作
		switch valueobject.ProjectStatus(newStatus) {
		case valueobject.ProjectStatusActive:
			if err := project.Activate(valueobject.UserID(userID)); err != nil {
				return fmt.Errorf("激活项目失败: %w", err)
			}
		case valueobject.ProjectStatusPaused:
			if err := project.Pause(valueobject.UserID(userID), reason); err != nil {
				return fmt.Errorf("暂停项目失败: %w", err)
			}
		case valueobject.ProjectStatusCompleted:
			// 验证完成条件
			if err := s.projectDomainService.ValidateProjectCompletion(ctx, valueobject.ProjectID(projectID)); err != nil {
				return fmt.Errorf("项目完成条件验证失败: %w", err)
			}
			if err := project.Complete(valueobject.UserID(userID)); err != nil {
				return fmt.Errorf("完成项目失败: %w", err)
			}
		case valueobject.ProjectStatusCancelled:
			if err := project.Cancel(valueobject.UserID(userID), reason); err != nil {
				return fmt.Errorf("取消项目失败: %w", err)
			}
		default:
			return fmt.Errorf("不支持的项目状态: %s", newStatus)
		}

		// 4. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// CreateSubProject 创建子项目（需要事务）
func (s *ProjectAppService) CreateSubProject(ctx context.Context, parentID, name, description, createdBy string) (*ProjectResponse, error) {
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 验证是否可以创建子项目
		canCreate, err := s.projectDomainService.CanCreateSubProject(
			ctx,
			valueobject.ProjectID(parentID),
			valueobject.UserID(createdBy),
		)
		if err != nil {
			return nil, fmt.Errorf("子项目创建验证失败: %w", err)
		}
		if !canCreate {
			return nil, fmt.Errorf("用户无权限创建子项目")
		}

		// 2. 查找父项目
		parentProject, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(parentID))
		if err != nil {
			return nil, fmt.Errorf("父项目不存在: %w", err)
		}

		// 3. 创建子项目
		subProjectID := generateProjectID()
		subProject, err := parentProject.CreateSubProject(
			valueobject.ProjectID(subProjectID),
			name,
			description,
			valueobject.UserID(createdBy),
		)
		if err != nil {
			return nil, fmt.Errorf("创建子项目失败: %w", err)
		}

		// 4. 保存父项目和子项目
		if err := s.projectRepo.Save(ctx, *parentProject); err != nil {
			return nil, fmt.Errorf("保存父项目失败: %w", err)
		}

		if concreteSubProject, ok := subProject.(*aggregate.Project); ok {
			if err := s.projectRepo.Save(ctx, *concreteSubProject); err != nil {
				return nil, fmt.Errorf("保存子项目失败: %w", err)
			}

			// 5. 返回响应
			return s.buildProjectResponse(*concreteSubProject), nil
		}

		return nil, fmt.Errorf("子项目类型转换失败")
	})

	if err != nil {
		return nil, err
	}

	if projectResponse, ok := result.(*ProjectResponse); ok {
		return projectResponse, nil
	}

	return nil, fmt.Errorf("unexpected result type")
}

// ListProjects 获取项目列表（不需要事务）
func (s *ProjectAppService) ListProjects(ctx context.Context, req *ProjectListRequest) (*ProjectListResponse, error) {
	// 构建查询条件
	criteria := aggregate.ProjectSearchCriteria{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	// 设置状态过滤
	if req.Status != "" {
		status := valueobject.ProjectStatus(req.Status)
		criteria.Status = &status
	}

	// 设置项目类型过滤
	if req.Type != "" {
		projectType := valueobject.ProjectType(req.Type)
		criteria.ProjectType = &projectType
	}

	// 设置所有者过滤
	if req.OwnerID != "" {
		ownerID := valueobject.UserID(req.OwnerID)
		criteria.OwnerID = &ownerID
	}

	// 设置管理者过滤
	if req.ManagerID != "" {
		managerID := valueobject.UserID(req.ManagerID)
		criteria.ManagerID = &managerID
	}

	// 设置名称搜索
	if req.Search != "" {
		criteria.Name = &req.Search
	}

	// 设置排序
	if req.SortBy != "" {
		criteria.OrderBy = req.SortBy
	}
	if req.SortOrder != "" {
		criteria.OrderDir = req.SortOrder
	}

	// 查询项目
	projects, total, err := s.projectRepo.SearchProjects(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("查询项目失败: %w", err)
	}

	// 构建响应
	projectResponses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = *s.buildProjectResponse(project)
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	return &ProjectListResponse{
		Projects:   projectResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetProjectHierarchy 获取项目层级结构（不需要事务）
func (s *ProjectAppService) GetProjectHierarchy(ctx context.Context, projectID string) (*ProjectHierarchyResponse, error) {
	hierarchy, err := s.projectDomainService.GetProjectHierarchy(ctx, valueobject.ProjectID(projectID))
	if err != nil {
		return nil, fmt.Errorf("获取项目层级失败: %w", err)
	}

	response := &ProjectHierarchyResponse{
		Project:       s.buildProjectResponse(*hierarchy.Project),
		Depth:         hierarchy.Depth,
		TotalProjects: hierarchy.TotalProjects,
	}

	if hierarchy.Parent != nil {
		response.Parent = s.buildProjectResponse(*hierarchy.Parent)
	}

	if len(hierarchy.Children) > 0 {
		response.Children = make([]ProjectResponse, len(hierarchy.Children))
		for i, child := range hierarchy.Children {
			response.Children[i] = *s.buildProjectResponse(child)
		}
	}

	return response, nil
}

// DeleteProject 删除项目（需要事务）
func (s *ProjectAppService) DeleteProject(ctx context.Context, projectID, deletedBy string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找项目
		project, err := s.projectRepo.FindByID(ctx, valueobject.ProjectID(projectID))
		if err != nil {
			return fmt.Errorf("项目不存在: %w", err)
		}

		// 2. 执行软删除
		if err := project.Delete(valueobject.UserID(deletedBy)); err != nil {
			return fmt.Errorf("删除项目失败: %w", err)
		}

		// 3. 保存更新
		if err := s.projectRepo.Save(ctx, *project); err != nil {
			return fmt.Errorf("保存项目失败: %w", err)
		}

		return nil
	})
}

// 辅助方法

// buildProjectResponse 构建项目响应
func (s *ProjectAppService) buildProjectResponse(project aggregate.Project) *ProjectResponse {
	// 转换成员列表
	members := make([]ProjectMemberResponse, len(project.Members))
	for i, member := range project.Members {
		members[i] = ToProjectMemberResponse(member)
	}

	// 转换子项目ID列表
	children := make([]string, len(project.Children))
	for i, childID := range project.Children {
		children[i] = string(childID)
	}

	response := &ProjectResponse{
		ID:          string(project.ID),
		Name:        project.Name,
		Description: project.Description,
		ProjectType: string(project.ProjectType),
		Status:      string(project.Status),
		OwnerID:     string(project.OwnerID),
		Members:     members,
		Children:    children,
		StartDate:   project.StartDate,
		EndDate:     project.EndDate,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	// 设置管理者ID
	if project.ManagerID != nil {
		managerID := string(*project.ManagerID)
		response.ManagerID = &managerID
	}

	// 设置父项目ID
	if project.ParentID != nil {
		parentID := string(*project.ParentID)
		response.ParentID = &parentID
	}

	// 设置统计信息
	response.Statistics = ToProjectStatisticsResponse(
		project.TaskCount,
		project.CompletedTasks,
		len(project.Members)+1, // +1 for owner
	)

	return response
}

// generateProjectID 生成项目ID
func generateProjectID() string {
	// 这里可以使用UUID或其他ID生成策略
	return "proj_" + fmt.Sprintf("%d", time.Now().UnixNano())
}
