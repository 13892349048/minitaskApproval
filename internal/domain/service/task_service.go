package service

import (
	"context"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

// TaskDomainService 任务领域服务接口
type TaskDomainService interface {
	// 业务规则验证
	ValidateTaskCreation(task aggregate.TaskAggregate, createdBy valueobject.UserID) error
	ValidateTaskAssignment(task aggregate.TaskAggregate, responsibleID valueobject.UserID, assignedBy valueobject.UserID) error
	ValidateParticipantAddition(task aggregate.TaskAggregate, participantID valueobject.UserID, addedBy valueobject.UserID) error
	ValidateStatusTransition(task aggregate.TaskAggregate, fromStatus, toStatus valueobject.TaskStatus, changedBy valueobject.UserID) error
	ValidateTaskCompletion(task aggregate.TaskAggregate, completedBy valueobject.UserID) error

	// 复杂业务逻辑
	TransferTaskResponsibility(ctx context.Context, task aggregate.TaskAggregate, newResponsibleID valueobject.UserID, transferredBy valueobject.UserID) error
	BulkUpdateTaskStatus(ctx context.Context, taskIDs []valueobject.TaskID, newStatus valueobject.TaskStatus, updatedBy valueobject.UserID) error
	ArchiveCompletedTasks(ctx context.Context, projectID valueobject.ProjectID, archivedBy valueobject.UserID) error

	// 权限相关
	CanUserCreateTaskInProject(userID valueobject.UserID, projectID valueobject.ProjectID) bool
	CanUserManageTask(userID valueobject.UserID, task aggregate.TaskAggregate) bool
	GetUserTaskPermissions(userID valueobject.UserID, task aggregate.TaskAggregate) valueobject.TaskPermissions

	// 工作流相关
	GetNextApprovers(task aggregate.TaskAggregate) ([]valueobject.UserID, error)
	ValidateWorkflowTransition(task aggregate.TaskAggregate, fromStatus, toStatus valueobject.TaskStatus) error
	ProcessWorkflowStep(ctx context.Context, task aggregate.TaskAggregate, stepData valueobject.WorkflowStepData) error
}
