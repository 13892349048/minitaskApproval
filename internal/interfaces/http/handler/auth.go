package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/internal/application/shared"
	"github.com/taskflow/internal/application/user"
	"github.com/taskflow/pkg/errors"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	jwtService  shared.JWTService
	userService *user.UserAppService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(jwtService shared.JWTService, userService *user.UserAppService) *AuthHandler {
	return &AuthHandler{
		jwtService:  jwtService,
		userService: userService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Phone    string `json:"phone,omitempty"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User   *UserInfo         `json:"user"`
	Tokens *shared.TokenPair `json:"tokens"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Phone  *string  `json:"phone,omitempty"`
	Roles  []string `json:"roles"`
	Status string   `json:"status"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户使用邮箱和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} AuthResponse "登录成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "认证失败"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 验证用户凭据
	userResp, err := h.userService.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logger.Warn("Login failed",
			zap.String("email", req.Email),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "邮箱或密码错误")
		return
	}

	// 生成JWT令牌
	tokens, err := h.jwtService.GenerateTokens(userResp.ID, userResp.Email, userResp.Roles)
	if err != nil {
		logger.Error("Failed to generate tokens",
			zap.String("user_id", userResp.ID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "令牌生成失败")
		return
	}

	// 记录登录日志
	logger.Info("User logged in successfully",
		zap.String("user_id", userResp.ID),
		zap.String("email", userResp.Email))

	// 返回响应
	response := &AuthResponse{
		User: &UserInfo{
			ID:     userResp.ID,
			Name:   userResp.Name,
			Email:  userResp.Email,
			Phone:  userResp.Phone,
			Roles:  userResp.Roles,
			Status: userResp.Status,
		},
		Tokens: tokens,
	}

	errors.RespondWithSuccess(c, response, "登录成功")
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} AuthResponse "注册成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 409 {object} errors.ErrorResponse "邮箱已存在"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 创建用户
	createReq := &user.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}

	userResp, err := h.userService.CreateUser(c.Request.Context(), createReq)
	if err != nil {
		logger.Warn("Registration failed",
			zap.String("email", req.Email),
			zap.Error(err))

		// 检查是否是邮箱重复错误
		if isEmailExistsError(err) {
			errors.RespondWithError(c, http.StatusConflict, "EMAIL_EXISTS", "邮箱已存在")
			return
		}

		errors.RespondWithError(c, http.StatusInternalServerError, "REGISTRATION_FAILED", "注册失败")
		return
	}

	// 生成JWT令牌
	tokens, err := h.jwtService.GenerateTokens(userResp.ID, userResp.Email, []string{shared.RoleEmployee})
	if err != nil {
		logger.Error("Failed to generate tokens after registration",
			zap.String("user_id", userResp.ID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "令牌生成失败")
		return
	}

	// 记录注册日志
	logger.Info("User registered successfully",
		zap.String("user_id", userResp.ID),
		zap.String("email", userResp.Email))

	// 返回响应
	response := &AuthResponse{
		User: &UserInfo{
			ID:     userResp.ID,
			Name:   userResp.Name,
			Email:  userResp.Email,
			Phone:  &req.Phone,
			Roles:  []string{shared.RoleEmployee},
			Status: "active",
		},
		Tokens: tokens,
	}

	c.JSON(http.StatusCreated, errors.SuccessResponse{
		Success: true,
		Data:    response,
		Message: "注册成功",
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} shared.TokenPair "新的令牌对"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "刷新令牌无效"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 刷新令牌
	tokens, err := h.jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Warn("Token refresh failed", zap.Error(err))
		errors.RespondWithError(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "刷新令牌无效或已过期")
		return
	}

	logger.Debug("Token refreshed successfully")
	errors.RespondWithSuccess(c, tokens, "令牌刷新成功")
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出系统，撤销令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} errors.SuccessResponse "登出成功"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从中间件获取令牌
	token := c.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		tokenString := token[7:] // 去掉 "Bearer " 前缀

		// 撤销令牌
		if err := h.jwtService.RevokeToken(tokenString); err != nil {
			logger.Warn("Failed to revoke token", zap.Error(err))
		}
	}

	logger.Info("User logged out")
	errors.RespondWithSuccess(c, gin.H{"message": "登出成功"}, "登出成功")
}

// GetProfile 获取用户资料
// @Summary 获取当前用户资料
// @Description 获取当前登录用户的详细资料
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} UserInfo "用户资料"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从JWT中间件获取用户信息
	claims, exists := c.Get("user")
	if !exists {
		errors.RespondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "未找到用户信息")
		return
	}

	userClaims, ok := claims.(*shared.Claims)
	if !ok {
		errors.RespondWithError(c, http.StatusInternalServerError, "INVALID_USER_CLAIMS", "用户信息格式错误")
		return
	}

	// 获取完整的用户信息
	userResp, err := h.userService.GetUser(c.Request.Context(), userClaims.UserID)
	if err != nil {
		logger.Error("Failed to get user profile",
			zap.String("user_id", userClaims.UserID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "PROFILE_FETCH_FAILED", "获取用户资料失败")
		return
	}

	profile := &UserInfo{
		ID:     userResp.ID,
		Name:   userResp.Name,
		Email:  userResp.Email,
		Phone:  userResp.Phone,
		Roles:  userClaims.Roles,
		Status: userResp.Status,
	}

	errors.RespondWithSuccess(c, profile, "获取用户资料成功")
}

// 辅助函数：检查是否是邮箱已存在错误
func isEmailExistsError(err error) bool {
	// 这里可以根据具体的错误类型或错误消息来判断
	// 简单实现：检查错误消息中是否包含邮箱相关的关键词
	errMsg := err.Error()
	return contains(errMsg, "email") && (contains(errMsg, "exists") || contains(errMsg, "duplicate") || contains(errMsg, "unique"))
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > len(substr) && (str[:len(substr)] == substr || str[len(str)-len(substr):] == substr || containsSubstring(str, substr)))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
