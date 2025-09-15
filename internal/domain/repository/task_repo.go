package repository

import (
	"context"
	"time"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

// TaskRepository 任务仓储接口
type TaskRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, task aggregate.TaskAggregate) error
	FindByID(ctx context.Context, id valueobject.TaskID) (*aggregate.TaskAggregate, error)
	FindByIDs(ctx context.Context, ids []valueobject.TaskID) ([]aggregate.TaskAggregate, error)
	Delete(ctx context.Context, id valueobject.TaskID) error

	// 查询方法
	FindByProject(ctx context.Context, projectID valueobject.ProjectID) ([]aggregate.TaskAggregate, error)
	FindByCreator(ctx context.Context, creatorID valueobject.UserID) ([]aggregate.TaskAggregate, error)
	FindByResponsible(ctx context.Context, responsibleID valueobject.UserID) ([]aggregate.TaskAggregate, error)
	FindByParticipant(ctx context.Context, participantID valueobject.UserID) ([]aggregate.TaskAggregate, error)
	FindByStatus(ctx context.Context, status valueobject.TaskStatus) ([]aggregate.TaskAggregate, error)
	FindByPriority(ctx context.Context, priority valueobject.TaskPriority) ([]aggregate.TaskAggregate, error)
	FindByType(ctx context.Context, taskType valueobject.TaskType) ([]aggregate.TaskAggregate, error)

	// 复杂查询
	SearchTasks(ctx context.Context, criteria valueobject.TaskSearchCriteria) ([]aggregate.TaskAggregate, int, error)
	FindOverdueTasks(ctx context.Context, asOfDate time.Time) ([]aggregate.TaskAggregate, error)
	FindTasksDueWithin(ctx context.Context, duration time.Duration) ([]aggregate.TaskAggregate, error)
	FindUserAccessibleTasks(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]aggregate.TaskAggregate, int, error)

	// 统计查询
	CountByProject(ctx context.Context, projectID valueobject.ProjectID) (int, error)
	CountByStatus(ctx context.Context, status valueobject.TaskStatus) (int, error)
	CountByResponsible(ctx context.Context, responsibleID valueobject.UserID) (int, error)
	GetTaskStatistics(ctx context.Context, taskID valueobject.TaskID) (*valueobject.TaskStatistics, error)
	GetProjectTaskStatistics(ctx context.Context, projectID valueobject.ProjectID) (*valueobject.ProjectTaskStatistics, error)
}
