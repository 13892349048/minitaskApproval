package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
	"gorm.io/gorm"
)

// TaskRepositoryImpl 任务仓储实现
type TaskRepositoryImpl struct {
	*BaseRepository
}

// NewTaskRepository 创建任务仓储
func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &TaskRepositoryImpl{
		BaseRepository: NewBaseRepository(db),
	}
}

// TaskPO 任务持久化对象
type TaskPO struct {
	ID             string     `gorm:"primaryKey;column:id" json:"id"`
	Title          string     `gorm:"column:title;not null" json:"title"`
	Description    string     `gorm:"column:description;type:text" json:"description"`
	ProjectID      string     `gorm:"column:project_id;not null;index" json:"project_id"`
	CreatorID      string     `gorm:"column:creator_id;not null;index" json:"creator_id"`
	AssigneeID     *string    `gorm:"column:assignee_id;index" json:"assignee_id"`
	Status         string     `gorm:"column:status;not null;index" json:"status"`
	Priority       string     `gorm:"column:priority;not null" json:"priority"`
	Type           string     `gorm:"column:type;not null" json:"type"`
	StartDate      *time.Time `gorm:"column:start_date" json:"start_date"`
	DueDate        *time.Time `gorm:"column:due_date;index" json:"due_date"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completed_at"`
	EstimatedHours *float64   `gorm:"column:estimated_hours" json:"estimated_hours"`
	ActualHours    *float64   `gorm:"column:actual_hours" json:"actual_hours"`
	Tags           string     `gorm:"column:tags;type:json" json:"tags"`
	Participants   string     `gorm:"column:participants;type:json" json:"participants"`
	Attachments    string     `gorm:"column:attachments;type:json" json:"attachments"`
	RecurrenceRule *string    `gorm:"column:recurrence_rule" json:"recurrence_rule"`
	ParentTaskID   *string    `gorm:"column:parent_task_id;index" json:"parent_task_id"`
	WorkflowStepID *string    `gorm:"column:workflow_step_id" json:"workflow_step_id"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 表名
func (TaskPO) TableName() string {
	return "tasks"
}

// Save 保存任务
func (r *TaskRepositoryImpl) Save(ctx context.Context, task aggregate.TaskAggregate) error {
	po := r.aggregateToTaskPO(task)
	return r.db.WithContext(ctx).Create(&po).Error
}

// FindByID 根据ID查找任务
func (r *TaskRepositoryImpl) FindByID(ctx context.Context, id valueobject.TaskID) (*aggregate.TaskAggregate, error) {
	var po TaskPO
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", string(id)).First(&po).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOToAggregate(po), nil
}

// Update 更新任务
func (r *TaskRepositoryImpl) Update(ctx context.Context, task aggregate.TaskAggregate) error {
	po := r.aggregateToTaskPO(task)
	return r.db.WithContext(ctx).Where("id = ?", po.ID).Updates(&po).Error
}

// Delete 删除任务
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id valueobject.TaskID) error {
	return r.db.WithContext(ctx).Model(&TaskPO{}).Where("id = ?", string(id)).Update("deleted_at", time.Now()).Error
}

// FindByProjectID 根据项目ID查找任务
func (r *TaskRepositoryImpl) FindByProjectID(ctx context.Context, projectID valueobject.ProjectID) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("project_id = ? AND deleted_at IS NULL", string(projectID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// FindByAssigneeID 根据负责人ID查找任务
func (r *TaskRepositoryImpl) FindByAssigneeID(ctx context.Context, assigneeID valueobject.UserID) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("assignee_id = ? AND deleted_at IS NULL", string(assigneeID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// FindByCreatorID 根据创建者ID查找任务
func (r *TaskRepositoryImpl) FindByCreatorID(ctx context.Context, creatorID valueobject.UserID) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("creator_id = ? AND deleted_at IS NULL", string(creatorID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// FindByDateRange 根据日期范围查找任务
func (r *TaskRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// Search 搜索任务
func (r *TaskRepositoryImpl) Search(ctx context.Context, criteria valueobject.TaskSearchCriteria) ([]*aggregate.TaskAggregate, error) {
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if criteria.ProjectID != nil {
		query = query.Where("project_id = ?", string(*criteria.ProjectID))
	}
	if criteria.ResponsibleID != nil {
		query = query.Where("assignee_id = ?", string(*criteria.ResponsibleID))
	}
	if criteria.CreatorID != nil {
		query = query.Where("creator_id = ?", string(*criteria.CreatorID))
	}
	if criteria.Status != nil {
		query = query.Where("status = ?", string(*criteria.Status))
	}
	if criteria.Priority != nil {
		query = query.Where("priority = ?", string(*criteria.Priority))
	}
	if criteria.TaskType != nil {
		query = query.Where("type = ?", string(*criteria.TaskType))
	}
	if criteria.Title != nil {
		query = query.Where("title LIKE ?", "%"+*criteria.Title+"%")
	}
	if criteria.Description != nil {
		query = query.Where("description LIKE ?", "%"+*criteria.Description+"%")
	}

	var pos []TaskPO
	err := query.Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// Count 统计任务数量
func (r *TaskRepositoryImpl) Count(ctx context.Context, criteria valueobject.TaskSearchCriteria) (int64, error) {
	query := r.db.WithContext(ctx).Model(&TaskPO{}).Where("deleted_at IS NULL")

	if criteria.ProjectID != nil {
		query = query.Where("project_id = ?", string(*criteria.ProjectID))
	}
	if criteria.ResponsibleID != nil {
		query = query.Where("assignee_id = ?", string(*criteria.ResponsibleID))
	}
	if criteria.CreatorID != nil {
		query = query.Where("creator_id = ?", string(*criteria.CreatorID))
	}
	if criteria.Status != nil {
		query = query.Where("status = ?", string(*criteria.Status))
	}
	if criteria.Priority != nil {
		query = query.Where("priority = ?", string(*criteria.Priority))
	}
	if criteria.TaskType != nil {
		query = query.Where("type = ?", string(*criteria.TaskType))
	}
	if criteria.Title != nil {
		query = query.Where("title LIKE ?", "%"+*criteria.Title+"%")
	}
	if criteria.Description != nil {
		query = query.Where("description LIKE ?", "%"+*criteria.Description+"%")
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// FindWithPagination 分页查找任务
func (r *TaskRepositoryImpl) FindWithPagination(ctx context.Context, criteria valueobject.TaskSearchCriteria, offset, limit int) ([]*aggregate.TaskAggregate, int64, error) {
	// 先获取总数
	total, err := r.Count(ctx, criteria)
	if err != nil {
		return nil, 0, err
	}

	// 构建查询
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if criteria.ProjectID != nil {
		query = query.Where("project_id = ?", string(*criteria.ProjectID))
	}
	if criteria.ResponsibleID != nil {
		query = query.Where("assignee_id = ?", string(*criteria.ResponsibleID))
	}
	if criteria.CreatorID != nil {
		query = query.Where("creator_id = ?", string(*criteria.CreatorID))
	}
	if criteria.Status != nil {
		query = query.Where("status = ?", string(*criteria.Status))
	}
	if criteria.Priority != nil {
		query = query.Where("priority = ?", string(*criteria.Priority))
	}
	if criteria.TaskType != nil {
		query = query.Where("type = ?", string(*criteria.TaskType))
	}
	if criteria.Title != nil {
		query = query.Where("title LIKE ?", "%"+*criteria.Title+"%")
	}
	if criteria.Description != nil {
		query = query.Where("description LIKE ?", "%"+*criteria.Description+"%")
	}

	var pos []TaskPO
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, 0, err
	}

	return r.taskPOsToAggregates(pos), total, nil
}

// FindByParticipantID 根据参与者ID查找任务
func (r *TaskRepositoryImpl) FindByParticipantID(ctx context.Context, participantID valueobject.UserID) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("JSON_CONTAINS(participants, ?) AND deleted_at IS NULL", fmt.Sprintf(`"%s"`, string(participantID))).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// FindOverdueTasks 查找过期任务
// func (r *TaskRepositoryImpl) FindOverdueTasks(ctx context.Context) ([]*aggregate.TaskAggregate, error) {
// 	var pos []TaskPO
// 	err := r.db.WithContext(ctx).Where("due_date < ? AND status NOT IN (?, ?) AND deleted_at IS NULL",
// 		time.Now(), string(valueobject.TaskStatusCompleted), string(valueobject.TaskStatusCancelled)).Find(&pos).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r.taskPOsToAggregates(pos), nil
// }

// FindUpcomingTasks 查找即将到期的任务
func (r *TaskRepositoryImpl) FindUpcomingTasks(ctx context.Context, days int) ([]*aggregate.TaskAggregate, error) {
	upcomingDate := time.Now().AddDate(0, 0, days)
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("due_date BETWEEN ? AND ? AND status NOT IN (?, ?) AND deleted_at IS NULL",
		time.Now(), upcomingDate, string(valueobject.TaskStatusCompleted), string(valueobject.TaskStatusCancelled)).Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// FindRecurringTasks 查找循环任务
func (r *TaskRepositoryImpl) FindRecurringTasks(ctx context.Context) ([]*aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("recurrence_rule IS NOT NULL AND deleted_at IS NULL").Find(&pos).Error
	if err != nil {
		return nil, err
	}
	return r.taskPOsToAggregates(pos), nil
}

// BatchSave 批量保存任务
func (r *TaskRepositoryImpl) BatchSave(ctx context.Context, tasks []*aggregate.TaskAggregate) error {
	pos := make([]TaskPO, len(tasks))
	for i, task := range tasks {
		pos[i] = r.aggregateToTaskPO(*task)
	}
	return r.db.WithContext(ctx).CreateInBatches(pos, 100).Error
}

// BatchUpdate 批量更新任务
func (r *TaskRepositoryImpl) BatchUpdate(ctx context.Context, tasks []*aggregate.TaskAggregate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, task := range tasks {
			po := r.aggregateToTaskPO(*task)
			if err := tx.Where("id = ?", po.ID).Updates(&po).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除任务
func (r *TaskRepositoryImpl) BatchDelete(ctx context.Context, ids []valueobject.TaskID) error {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = string(id)
	}
	return r.db.WithContext(ctx).Model(&TaskPO{}).Where("id IN ?", strIDs).Update("deleted_at", time.Now()).Error
}

// aggregateToTaskPO 将聚合根转换为持久化对象
func (r *TaskRepositoryImpl) aggregateToTaskPO(task aggregate.TaskAggregate) TaskPO {
	po := TaskPO{
		ID:        string(task.ID),
		Title:     task.Title,
		ProjectID: string(task.ProjectID),
		CreatorID: string(task.CreatorID),
		Status:    string(task.Status),
		Priority:  string(task.Priority),
		Type:      string(task.TaskType),
		DueDate:   task.DueDate,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	// 处理可选的Description字段
	if task.Description != nil {
		po.Description = *task.Description
	}

	// 处理ResponsibleID
	assigneeID := string(task.ResponsibleID)
	po.AssigneeID = &assigneeID

	// 处理EstimatedHours转换
	if task.EstimatedHours > 0 {
		hours := float64(task.EstimatedHours)
		po.EstimatedHours = &hours
	}

	// 处理ActualHours
	if task.ActualHours > 0 {
		po.ActualHours = &task.ActualHours
	}

	return po
}

// taskPOToAggregate 将持久化对象转换为聚合根
func (r *TaskRepositoryImpl) taskPOToAggregate(po TaskPO) *aggregate.TaskAggregate {
	task := &aggregate.TaskAggregate{
		ID:           valueobject.TaskID(po.ID),
		Title:        po.Title,
		ProjectID:    valueobject.ProjectID(po.ProjectID),
		CreatorID:    valueobject.UserID(po.CreatorID),
		Status:       valueobject.TaskStatus(po.Status),
		Priority:     valueobject.TaskPriority(po.Priority),
		TaskType:     valueobject.TaskType(po.Type),
		DueDate:      po.DueDate,
		WorkflowID:   "",
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
		Participants: make([]valueobject.TaskParticipant, 0),
		Events:       make([]event.DomainEvent, 0),
	}

	// 处理可选的Description字段
	if po.Description != "" {
		task.Description = &po.Description
	}

	// 处理ResponsibleID
	if po.AssigneeID != nil {
		task.ResponsibleID = valueobject.UserID(*po.AssigneeID)
	}

	// 处理EstimatedHours转换
	if po.EstimatedHours != nil {
		task.EstimatedHours = int(*po.EstimatedHours)
	}

	// 处理ActualHours
	if po.ActualHours != nil {
		task.ActualHours = *po.ActualHours
	}

	return task
}

// taskPOsToAggregates 将持久化对象数组转换为聚合根数组
func (r *TaskRepositoryImpl) taskPOsToAggregates(pos []TaskPO) []*aggregate.TaskAggregate {
	aggregates := make([]*aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = r.taskPOToAggregate(po)
	}
	return aggregates
}

// FindByIDs 根据ID列表查找任务
func (r *TaskRepositoryImpl) FindByIDs(ctx context.Context, ids []valueobject.TaskID) ([]aggregate.TaskAggregate, error) {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = string(id)
	}

	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", strIDs).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByProject 根据项目ID查找任务
func (r *TaskRepositoryImpl) FindByProject(ctx context.Context, projectID valueobject.ProjectID) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("project_id = ? AND deleted_at IS NULL", string(projectID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByCreator 根据创建者ID查找任务
func (r *TaskRepositoryImpl) FindByCreator(ctx context.Context, creatorID valueobject.UserID) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("creator_id = ? AND deleted_at IS NULL", string(creatorID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByResponsible 根据负责人ID查找任务
func (r *TaskRepositoryImpl) FindByResponsible(ctx context.Context, responsibleID valueobject.UserID) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("assignee_id = ? AND deleted_at IS NULL", string(responsibleID)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByParticipant 根据参与者ID查找任务
func (r *TaskRepositoryImpl) FindByParticipant(ctx context.Context, participantID valueobject.UserID) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("JSON_CONTAINS(participants, ?) AND deleted_at IS NULL", fmt.Sprintf(`"%s"`, string(participantID))).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByStatus 根据状态查找任务
func (r *TaskRepositoryImpl) FindByStatus(ctx context.Context, status valueobject.TaskStatus) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("status = ? AND deleted_at IS NULL", string(status)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByPriority 根据优先级查找任务
func (r *TaskRepositoryImpl) FindByPriority(ctx context.Context, priority valueobject.TaskPriority) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("priority = ? AND deleted_at IS NULL", string(priority)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindByType 根据类型查找任务
func (r *TaskRepositoryImpl) FindByType(ctx context.Context, taskType valueobject.TaskType) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("type = ? AND deleted_at IS NULL", string(taskType)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// FindOverdueTasks 查找过期任务
func (r *TaskRepositoryImpl) FindOverdueTasks(ctx context.Context, asOfDate time.Time) ([]aggregate.TaskAggregate, error) {
	var pos []TaskPO
	err := r.db.WithContext(ctx).Where("due_date < ? AND status NOT IN (?, ?) AND deleted_at IS NULL",
		asOfDate, string(valueobject.TaskStatusCompleted), string(valueobject.TaskStatusCancelled)).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	aggregates := make([]aggregate.TaskAggregate, len(pos))
	for i, po := range pos {
		aggregates[i] = *r.taskPOToAggregate(po)
	}
	return aggregates, nil
}

// SearchTasks 搜索任务
func (r *TaskRepositoryImpl) SearchTasks(ctx context.Context, criteria valueobject.TaskSearchCriteria) ([]aggregate.TaskAggregate, int, error) {
	return nil, 0, fmt.Errorf("not implemented yet")
}

// FindTasksDueWithin 查找指定时间内到期的任务
func (r *TaskRepositoryImpl) FindTasksDueWithin(ctx context.Context, duration time.Duration) ([]aggregate.TaskAggregate, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// FindUserAccessibleTasks 查找用户可访问的任务
func (r *TaskRepositoryImpl) FindUserAccessibleTasks(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]aggregate.TaskAggregate, int, error) {
	return nil, 0, fmt.Errorf("not implemented yet")
}

// CountByProject 按项目统计任务数量
func (r *TaskRepositoryImpl) CountByProject(ctx context.Context, projectID valueobject.ProjectID) (int, error) {
	return 0, fmt.Errorf("not implemented yet")
}

// CountByStatus 按状态统计任务数量
func (r *TaskRepositoryImpl) CountByStatus(ctx context.Context, status valueobject.TaskStatus) (int, error) {
	return 0, fmt.Errorf("not implemented yet")
}

// CountByResponsible 按负责人统计任务数量
func (r *TaskRepositoryImpl) CountByResponsible(ctx context.Context, responsibleID valueobject.UserID) (int, error) {
	return 0, fmt.Errorf("not implemented yet")
}

// GetTaskStatistics 获取任务统计信息
func (r *TaskRepositoryImpl) GetTaskStatistics(ctx context.Context, taskID valueobject.TaskID) (*valueobject.TaskStatistics, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// GetProjectTaskStatistics 获取项目任务统计信息
func (r *TaskRepositoryImpl) GetProjectTaskStatistics(ctx context.Context, projectID valueobject.ProjectID) (*valueobject.ProjectTaskStatistics, error) {
	return nil, fmt.Errorf("not implemented yet")
}
