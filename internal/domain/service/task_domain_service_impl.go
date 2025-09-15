package service

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
)

// TaskDomainServiceImpl 任务领域服务实现
type TaskDomainServiceImpl struct {
	taskRepo    repository.TaskRepository
	userRepo    repository.UserRepository
	projectRepo repository.ProjectRepository
}

// NewTaskDomainService 创建任务领域服务
func NewTaskDomainService(
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	projectRepo repository.ProjectRepository,
) TaskDomainService {
	return &TaskDomainServiceImpl{
		taskRepo:    taskRepo,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

// ValidateTaskCreation 验证任务创建
func (s *TaskDomainServiceImpl) ValidateTaskCreation(task aggregate.TaskAggregate, createdBy valueobject.UserID) error {
	// 1. 验证创建者权限
	if !s.CanUserCreateTaskInProject(createdBy, task.ProjectID) {
		return fmt.Errorf("user does not have permission to create task in project")
	}

	// 2. 验证负责人存在
	if task.ResponsibleID != "" {
		_, err := s.userRepo.FindByID(context.Background(), string(task.ResponsibleID))
		if err != nil {
			return fmt.Errorf("responsible user not found: %w", err)
		}
	}

	// 3. 验证项目存在且活跃
	project, err := s.projectRepo.FindByID(context.Background(), task.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}
	if project.Status != valueobject.ProjectStatusActive {
		return fmt.Errorf("cannot create task in inactive project")
	}

	return nil
}

// ValidateTaskAssignment 验证任务分配
func (s *TaskDomainServiceImpl) ValidateTaskAssignment(task aggregate.TaskAggregate, responsibleID valueobject.UserID, assignedBy valueobject.UserID) error {
	// 1. 验证分配者权限
	if !s.CanUserManageTask(assignedBy, task) {
		return fmt.Errorf("user does not have permission to assign task")
	}

	// 2. 验证被分配者存在
	_, err := s.userRepo.FindByID(context.Background(), string(responsibleID))
	if err != nil {
		return fmt.Errorf("assignee user not found: %w", err)
	}

	// 3. 验证任务状态允许分配
	if task.Status == valueobject.TaskStatusCompleted || task.Status == valueobject.TaskStatusCancelled {
		return fmt.Errorf("cannot assign completed or cancelled task")
	}

	return nil
}

// ValidateParticipantAddition 验证参与者添加
func (s *TaskDomainServiceImpl) ValidateParticipantAddition(task aggregate.TaskAggregate, participantID valueobject.UserID, addedBy valueobject.UserID) error {
	// 1. 验证添加者权限
	if !s.CanUserManageTask(addedBy, task) {
		return fmt.Errorf("user does not have permission to add participant")
	}

	// 2. 验证参与者存在
	_, err := s.userRepo.FindByID(context.Background(), string(participantID))
	if err != nil {
		return fmt.Errorf("participant user not found: %w", err)
	}

	// 3. 验证不重复添加
	if task.IsParticipant(participantID) {
		return fmt.Errorf("user is already a participant")
	}

	return nil
}

// ValidateStatusTransition 验证状态转换
func (s *TaskDomainServiceImpl) ValidateStatusTransition(task aggregate.TaskAggregate, fromStatus, toStatus valueobject.TaskStatus, changedBy valueobject.UserID) error {
	// 定义允许的状态转换
	allowedTransitions := map[valueobject.TaskStatus][]valueobject.TaskStatus{
		valueobject.TaskStatusDraft:           {valueobject.TaskStatusPendingApproval, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusPendingApproval: {valueobject.TaskStatusApproved, valueobject.TaskStatusRejected, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusApproved:        {valueobject.TaskStatusInProgress, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusRejected:        {valueobject.TaskStatusDraft, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusInProgress:      {valueobject.TaskStatusPaused, valueobject.TaskStatusCompleted, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusPaused:          {valueobject.TaskStatusInProgress, valueobject.TaskStatusCancelled},
		valueobject.TaskStatusCompleted:       {}, // 完成状态不允许转换
		valueobject.TaskStatusCancelled:       {}, // 取消状态不允许转换
	}

	// 检查转换是否允许
	allowed, exists := allowedTransitions[fromStatus]
	if !exists {
		return fmt.Errorf("invalid from status: %s", fromStatus)
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == toStatus {
			return nil // 转换允许
		}
	}

	return fmt.Errorf("status transition from %s to %s is not allowed", fromStatus, toStatus)
}

// ValidateTaskCompletion 验证任务完成
func (s *TaskDomainServiceImpl) ValidateTaskCompletion(task aggregate.TaskAggregate, completedBy valueobject.UserID) error {
	// 1. 验证完成者权限
	if !task.CanUserExecute(completedBy) {
		return fmt.Errorf("user does not have permission to complete task")
	}

	// 2. 验证任务状态
	if task.Status != valueobject.TaskStatusInProgress {
		return fmt.Errorf("only in-progress tasks can be completed")
	}

	// 3. 验证所有必要工作已提交（简化实现）
	// 实际实现中应该检查工作提交和审核状态

	return nil
}

// TransferTaskResponsibility 转移任务责任
func (s *TaskDomainServiceImpl) TransferTaskResponsibility(ctx context.Context, task aggregate.TaskAggregate, newResponsibleID valueobject.UserID, transferredBy valueobject.UserID) error {
	// 1. 验证转移者权限
	if !s.CanUserManageTask(transferredBy, task) {
		return fmt.Errorf("user does not have permission to transfer task responsibility")
	}

	// 2. 验证新负责人存在
	_, err := s.userRepo.FindByID(ctx, string(newResponsibleID))
	if err != nil {
		return fmt.Errorf("new responsible user not found: %w", err)
	}

	// 3. 执行转移
	return task.AssignResponsible(newResponsibleID, transferredBy)
}

// BulkUpdateTaskStatus 批量更新任务状态
func (s *TaskDomainServiceImpl) BulkUpdateTaskStatus(ctx context.Context, taskIDs []valueobject.TaskID, newStatus valueobject.TaskStatus, updatedBy valueobject.UserID) error {
	for _, taskID := range taskIDs {
		task, err := s.taskRepo.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to find task %s: %w", taskID, err)
		}

		// 验证状态转换
		if err := s.ValidateStatusTransition(*task, task.Status, newStatus, updatedBy); err != nil {
			return fmt.Errorf("invalid status transition for task %s: %w", taskID, err)
		}

		// 更新状态（简化实现，实际应该调用具体的状态转换方法）
		// 这里需要根据newStatus调用相应的方法
	}

	return nil
}

// ArchiveCompletedTasks 归档已完成任务
func (s *TaskDomainServiceImpl) ArchiveCompletedTasks(ctx context.Context, projectID valueobject.ProjectID, archivedBy valueobject.UserID) error {
	// 1. 查找项目中已完成的任务
	tasks, err := s.taskRepo.FindByProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to find tasks in project: %w", err)
	}

	// 2. 筛选已完成的任务并归档
	for _, task := range tasks {
		if task.Status == valueobject.TaskStatusCompleted {
			// 实际实现中应该有归档操作，这里简化处理
			// 可以添加归档标记或移动到归档表
		}
	}

	return nil
}

// CanUserCreateTaskInProject 检查用户是否可以在项目中创建任务
func (s *TaskDomainServiceImpl) CanUserCreateTaskInProject(userID valueobject.UserID, projectID valueobject.ProjectID) bool {
	// 简化实现：检查用户是否为项目成员
	project, err := s.projectRepo.FindByID(context.Background(), projectID)
	if err != nil {
		return false
	}

	return project.CanUserAccess(userID)
}

// CanUserManageTask 检查用户是否可以管理任务
func (s *TaskDomainServiceImpl) CanUserManageTask(userID valueobject.UserID, task aggregate.TaskAggregate) bool {
	return task.CanUserModify(userID)
}

// GetUserTaskPermissions 获取用户任务权限
func (s *TaskDomainServiceImpl) GetUserTaskPermissions(userID valueobject.UserID, task aggregate.TaskAggregate) valueobject.TaskPermissions {
	permissions := valueobject.TaskPermissions{
		CanView:    task.CanUserView(userID),
		CanModify:  task.CanUserModify(userID),
		CanExecute: task.CanUserExecute(userID),
		CanApprove: task.CanUserApprove(userID),
		CanDelete:  task.CanUserModify(userID), // 简化实现
		CanAssign:  task.CanUserModify(userID), // 简化实现
	}

	return permissions
}

// GetNextApprovers 获取下一个审批者
func (s *TaskDomainServiceImpl) GetNextApprovers(task aggregate.TaskAggregate) ([]valueobject.UserID, error) {
	// 简化实现：返回任务创建者作为审批者
	return []valueobject.UserID{task.CreatorID}, nil
}

// ValidateWorkflowTransition 验证工作流转换
func (s *TaskDomainServiceImpl) ValidateWorkflowTransition(task aggregate.TaskAggregate, fromStatus, toStatus valueobject.TaskStatus) error {
	// 委托给状态转换验证
	return s.ValidateStatusTransition(task, fromStatus, toStatus, task.CreatorID)
}

// ProcessWorkflowStep 处理工作流步骤
func (s *TaskDomainServiceImpl) ProcessWorkflowStep(ctx context.Context, task aggregate.TaskAggregate, stepData valueobject.WorkflowStepData) error {
	// 简化实现：根据步骤动作执行相应操作
	switch stepData.Action {
	case "approval":
		// 处理审批步骤
		return task.Approve(stepData.ActorID, stepData.Comment)
	case "assignment":
		// 处理分配步骤 - 从Data中获取目标用户ID
		if targetUserID, ok := stepData.Data["target_user_id"].(string); ok {
			return task.AssignResponsible(valueobject.UserID(targetUserID), stepData.ActorID)
		}
		return fmt.Errorf("target_user_id not found in workflow step data")
	default:
		return fmt.Errorf("unsupported workflow step action: %s", stepData.Action)
	}
}
