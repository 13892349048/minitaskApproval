package mysql

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/valueobject"
	"gorm.io/gorm"
)

// DepartmentRepositoryImpl 部门仓储实现
type DepartmentRepositoryImpl struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓储实例
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepositoryImpl {
	return &DepartmentRepositoryImpl{db: db}
}

// FindByID 根据ID查找部门
func (r *DepartmentRepositoryImpl) FindByID(ctx context.Context, id valueobject.DepartmentID) (*valueobject.DepartmentInfo, error) {
	// 简化实现 - 返回默认部门信息
	// 在实际项目中，这里应该查询departments表
	if id == "" {
		return nil, fmt.Errorf("department not found: %s", id)
	}
	
	return &valueobject.DepartmentInfo{
		ID:   id,
		Name: fmt.Sprintf("Department_%s", id),
	}, nil
}

// FindByUserID 根据用户ID查找部门
func (r *DepartmentRepositoryImpl) FindByUserID(ctx context.Context, userID valueobject.UserID) (*valueobject.DepartmentInfo, error) {
	// 简化实现 - 从users表查询department_id
	var departmentID string
	err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Select("department_id").
		Where("id = ?", string(userID)).
		First(&departmentID).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user department not found: %s", userID)
		}
		return nil, fmt.Errorf("failed to find user department: %w", err)
	}
	
	if departmentID == "" {
		return nil, fmt.Errorf("user has no department assigned: %s", userID)
	}
	
	return r.FindByID(ctx, valueobject.DepartmentID(departmentID))
}

// IsActive 检查部门是否活跃
func (r *DepartmentRepositoryImpl) IsActive(ctx context.Context, id valueobject.DepartmentID) (bool, error) {
	// 简化实现 - 假设所有部门都是活跃的
	// 在实际项目中，这里应该查询departments表的status字段
	if id == "" {
		return false, nil
	}
	return true, nil
}
