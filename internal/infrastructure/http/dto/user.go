package dto

import "time"

// UserLoginRequest 用户登录请求
// @Description 用户登录请求参数
type UserLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`     // 用户名
	Password string `json:"password" binding:"required" example:"password"`  // 密码
} // @name UserLoginRequest

// UserLoginResponse 用户登录响应
// @Description 用户登录响应数据
type UserLoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`  // JWT令牌
	ExpiresAt time.Time `json:"expires_at" example:"2023-01-01T00:00:00Z"`                // 令牌过期时间
	User      UserInfo  `json:"user"`                                                     // 用户信息
} // @name UserLoginResponse

// UserInfo 用户信息
// @Description 用户基本信息
type UserInfo struct {
	BaseEntity
	Username    string   `json:"username" example:"admin"`                    // 用户名
	Email       string   `json:"email" example:"admin@example.com"`           // 邮箱
	DisplayName string   `json:"display_name" example:"管理员"`                 // 显示名称
	Avatar      string   `json:"avatar" example:"https://example.com/avatar"` // 头像URL
	Status      string   `json:"status" example:"active"`                     // 状态
	Roles       []string `json:"roles" example:"admin,user"`                  // 角色列表
} // @name UserInfo

// UserCreateRequest 创建用户请求
// @Description 创建用户请求参数
type UserCreateRequest struct {
	Username    string   `json:"username" binding:"required" example:"newuser"`        // 用户名
	Email       string   `json:"email" binding:"required,email" example:"user@example.com"` // 邮箱
	Password    string   `json:"password" binding:"required,min=6" example:"password123"`   // 密码
	DisplayName string   `json:"display_name" example:"新用户"`                           // 显示名称
	Roles       []string `json:"roles" example:"user"`                                 // 角色列表
} // @name UserCreateRequest

// UserUpdateRequest 更新用户请求
// @Description 更新用户请求参数
type UserUpdateRequest struct {
	Email       string   `json:"email" binding:"omitempty,email" example:"user@example.com"` // 邮箱
	DisplayName string   `json:"display_name" example:"更新的用户"`                             // 显示名称
	Avatar      string   `json:"avatar" example:"https://example.com/new-avatar"`            // 头像URL
	Status      string   `json:"status" example:"active"`                                    // 状态
	Roles       []string `json:"roles" example:"user,manager"`                               // 角色列表
} // @name UserUpdateRequest

// ChangePasswordRequest 修改密码请求
// @Description 修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword"`  // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newpassword123"` // 新密码
} // @name ChangePasswordRequest
