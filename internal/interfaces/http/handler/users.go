package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/internal/application/service"
	"github.com/taskflow/pkg/errors"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserAppService
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []*service.UserResponse `json:"users"`
	Pagination PaginationInfo          `json:"pagination"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserAppService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取系统中的用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页大小" default(10)
// @Param role query string false "角色过滤"
// @Param status query string false "状态过滤"
// @Success 200 {object} UserListResponse "用户列表"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 解析查询参数
	page := 1
	size := 10
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if s := c.Query("size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 && parsed <= 100 {
			size = parsed
		}
	}

	role := c.Query("role")
	status := c.Query("status")

	// 调用应用服务获取用户列表
	users, total, err := h.userService.ListUsers(c.Request.Context(), &service.ListUsersRequest{
		Page:   page,
		Size:   size,
		Role:   role,
		Status: status,
	})
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "LIST_USERS_FAILED", "获取用户列表失败")
		return
	}

	response := &UserListResponse{
		Users: users,
		Pagination: PaginationInfo{
			Page:  page,
			Size:  size,
			Total: total,
		},
	}

	errors.RespondWithSuccess(c, response, "获取用户列表成功")
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} user.UserResponse "用户信息"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 404 {object} errors.ErrorResponse "用户不存在"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_USER_ID", "用户ID不能为空")
		return
	}

	userResp, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user",
			zap.String("user_id", userID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusNotFound, "USER_NOT_FOUND", "用户不存在")
		return
	}

	errors.RespondWithSuccess(c, userResp, "获取用户信息成功")
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新指定用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param request body user.UpdateUserRequest true "更新信息"
// @Success 200 {object} errors.SuccessResponse "更新成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 404 {object} errors.ErrorResponse "用户不存在"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_USER_ID", "用户ID不能为空")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	req.ID = userID
	if err := h.userService.UpdateUserProfile(c.Request.Context(), &req); err != nil {
		logger.Error("Failed to update user",
			zap.String("user_id", userID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "UPDATE_FAILED", "更新用户失败")
		return
	}

	errors.RespondWithSuccess(c, gin.H{"message": "用户信息更新成功"}, "更新成功")
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定的用户（软删除）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} errors.SuccessResponse "删除成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 404 {object} errors.ErrorResponse "用户不存在"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_USER_ID", "用户ID不能为空")
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to delete user",
			zap.String("user_id", userID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "DELETE_FAILED", "删除用户失败")
		return
	}

	errors.RespondWithSuccess(c, gin.H{"message": "用户删除成功"}, "删除成功")
}

// 兼容性函数 - 保持现有路由工作
func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint - to be implemented"})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint - to be implemented"})
}

func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user endpoint - to be implemented"})
}

func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete user endpoint - to be implemented"})
}
