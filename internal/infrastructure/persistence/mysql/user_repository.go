package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	*BaseRepository // 嵌入基础仓储，自动获得事务支持
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	var user User
	// 使用GetDB(ctx)自动获取正确的数据库连接（事务或非事务）
	if err := r.GetDB(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, user *User) error {
	// 自动支持事务：如果在事务中，会使用事务连接
	if err := r.GetDB(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	if err := r.GetDB(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	if err := r.GetDB(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.GetDB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if err := r.GetDB(ctx).Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// 使用示例：
//
// 1. 非事务操作：
//    user, err := userRepo.FindByID(ctx, "user123")
//    // 直接使用普通数据库连接
//
// 2. 事务操作：
//    err := transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
//        user, err := userRepo.FindByID(ctx, "user123") // 自动使用事务连接
//        if err != nil {
//            return err
//        }
//        user.Name = "新名字"
//        return userRepo.Update(ctx, user) // 自动使用事务连接
//    })
//
// Repository完全不需要知道自己是否在事务中！
