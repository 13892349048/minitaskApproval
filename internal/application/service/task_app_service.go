package service

import (
	"context"
	"fmt"
	"time"

	"github.com/taskflow/internal/application/dto"
	"github.com/taskflow/internal/domain/aggregate"
	authService "github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/service"
	"github.com/taskflow/internal/domain/valueobject"
)

// TaskAppService 任务应用服务
type TaskAppService struct {
	taskDomainService service.TaskDomainService
	transactionMgr    authService.TransactionManager
	taskRepo          repository.TaskRepository
	taskFactory       *aggregate.TaskFactory
}

// NewTaskAppService 创建任务应用服务
func NewTaskAppService(
	taskDomainService service.TaskDomainService,
	transactionMgr authService.TransactionManager,
	taskRepo repository.TaskRepository,
	taskFactory *aggregate.TaskFactory,
) *TaskAppService {
	return &TaskAppService{
		taskDomainService: taskDomainService,
		transactionMgr:    transactionMgr,
		taskRepo:          taskRepo,
		taskFactory:       taskFactory,
	}
}

// CreateTask 创建任务（需要事务）
func (s *TaskAppService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 创建任务聚合
		task, err := s.taskFactory.CreateTask(
			valueobject.TaskID(""), // Generate ID in factory
			req.Title,
			s.stringPtrToString(req.Description),
			valueobject.TaskType(req.TaskType),
			valueobject.TaskPriority(req.Priority),
			valueobject.ProjectID(req.ProjectID),
			valueobject.UserID(req.CreatorID),
			valueobject.UserID(req.ResponsibleID),
			req.DueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("创建任务失败: %w", err)
		}

		// 2. 保存任务
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return nil, fmt.Errorf("保存任务失败: %w", err)
		}

		// 3. 返回结果
		return &dto.CreateTaskResponse{
			ID:            string((*task).ID),
			Title:         (*task).Title,
			Description:   (*task).Description,
			TaskType:      string((*task).TaskType),
			Priority:      string((*task).Priority),
			Status:        string((*task).Status),
			ProjectID:     string((*task).ProjectID),
			CreatorID:     string((*task).CreatorID),
			ResponsibleID: string((*task).ResponsibleID),
			DueDate:       (*task).DueDate,
			CreatedAt:     (*task).CreatedAt,
			UpdatedAt:     (*task).UpdatedAt,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	if taskResponse, ok := result.(*dto.CreateTaskResponse); ok {
		return taskResponse, nil
	}

	return nil, fmt.Errorf("unexpected result type")
}

// GetTask 获取任务（不需要事务）
func (s *TaskAppService) GetTask(ctx context.Context, id string) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(id))
	if err != nil {
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}

	return &dto.TaskResponse{
		ID:            string(task.ID),
		Title:         task.Title,
		Description:   task.Description,
		TaskType:      string(task.TaskType),
		Priority:      string(task.Priority),
		Status:        string(task.Status),
		ProjectID:     string(task.ProjectID),
		CreatorID:     string(task.CreatorID),
		ResponsibleID: string(task.ResponsibleID),
		DueDate:       task.DueDate,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}, nil
}

// UpdateTask 更新任务（需要事务）
func (s *TaskAppService) UpdateTask(ctx context.Context, req dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 查找任务
		task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(req.ID))
		if err != nil {
			return nil, fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 更新任务信息
		title := task.Title
		if req.Title != nil {
			title = *req.Title
		}
		description := s.stringPtrToString(task.Description)
		if req.Description != nil {
			description = *req.Description
		}
		if err := task.UpdateBasicInfo(title, description); err != nil {
			return nil, fmt.Errorf("更新任务信息失败: %w", err)
		}

		// 3. 保存更新
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return nil, fmt.Errorf("保存任务失败: %w", err)
		}

		// 4. 返回更新后的任务
		return &dto.UpdateTaskResponse{
			ID:            string(task.ID),
			Title:         task.Title,
			Description:   task.Description,
			TaskType:      string(task.TaskType),
			Priority:      string(task.Priority),
			Status:        string(task.Status),
			ProjectID:     string(task.ProjectID),
			CreatorID:     string(task.CreatorID),
			ResponsibleID: string(task.ResponsibleID),
			DueDate:       task.DueDate,
			EstimatedHours: task.EstimatedHours,
			CreatedAt:     task.CreatedAt,
			UpdatedAt:     task.UpdatedAt,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	if updateResponse, ok := result.(*dto.UpdateTaskResponse); ok {
		return updateResponse, nil
	}

	return nil, fmt.Errorf("unexpected result type")
}

// AssignTask 分配任务（需要事务）
func (s *TaskAppService) AssignTask(ctx context.Context, req dto.AssignTaskRequest) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找任务
		task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(req.TaskID))
		if err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 分配负责人
		if err := task.AssignResponsible(
			valueobject.UserID(req.ResponsibleID),
			valueobject.UserID(req.AssignedBy),
		); err != nil {
			return fmt.Errorf("分配任务失败: %w", err)
		}

		// 3. 保存更新
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}

		return nil
	})
}

// DeleteTask 删除任务（需要事务）
func (s *TaskAppService) DeleteTask(ctx context.Context, taskID valueobject.TaskID) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 验证任务存在
		_, err := s.taskRepo.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 删除任务
		if err := s.taskRepo.Delete(ctx, taskID); err != nil {
			return fmt.Errorf("删除任务失败: %w", err)
		}

		return nil
	})
}

// ListTasks 获取任务列表
func (s *TaskAppService) ListTasks(ctx context.Context, req dto.ListTasksRequest) (*dto.ListTasksResponse, error) {
	// 转换搜索条件
	criteria := s.convertSearchCriteria(req.Criteria)
	
	// 查询任务
	tasks, total, err := s.taskRepo.SearchTasks(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	// 转换为响应DTO
	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		participants := make([]dto.TaskParticipantDTO, len(task.Participants))
		for j, p := range task.Participants {
			participants[j] = dto.TaskParticipantDTO{
				UserID:  string(p.UserID),
				Role:    string(p.Role),
				AddedAt: p.AddedAt,
				AddedBy: string(p.AddedBy),
			}
		}

		taskResponses[i] = dto.TaskResponse{
			ID:            string(task.ID),
			Title:         task.Title,
			Description:   task.Description,
			TaskType:      string(task.TaskType),
			Priority:      string(task.Priority),
			Status:        string(task.Status),
			ProjectID:     string(task.ProjectID),
			CreatorID:     string(task.CreatorID),
			ResponsibleID: string(task.ResponsibleID),
			DueDate:       task.DueDate,
			EstimatedHours: task.EstimatedHours,
			ActualHours:   task.ActualHours,
			Participants:  participants,
			CreatedAt:     task.CreatedAt,
			UpdatedAt:     task.UpdatedAt,
		}
	}

	// 计算总页数
	totalPages := int((int64(total) + int64(req.PageSize) - 1) / int64(req.PageSize))

	response := &dto.ListTasksResponse{
		Tasks:      taskResponses,
		Total:      int64(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// UpdateTaskStatus 更新任务状态（需要事务）
func (s *TaskAppService) UpdateTaskStatus(ctx context.Context, req dto.UpdateTaskStatusRequest) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找任务
		task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(req.TaskID))
		if err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		userID := valueobject.UserID(req.UpdatedBy)
		status := valueobject.TaskStatus(req.Status)

		// 2. 根据状态执行相应操作
		switch status {
		case valueobject.TaskStatusDraft:
			// 草稿状态 - 通常不需要特殊处理
		case valueobject.TaskStatusPendingApproval:
			err = task.SubmitForApproval(userID)
		case valueobject.TaskStatusApproved:
			err = task.Approve(userID, req.Comment)
		case valueobject.TaskStatusRejected:
			err = task.Reject(userID, req.Comment)
		case valueobject.TaskStatusInProgress:
			err = task.Start(userID)
		case valueobject.TaskStatusPaused:
			err = task.Pause(userID, req.Comment)
		case valueobject.TaskStatusCompleted:
			err = task.Complete(userID)
		case valueobject.TaskStatusCancelled:
			err = task.Cancel(userID, req.Comment)
		default:
			return fmt.Errorf("不支持的状态: %s", status)
		}

		if err != nil {
			return fmt.Errorf("更新任务状态失败: %w", err)
		}

		// 3. 保存更新
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}

		return nil
	})
}

// AddTaskParticipant 添加任务参与者（需要事务）
func (s *TaskAppService) AddTaskParticipant(ctx context.Context, req dto.AddTaskParticipantRequest) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找任务
		task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(req.TaskID))
		if err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 添加参与者
		if err := task.AddParticipant(valueobject.UserID(req.ParticipantID), valueobject.UserID(req.AddedBy)); err != nil {
			return fmt.Errorf("添加参与者失败: %w", err)
		}

		// 3. 保存更新
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}

		return nil
	})
}

// RemoveTaskParticipant 移除任务参与者（需要事务）
func (s *TaskAppService) RemoveTaskParticipant(ctx context.Context, req dto.RemoveTaskParticipantRequest) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找任务
		task, err := s.taskRepo.FindByID(ctx, valueobject.TaskID(req.TaskID))
		if err != nil {
			return fmt.Errorf("任务不存在: %w", err)
		}

		// 2. 移除参与者
		if err := task.RemoveParticipant(valueobject.UserID(req.ParticipantID), valueobject.UserID(req.RemovedBy)); err != nil {
			return fmt.Errorf("移除参与者失败: %w", err)
		}

		// 3. 保存更新
		if err := s.taskRepo.Save(ctx, *task); err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}

		return nil
	})
}

// GetTaskStatistics 获取任务统计信息
func (s *TaskAppService) GetTaskStatistics(ctx context.Context, projectID *valueobject.ProjectID) (*dto.TaskStatisticsResponse, error) {
	// 构建搜索条件
	criteria := valueobject.TaskSearchCriteria{}
	if projectID != nil {
		criteria.ProjectID = projectID
	}

	// 获取所有任务
	tasks, _, err := s.taskRepo.SearchTasks(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("查询任务失败: %w", err)
	}

	// 计算统计信息
	stats := &dto.TaskStatisticsResponse{
		TotalTasks:      len(tasks),
		TasksByStatus:   make(map[string]int),
		TasksByPriority: make(map[string]int),
		TasksByType:     make(map[string]int),
	}

	var totalHours float64
	var completedTasks int
	overdueTasks := 0

	for _, task := range tasks {
		// 按状态统计
		stats.TasksByStatus[string(task.Status)]++
		
		// 按优先级统计
		stats.TasksByPriority[string(task.Priority)]++
		
		// 按类型统计
		stats.TasksByType[string(task.TaskType)]++
		
		// 计算完成率
		if task.Status == valueobject.TaskStatusCompleted {
			completedTasks++
		}
		
		// 计算过期任务
		if task.DueDate != nil && task.DueDate.Before(time.Now()) && 
		   task.Status != valueobject.TaskStatusCompleted && 
		   task.Status != valueobject.TaskStatusCancelled {
			overdueTasks++
		}
		
		// 累计工时
		totalHours += task.ActualHours
	}

	stats.OverdueTasks = overdueTasks
	
	// 计算完成率
	if stats.TotalTasks > 0 {
		stats.CompletionRate = float64(completedTasks) / float64(stats.TotalTasks) * 100
	}
	
	// 计算平均工时
	if stats.TotalTasks > 0 {
		stats.AverageHours = totalHours / float64(stats.TotalTasks)
	}

	return stats, nil
}

// convertSearchCriteria 转换搜索条件
func (s *TaskAppService) convertSearchCriteria(dto dto.TaskSearchCriteria) valueobject.TaskSearchCriteria {
	return valueobject.TaskSearchCriteria{
		Title:         dto.Title,
		Description:   dto.Description,
		TaskType:      dto.TaskType,
		Priority:      dto.Priority,
		Status:        dto.Status,
		ProjectID:     dto.ProjectID,
		CreatorID:     dto.CreatorID,
		ResponsibleID: dto.ResponsibleID,
		ParticipantID: dto.ParticipantID,
		StartDate:     dto.StartDate,
		DueDate:       dto.DueDate,
		CreatedAfter:  dto.CreatedAfter,
		CreatedBefore: dto.CreatedBefore,
	}
}

// stringPtrToString 将字符串指针转换为字符串
func (s *TaskAppService) stringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
