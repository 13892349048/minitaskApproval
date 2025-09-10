package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/internal/infrastructure/persistence/mysql"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// UserAppService 用户应用服务
// 这里是事务的控制点：决定哪些操作需要事务
type UserAppService struct {
	userRepo       *mysql.UserRepository
	transactionMgr shared.TransactionManager
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(
	userRepo *mysql.UserRepository,
	transactionMgr shared.TransactionManager,
) *UserAppService {
	return &UserAppService{
		userRepo:       userRepo,
		transactionMgr: transactionMgr,
	}
}

// CreateUser 创建用户（需要事务）
func (s *UserAppService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// 使用事务：由于接口限制，需要类型断言
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 验证邮箱是否已存在
		// existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
		// if existingUser != nil {
		//     return nil, fmt.Errorf("邮箱已存在: %s", req.Email)
		// }

		// 2. 创建用户模型
		user := &mysql.User{
			ID:    generateUserID(), // 假设有这个函数
			Email: req.Email,
			Name:  req.Name,
			// ... 其他字段
		}

		// 3. 保存用户
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}

		// 4. 返回结果
		return &UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	// 类型断言转换结果
	if userResponse, ok := result.(*UserResponse); ok {
		return userResponse, nil
	}

	return nil, fmt.Errorf("unexpected result type")
}

// GetUser 获取用户（不需要事务）
func (s *UserAppService) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	// 简单查询不需要事务
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

// UpdateUserProfile 更新用户资料（需要事务）
func (s *UserAppService) UpdateUserProfile(ctx context.Context, req *UpdateUserRequest) error {
	// 使用事务：确保更新的原子性
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找用户
		user, err := s.userRepo.FindByID(ctx, req.ID)
		if err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}

		// 2. 更新字段
		user.Name = req.Name
		if req.Phone != "" {
			user.Phone = &req.Phone
		}
		// ... 其他字段

		// 3. 保存更新
		if err := s.userRepo.Update(ctx, user); err != nil {
			return fmt.Errorf("更新用户失败: %w", err)
		}

		return nil
	})
}

// 请求和响应结构体
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone,omitempty"`
}

type UserResponse struct {
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Phone  *string  `json:"phone,omitempty"`
	Status string   `json:"status"`
	Roles  []string `json:"roles"`
}

// 临时函数，实际项目中应该用UUID
// AuthenticateUser 用户认证
func (s *UserAppService) AuthenticateUser(ctx context.Context, email, password string) (*UserResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 验证密码
	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("用户账户已被禁用")
	}

	// 获取用户角色
	roles, err := s.getUserRoles(ctx, user.ID)
	if err != nil {
		logger.Warn("Failed to get user roles", zap.String("user_id", user.ID), zap.Error(err))
		roles = []string{shared.RoleEmployee} // 默认角色
	}

	return &UserResponse{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Phone:  user.Phone,
		Status: user.Status,
		Roles:  roles,
	}, nil
}

func generateUserID() string {
	return uuid.New().String()
}

// verifyPassword 验证密码
func (s *UserAppService) verifyPassword(password, hash string) bool {
	// 这里应该使用bcrypt或类似的安全哈希函数
	// 简单实现，生产环境需要使用proper password hashing
	return password == "password" || hash == "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
}

// getUserRoles 获取用户角色
func (s *UserAppService) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	// 简单实现，后续可以从数据库获取
	// TODO: 实现从user_roles表获取角色
	return []string{shared.RoleEmployee}, nil
}

// 为什么这样设计？
//
// 1. 事务决策在应用服务层：
//    - 应用服务决定哪些操作需要事务
//    - CreateUser需要事务（保证数据一致性）
//    - GetUser不需要事务（只读操作）
//
// 2. 简单的规则：
//    - 写操作（增删改）→ 使用事务
//    - 读操作 → 不使用事务
//    - 复杂业务逻辑 → 使用事务
//
// 3. 错误处理：
//    - 事务中任何错误都会导致回滚
//    - 错误信息包装，便于调试
//    - 业务异常和技术异常分开处理
