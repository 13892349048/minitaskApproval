package mysql

import (
	"context"
	"fmt"

	"github.com/taskflow/internal/domain/aggregate"
	"github.com/taskflow/internal/domain/valueobject"
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储实现 - 实现Domain层接口
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实现
func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

// Save 保存用户
func (r *UserRepositoryImpl) Save(ctx context.Context, domainUser *aggregate.User) error {
	userModel := r.domainToModel(domainUser)

	if err := r.db.WithContext(ctx).Save(userModel).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, domainUser *aggregate.User) error {
	UserModel := r.domainToModel(domainUser)

	if err := r.db.WithContext(ctx).Updates(UserModel).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*aggregate.User, error) {
	var userModel UserModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
	var userModel UserModel

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	var userModel UserModel

	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return r.modelToDomain(&userModel), nil
}

// Delete 删除用户（软删除）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// FindByRole 根据角色查找用户
func (r *UserRepositoryImpl) FindByRole(ctx context.Context, role valueobject.UserRole) ([]*aggregate.User, error) {
	var userModels []UserModel

	if err := r.db.WithContext(ctx).Where("role = ?", string(role)).Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, nil
}

// FindByStatus 根据状态查找用户
func (r *UserRepositoryImpl) FindByStatus(ctx context.Context, status valueobject.UserStatus) ([]*aggregate.User, error) {
	var userModels []UserModel

	if err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by status: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, nil
}

// FindByDepartment 根据部门查找用户
func (r *UserRepositoryImpl) FindByDepartment(ctx context.Context, departmentID string) ([]*aggregate.User, error) {
	var userModels []UserModel

	if err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by department: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, nil
}

// FindByManager 根据管理者查找用户
func (r *UserRepositoryImpl) FindByManager(ctx context.Context, managerID valueobject.UserID) ([]*aggregate.User, error) {
	var userModels []UserModel

	if err := r.db.WithContext(ctx).Where("manager_id = ?", string(managerID)).Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by manager: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, nil
}

// SearchUsers 搜索用户
func (r *UserRepositoryImpl) SearchUsers(ctx context.Context, criteria valueobject.UserSearchCriteria) ([]*aggregate.User, int, error) {
	query := r.db.WithContext(ctx).Model(&UserModel{})

	// 构建查询条件
	if criteria.Username != nil {
		query = query.Where("username LIKE ?", "%"+*criteria.Username+"%")
	}
	if criteria.Email != nil {
		query = query.Where("email LIKE ?", "%"+*criteria.Email+"%")
	}
	if criteria.FullName != nil {
		query = query.Where("full_name LIKE ?", "%"+*criteria.FullName+"%")
	}
	if criteria.Role != nil {
		query = query.Where("role = ?", string(*criteria.Role))
	}
	if criteria.Status != nil {
		query = query.Where("status = ?", string(*criteria.Status))
	}
	if criteria.DepartmentID != nil {
		query = query.Where("department_id = ?", *criteria.DepartmentID)
	}
	if criteria.ManagerID != nil {
		query = query.Where("manager_id = ?", criteria.ManagerID.String())
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 排序和分页
	if criteria.OrderBy != "" {
		orderDir := "ASC"
		if criteria.OrderDir == "DESC" {
			orderDir = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", criteria.OrderBy, orderDir))
	}

	if criteria.Limit > 0 {
		query = query.Limit(criteria.Limit)
	}
	if criteria.Offset > 0 {
		query = query.Offset(criteria.Offset)
	}

	// 执行查询
	var userModels []UserModel
	if err := query.Find(&userModels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, int(total), nil
}

// FindUsersByRole 根据角色名称查找用户（分页）
func (r *UserRepositoryImpl) FindUsersByRole(ctx context.Context, roleName string, limit, offset int) ([]*aggregate.User, int, error) {
	var userModels []UserModel
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("role = ?", roleName).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	// 查询分页数据
	if err := r.db.WithContext(ctx).Where("role = ?", roleName).
		Limit(limit).Offset(offset).Find(&userModels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find users by role: %w", err)
	}

	users := make([]*aggregate.User, len(userModels))
	for i, model := range userModels {
		users[i] = r.modelToDomain(&model)
	}

	return users, int(total), nil
}

// CountByStatus 根据状态统计用户数量
func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status valueobject.UserStatus) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("status = ?", string(status)).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}

	return int(count), nil
}

func (r *UserRepositoryImpl) CountByDepartment(ctx context.Context, department string) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("department = ?", department).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users by department: %w", err)
	}

	return int(count), nil
}

// domainToModel 将Domain User转换为持久化模型
func (r *UserRepositoryImpl) domainToModel(domainUser *aggregate.User) *UserModel {
	model := &UserModel{
		ID:        string(domainUser.ID),
		Username:  domainUser.Username,
		Email:     domainUser.Email,
		FullName:  domainUser.FullName,
		Role:      string(domainUser.Role),
		Status:    string(domainUser.Status),
		CreatedAt: domainUser.CreatedAt,
		UpdatedAt: domainUser.UpdatedAt,
	}

	if domainUser.DepartmentID != nil {
		model.DepartmentID = domainUser.DepartmentID
	}

	if domainUser.ManagerID != nil {
		managerID := string(*domainUser.ManagerID)
		model.ManagerID = &managerID
	}

	if domainUser.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *domainUser.DeletedAt, Valid: true}
	}

	return model
}

// modelToDomain 将持久化模型转换为Domain User
func (r *UserRepositoryImpl) modelToDomain(model *UserModel) *aggregate.User {
	domainUser := aggregate.NewUser(
		valueobject.UserID(model.ID),
		model.Username,
		model.Email,
		model.FullName,
		model.PasswordHash,
		valueobject.UserRole(model.Role),
	)

	// 设置其他属性 - 直接访问公共字段
	if model.DepartmentID != nil {
		domainUser.DepartmentID = model.DepartmentID
	}

	if model.ManagerID != nil {
		managerID := valueobject.UserID(*model.ManagerID)
		domainUser.ManagerID = &managerID
	}

	// 设置状态
	switch valueobject.UserStatus(model.Status) {
	case valueobject.UserStatusActive:
	case valueobject.UserStatusInactive:
		domainUser.Deactivate()
	case valueobject.UserStatusSuspended:
		domainUser.Suspend()
	}

	return domainUser
}
