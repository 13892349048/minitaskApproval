package repository

import (
	"context"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
)

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, project aggregate.Project) error
	FindByID(ctx context.Context, id valueobject.ProjectID) (*aggregate.Project, error)
	FindByIDs(ctx context.Context, ids []valueobject.ProjectID) ([]aggregate.Project, error)
	Delete(ctx context.Context, id valueobject.ProjectID) error

	// 查询方法
	FindByOwner(ctx context.Context, ownerID valueobject.UserID) ([]aggregate.Project, error)
	FindByManager(ctx context.Context, managerID valueobject.UserID) ([]aggregate.Project, error)
	FindByMember(ctx context.Context, userID valueobject.UserID) ([]aggregate.Project, error)
	FindByParent(ctx context.Context, parentID valueobject.ProjectID) ([]aggregate.Project, error)
	FindByStatus(ctx context.Context, status valueobject.ProjectStatus) ([]aggregate.Project, error)
	FindByType(ctx context.Context, projectType valueobject.ProjectType) ([]aggregate.Project, error)

	// 复杂查询
	SearchProjects(ctx context.Context, criteria aggregate.ProjectSearchCriteria) ([]aggregate.Project, int, error)
	FindUserAccessibleProjects(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]aggregate.Project, int, error)

	// 统计查询
	CountByOwner(ctx context.Context, ownerID valueobject.UserID) (int, error)
	CountByStatus(ctx context.Context, status valueobject.ProjectStatus) (int, error)
	GetProjectStatistics(ctx context.Context, projectID valueobject.ProjectID) (*aggregate.ProjectStatistics, error)
}
