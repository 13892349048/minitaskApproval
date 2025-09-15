package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/taskflow/internal/domain/aggregate"
	authService "github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/domain/repository"
	"github.com/taskflow/internal/domain/service"
	"github.com/taskflow/internal/domain/valueobject"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// UserAppService 用户应用服务
// 这里是事务的控制点：决定哪些操作需要事务
type UserAppService struct {
	userDomainService service.UserDomainService
	transactionMgr    authService.TransactionManager
	uv                service.UserValidator
	userRepo          repository.UserRepository
	passwordHasher    service.PasswordHasher
}

// NewUserAppService 创建用户应用服务
func NewUserAppService(
	userDomainService service.UserDomainService,
	transactionMgr authService.TransactionManager,
	uv service.UserValidator,
	userRepo repository.UserRepository,
	passwordHasher service.PasswordHasher,
) *UserAppService {
	return &UserAppService{
		userDomainService: userDomainService,
		transactionMgr:    transactionMgr,
		uv:                uv,
		passwordHasher:    passwordHasher,
		userRepo:          userRepo,
	}
}

// CreateUser 创建用户（需要事务）
func (s *UserAppService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// 使用事务：由于接口限制，需要类型断言
	result, err := s.transactionMgr.WithTransactionResult(ctx, func(ctx context.Context) (interface{}, error) {
		// 1. 验证邮箱是否已存在
		existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
		if existingUser != nil {
			return nil, fmt.Errorf("邮箱已存在: %s", req.Email)
		}

		// 2. 验证输入数据
		if err := s.uv.ValidateEmail(req.Email); err != nil {
			return nil, fmt.Errorf("邮箱验证失败: %w", err)
		}
		if err := s.uv.ValidatePassword(req.Password); err != nil {
			return nil, fmt.Errorf("密码验证失败: %w", err)
		}
		if err := s.uv.ValidateName(req.Name); err != nil {
			return nil, fmt.Errorf("姓名验证失败: %w", err)
		}

		// 3. 哈希密码
		passwordHash, err := s.passwordHasher.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("密码哈希失败: %w", err)
		}

		// 4. 创建用户聚合根
		user := aggregate.NewUser(
			valueobject.UserID(generateUserID()),
			req.Name,
			req.Email,
			req.Name, // FullName
			passwordHash,
			valueobject.UserRoleEmployee,
		)

		// 5. 保存用户
		if err := s.userRepo.Save(ctx, user); err != nil {
			return nil, fmt.Errorf("创建用户失败: %w", err)
		}

		// 6. 返回结果
		return &UserResponse{
			ID:     string(user.ID),
			Email:  user.Email,
			Name:   user.Username,
			Phone:  &req.Phone,
			Status: string(user.Status),
			Roles:  []string{string(user.Role)},
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
		ID:    string(user.ID),
		Email: user.Email,
		Name:  user.Username,
	}, nil
}

// ListUsers 获取用户列表（不需要事务）
func (s *UserAppService) ListUsers(ctx context.Context, req *ListUsersRequest) ([]*UserResponse, int, error) {
	// 构建搜索条件
	criteria := valueobject.UserSearchCriteria{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	// 设置角色过滤条件
	if req.Role != "" {
		role := valueobject.UserRole(req.Role)
		criteria.Role = &role
	}

	// 设置状态过滤条件
	if req.Status != "" {
		status := valueobject.UserStatus(req.Status)
		criteria.Status = &status
	}

	// 查询用户列表
	users, total, err := s.userRepo.SearchUsers(ctx, criteria)
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 转换为响应格式
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		// 获取用户角色
		roles, err := s.getUserRoles(ctx, string(user.ID))
		if err != nil {
			// 如果获取角色失败，使用默认角色
			roles = []string{string(user.Role)}
		}

		responses[i] = &UserResponse{
			ID:     string(user.ID),
			Email:  user.Email,
			Name:   user.Username,
			Phone:  nil, // TODO: 需要在User聚合中添加Phone字段
			Status: string(user.Status),
			Roles:  roles,
		}
	}

	return responses, total, nil
}

// DeleteUser 删除用户（软删除，需要事务）
func (s *UserAppService) DeleteUser(ctx context.Context, id string) error {
	return s.transactionMgr.WithTransaction(ctx, func(ctx context.Context) error {
		// 1. 查找用户
		user, err := s.userRepo.FindByID(ctx, id)
		if err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}

		// 2. 执行软删除（停用用户）
		user.Deactivate()

		// 3. 保存更新
		if err := s.userRepo.Update(ctx, user); err != nil {
			return fmt.Errorf("删除用户失败: %w", err)
		}

		return nil
	})
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

		// 2. 更新字段 - 使用Domain方法
		if err := user.UpdateProfile(req.Name, user.Email); err != nil {
			return fmt.Errorf("更新用户资料失败: %w", err)
		}

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

type ListUsersRequest struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Role   string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
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
	if !s.passwordHasher.VerifyPassword(user.PasswordHash, password) {
		return nil, fmt.Errorf("密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("用户账户已被禁用")
	}

	// 获取用户角色
	roles, err := s.getUserRoles(ctx, string(user.ID))
	if err != nil {
		logger.Warn("Failed to get user roles", zap.String("user_id", string(user.ID)), zap.Error(err))
		roles = []string{string(user.Role)} // 使用用户当前角色作为默认
	}

	return &UserResponse{
		ID:     string(user.ID),
		Email:  user.Email,
		Name:   user.Username,
		Phone:  nil, // TODO: 需要在User聚合中添加Phone字段
		Status: string(user.Status),
		Roles:  roles,
	}, nil
}

func generateUserID() string {
	return uuid.New().String()
}

// getUserRoles 获取用户角色
func (s *UserAppService) getUserRoles(ctx context.Context, userID string) ([]string, error) {
	// 简单实现，后续可以从数据库获取
	// TODO: 实现从user_roles表获取角色
	return []string{string(valueobject.UserRoleEmployee)}, nil
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
