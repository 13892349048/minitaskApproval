package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/valueobject"
	"github.com/taskflow/internal/infrastructure/persistence/cache"
	"gorm.io/gorm"
)

// ProjectRepository 项目仓储实现 - 基于现有架构扩展
type ProjectRepository struct {
	*BaseRepository // 嵌入基础仓储，自动获得事务支持
	cache           cache.Interface
	cacheTTL        time.Duration
	event.TransactionManager
}

// NewProjectRepository 创建项目仓储
func NewProjectRepository(db *gorm.DB, cache cache.Interface) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: NewBaseRepository(db),
		cache:          cache,
		cacheTTL:       30 * time.Minute,
	}
}

// Save 保存项目 - 写入数据库，清除缓存
func (r *ProjectRepository) Save(ctx context.Context, proj aggregate.Project) error {

	// 转换为数据库模型
	projectModel := r.aggregateToModel(proj)

	// 使用GetDB自动支持事务
	if err := r.GetDB(ctx).Save(projectModel).Error; err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	// 保存项目成员
	if err := r.saveProjectMembers(ctx, proj); err != nil {
		return fmt.Errorf("failed to save project members: %w", err)
	}

	// 异步清除缓存
	go r.invalidateCache(ctx, proj.ID)

	return nil
}

// FindByID 查找项目 - 先查缓存，再查数据库
func (r *ProjectRepository) FindByID(ctx context.Context, id valueobject.ProjectID) (*aggregate.Project, error) {

	// 1. 尝试从缓存获取
	if proj, err := r.getFromCache(ctx, id); err == nil {
		return proj, nil
	}

	// 2. 从数据库查询
	var projectModel Project
	if err := r.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&projectModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	// 3. 加载关联数据
	if err := r.loadProjectMembers(ctx, &projectModel); err != nil {
		return nil, fmt.Errorf("failed to load project members: %w", err)
	}

	// 4. 转换为聚合根
	proj := r.modelToAggregate(&projectModel)

	// 5. 异步写入缓存
	go r.setCache(ctx, *proj)

	return proj, nil
}

// FindByOwner 查找用户拥有的项目
func (r *ProjectRepository) FindByOwner(ctx context.Context, ownerID valueobject.UserID) ([]aggregate.Project, error) {

	var projectModels []Project
	if err := r.GetDB(ctx).Where("owner_id = ? AND deleted_at IS NULL", ownerID).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by owner: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindUserAccessibleProjects 查找用户可访问的项目
func (r *ProjectRepository) FindUserAccessibleProjects(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]aggregate.Project, int, error) {

	// 复杂查询使用原生SQL，但仍通过GORM执行
	query := `
		SELECT DISTINCT p.*, COUNT(*) OVER() as total_count
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		WHERE p.deleted_at IS NULL 
		  AND (p.owner_id = ? OR p.manager_id = ? OR pm.user_id = ?)
		ORDER BY p.updated_at DESC
		LIMIT ? OFFSET ?
	`

	var results []struct {
		Project
		TotalCount int `gorm:"column:total_count"`
	}

	if err := r.GetDB(ctx).Raw(query, userID, userID, userID, limit, offset).Scan(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find user accessible projects: %w", err)
	}

	if len(results) == 0 {
		return []aggregate.Project{}, 0, nil
	}

	// 提取项目和总数
	var projectModels []Project
	totalCount := results[0].TotalCount

	for _, result := range results {
		projectModels = append(projectModels, result.Project)
	}

	return r.modelsToAggregates(projectModels), totalCount, nil
}

// Delete 软删除项目
func (r *ProjectRepository) Delete(ctx context.Context, id valueobject.ProjectID) error {

	now := time.Now()
	if err := r.GetDB(ctx).Model(&Project{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	// 清除缓存
	go r.invalidateCache(ctx, id)

	return nil
}

// 私有方法 - 缓存操作

func (r *ProjectRepository) getFromCache(ctx context.Context, id valueobject.ProjectID) (*aggregate.Project, error) {
	if r.cache == nil {
		return nil, fmt.Errorf("cache not available")
	}

	key := fmt.Sprintf("project:%s", id)
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var projectData aggregate.ProjectData
	if err := json.Unmarshal([]byte(data), &projectData); err != nil {
		return nil, err
	}

	// 使用工厂恢复项目
	factory := aggregate.NewProjectFactory()
	return factory.RestoreProject(projectData), nil
}

func (r *ProjectRepository) setCache(ctx context.Context, proj aggregate.Project) error {
	if r.cache == nil {
		return nil // 缓存不可用时静默失败
	}

	key := fmt.Sprintf("project:%s", proj.ID)
	data := r.aggregateToData(proj)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.cache.Set(ctx, key, string(jsonData), r.cacheTTL)
}

func (r *ProjectRepository) invalidateCache(ctx context.Context, id valueobject.ProjectID) {
	if r.cache != nil {
		key := fmt.Sprintf("project:%s", id)
		r.cache.Del(ctx, key)
	}
}

// 私有方法 - 数据转换

func (r *ProjectRepository) aggregateToModel(proj aggregate.Project) *Project {
	model := &Project{
		ID:          string(proj.ID),
		Name:        proj.Name,
		Description: &proj.Description,
		ProjectType: string(proj.ProjectType),
		Status:      string(proj.Status),
		OwnerID:     string(proj.OwnerID),
		StartDate:   &proj.StartDate,
		CreatedAt:   proj.CreatedAt,
		UpdatedAt:   proj.UpdatedAt,
	}

	// 处理DeletedAt
	if proj.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *proj.DeletedAt, Valid: true}
	}

	if proj.ParentID != nil {
		parentID := string(*proj.ParentID)
		model.ParentProjectID = &parentID
	}

	if proj.ManagerID != nil {
		managerID := string(*proj.ManagerID)
		model.ManagerID = &managerID
	}

	if proj.EndDate != nil {
		model.EndDate = proj.EndDate
	}

	return model
}

func (r *ProjectRepository) modelToAggregate(model *Project) *aggregate.Project {
	// 这里需要实现从数据库模型到聚合根的转换
	// 由于聚合根构造函数是私有的，需要使用工厂方法
	data := aggregate.ProjectData{
		ID:          model.ID,
		Name:        model.Name,
		Description: "",
		Type:        model.ProjectType,
		Status:      model.Status,
		OwnerID:     model.OwnerID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	if model.Description != nil {
		data.Description = *model.Description
	}

	if model.DeletedAt.Valid {
		data.DeletedAt = &model.DeletedAt.Time
	}

	if model.StartDate != nil {
		data.StartDate = *model.StartDate
	}

	if model.EndDate != nil {
		data.EndDate = model.EndDate
	}

	if model.ParentProjectID != nil {
		data.ParentID = model.ParentProjectID
	}

	if model.ManagerID != nil {
		data.ManagerID = model.ManagerID
	}

	factory := aggregate.NewProjectFactory()
	return factory.RestoreProject(data)
}

func (r *ProjectRepository) modelsToAggregates(models []Project) []aggregate.Project {
	aggregates := make([]aggregate.Project, len(models))
	for i, model := range models {
		aggregates[i] = *r.modelToAggregate(&model)
	}
	return aggregates
}

func (r *ProjectRepository) aggregateToData(proj aggregate.Project) aggregate.ProjectData {
	data := aggregate.ProjectData{
		ID:          string(proj.ID),
		Name:        proj.Name,
		Description: "",
		Type:        string(proj.ProjectType),
		Status:      string(proj.Status),
		OwnerID:     string(proj.OwnerID),
		CreatedAt:   proj.CreatedAt,
		UpdatedAt:   proj.UpdatedAt,
		DeletedAt:   proj.DeletedAt,
	}

	if proj.Description != "" {
		data.Description = proj.Description
	}

	if proj.StartDate.IsZero() {
		data.StartDate = proj.StartDate
	}

	if proj.ParentID != nil {
		parentID := string(*proj.ParentID)
		data.ParentID = &parentID
	}

	if proj.ManagerID != nil {
		managerID := string(*proj.ManagerID)
		data.ManagerID = &managerID
	}

	if proj.EndDate != nil {
		data.EndDate = proj.EndDate
	}

	// 转换成员数据
	for _, member := range proj.Members {
		memberData := aggregate.ProjectMemberData{
			UserID:   string(member.UserID),
			Role:     string(member.Role),
			JoinedAt: member.JoinedAt,
			AddedBy:  string(member.AddedBy),
		}
		data.Members = append(data.Members, memberData)
	}

	// 转换子项目ID
	for _, childID := range proj.Children {
		data.Children = append(data.Children, string(childID))
	}

	return data
}

// 成员管理相关方法

func (r *ProjectRepository) saveProjectMembers(ctx context.Context, proj aggregate.Project) error {
	// 先删除现有成员
	if err := r.GetDB(ctx).Where("project_id = ?", proj.ID).Delete(&ProjectMember{}).Error; err != nil {
		return err
	}
	var value string
	// 插入新成员
	for _, member := range proj.Members {
		value = string(member.AddedBy)
		memberModel := &ProjectMember{
			ID:        generateID(), // 需要实现ID生成函数
			ProjectID: string(proj.ID),
			UserID:    string(member.UserID),
			Role:      string(member.Role),
			JoinedAt:  member.JoinedAt,
			AddedBy:   &value,
		}

		if err := r.GetDB(ctx).Create(memberModel).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ProjectRepository) loadProjectMembers(ctx context.Context, projectModel *Project) error {
	var memberModels []ProjectMember
	if err := r.GetDB(ctx).Where("project_id = ?", projectModel.ID).Find(&memberModels).Error; err != nil {
		return err
	}

	// 这里可以将成员数据设置到项目模型中，或者在转换时处理
	// 具体实现取决于Project模型的设计

	return nil
}

// 辅助函数
func generateID() string {
	return uuid.New().String()
}

// FindByIDs 批量查找项目
func (r *ProjectRepository) FindByIDs(ctx context.Context, ids []valueobject.ProjectID) ([]aggregate.Project, error) {
	var projectModels []Project

	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = string(id)
	}

	if err := r.GetDB(ctx).Where("id IN ? AND deleted_at IS NULL", strIDs).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by IDs: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindByManager 查找管理者管理的项目
func (r *ProjectRepository) FindByManager(ctx context.Context, managerID valueobject.UserID) ([]aggregate.Project, error) {
	var projectModels []Project

	if err := r.GetDB(ctx).Where("manager_id = ? AND deleted_at IS NULL", managerID).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by manager: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindByMember 查找用户参与的项目
func (r *ProjectRepository) FindByMember(ctx context.Context, userID valueobject.UserID) ([]aggregate.Project, error) {
	var projectModels []Project

	query := `
		SELECT DISTINCT p.*
		FROM projects p
		INNER JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = ? AND p.deleted_at IS NULL
	`

	if err := r.GetDB(ctx).Raw(query, userID).Scan(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by member: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindByParent 查找子项目
func (r *ProjectRepository) FindByParent(ctx context.Context, parentID valueobject.ProjectID) ([]aggregate.Project, error) {
	var projectModels []Project

	if err := r.GetDB(ctx).Where("parent_project_id = ? AND deleted_at IS NULL", parentID).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by parent: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindByStatus 按状态查找项目
func (r *ProjectRepository) FindByStatus(ctx context.Context, status valueobject.ProjectStatus) ([]aggregate.Project, error) {
	var projectModels []Project

	if err := r.GetDB(ctx).Where("status = ? AND deleted_at IS NULL", status).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by status: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// FindByType 按类型查找项目
func (r *ProjectRepository) FindByType(ctx context.Context, projectType valueobject.ProjectType) ([]aggregate.Project, error) {
	var projectModels []Project

	if err := r.GetDB(ctx).Where("project_type = ? AND deleted_at IS NULL", projectType).Find(&projectModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find projects by type: %w", err)
	}

	return r.modelsToAggregates(projectModels), nil
}

// SearchProjects 复杂搜索项目
func (r *ProjectRepository) SearchProjects(ctx context.Context, criteria aggregate.ProjectSearchCriteria) ([]aggregate.Project, int, error) {
	db := r.GetDB(ctx).Model(&Project{})

	// 构建查询条件
	db = db.Where("deleted_at IS NULL")

	if criteria.Name != nil {
		db = db.Where("name LIKE ?", "%"+*criteria.Name+"%")
	}
	if criteria.Description != nil {
		db = db.Where("description LIKE ?", "%"+*criteria.Description+"%")
	}
	if criteria.ProjectType != nil {
		db = db.Where("project_type = ?", *criteria.ProjectType)
	}
	if criteria.Status != nil {
		db = db.Where("status = ?", *criteria.Status)
	}
	if criteria.OwnerID != nil {
		db = db.Where("owner_id = ?", *criteria.OwnerID)
	}
	if criteria.ManagerID != nil {
		db = db.Where("manager_id = ?", *criteria.ManagerID)
	}
	if criteria.ParentID != nil {
		db = db.Where("parent_project_id = ?", *criteria.ParentID)
	}

	// 计算总数
	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	// 排序和分页
	orderBy := "created_at"
	orderDir := "DESC"
	if criteria.OrderBy != "" {
		orderBy = criteria.OrderBy
	}
	if criteria.OrderDir != "" {
		orderDir = criteria.OrderDir
	}

	db = db.Order(fmt.Sprintf("%s %s", orderBy, orderDir))

	if criteria.Limit > 0 {
		db = db.Limit(criteria.Limit)
	}
	if criteria.Offset > 0 {
		db = db.Offset(criteria.Offset)
	}

	var projectModels []Project
	if err := db.Find(&projectModels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search projects: %w", err)
	}

	return r.modelsToAggregates(projectModels), int(totalCount), nil
}

// CountByOwner 统计用户拥有的项目数
func (r *ProjectRepository) CountByOwner(ctx context.Context, ownerID valueobject.UserID) (int, error) {
	var count int64

	if err := r.GetDB(ctx).Model(&Project{}).Where("owner_id = ? AND deleted_at IS NULL", ownerID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count projects by owner: %w", err)
	}

	return int(count), nil
}

// CountByStatus 统计指定状态的项目数
func (r *ProjectRepository) CountByStatus(ctx context.Context, status valueobject.ProjectStatus) (int, error) {
	var count int64

	if err := r.GetDB(ctx).Model(&Project{}).Where("status = ? AND deleted_at IS NULL", status).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count projects by status: %w", err)
	}

	return int(count), nil
}

// GetProjectStatistics 获取项目统计信息
func (r *ProjectRepository) GetProjectStatistics(ctx context.Context, projectID valueobject.ProjectID) (*aggregate.ProjectStatistics, error) {

	// 这里需要根据实际的任务表结构来实现
	// 暂时返回基本统计信息
	stats := &aggregate.ProjectStatistics{
		ProjectID:       projectID,
		TotalTasks:      0,
		CompletedTasks:  0,
		InProgressTasks: 0,
		PendingTasks:    0,
		OverdueTasks:    0,
		TotalMembers:    0,
		ActiveMembers:   0,
		CompletionRate:  0.0,
		AverageTaskTime: 0.0,
		LastActivityAt:  time.Now(),
	}

	// 统计项目成员数
	var memberCount int64
	if err := r.GetDB(ctx).Model(&ProjectMember{}).Where("project_id = ?", projectID).Count(&memberCount).Error; err == nil {
		stats.TotalMembers = int(memberCount)
		stats.ActiveMembers = int(memberCount) // 简化处理
	}

	return stats, nil
}

// 确保实现了接口
var _ repository.ProjectRepository = (*ProjectRepository)(nil)
