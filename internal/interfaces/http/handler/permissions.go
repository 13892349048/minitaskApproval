package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/internal/domain/auth/service"
	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/taskflow/pkg/errors"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// PermissionHandler 权限管理处理器
type PermissionHandler struct {
	permissionService service.PermissionDomainService
}

// NewPermissionHandler 创建权限管理处理器
func NewPermissionHandler(permissionService service.PermissionDomainService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	UserID      string                 `json:"user_id" binding:"required"`
	Resource    string                 `json:"resource" binding:"required"`
	Action      string                 `json:"action" binding:"required"`
	ResourceCtx map[string]interface{} `json:"resource_context,omitempty"`
}

// CheckPermissionResponse 权限检查响应
type CheckPermissionResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

// RevokeRoleRequest 撤销角色请求
type RevokeRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

// CheckPermission 检查用户权限
// @Summary 检查用户权限
// @Description 检查用户是否有执行特定操作的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CheckPermissionRequest true "权限检查信息"
// @Success 200 {object} CheckPermissionResponse "权限检查结果"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/permissions/check [post]
func (h *PermissionHandler) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 执行权限检查
	allowed, err := h.permissionService.CanUserPerformAction(
		c.Request.Context(),
		req.UserID,
		valueobject.ResourceType(req.Resource),
		valueobject.ActionType(req.Action),
		req.ResourceCtx,
	)
	if err != nil {
		logger.Error("Permission check failed",
			zap.String("user_id", req.UserID),
			zap.String("resource", req.Resource),
			zap.String("action", req.Action),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "PERMISSION_CHECK_FAILED", "权限检查失败")
		return
	}

	response := &CheckPermissionResponse{
		Allowed: allowed,
	}

	if !allowed {
		response.Reason = "用户没有执行此操作的权限"
	}

	errors.RespondWithSuccess(c, response, "权限检查完成")
}

// AssignRole 为用户分配角色
// @Summary 为用户分配角色
// @Description 为指定用户分配角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body AssignRoleRequest true "角色分配信息"
// @Success 200 {object} errors.SuccessResponse "分配成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 409 {object} errors.ErrorResponse "角色已分配"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/permissions/assign-role [post]
func (h *PermissionHandler) AssignRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 分配角色
	err := h.permissionService.AssignRoleToUser(
		c.Request.Context(),
		req.UserID,
		valueobject.RoleID(req.RoleID),
	)
	if err != nil {
		logger.Error("Role assignment failed",
			zap.String("user_id", req.UserID),
			zap.String("role_id", req.RoleID),
			zap.Error(err))

		// 检查是否是角色已分配错误
		if isRoleAlreadyAssignedError(err) {
			errors.RespondWithError(c, http.StatusConflict, "ROLE_ALREADY_ASSIGNED", "角色已分配给用户")
			return
		}

		errors.RespondWithError(c, http.StatusInternalServerError, "ROLE_ASSIGNMENT_FAILED", "角色分配失败")
		return
	}

	logger.Info("Role assigned successfully",
		zap.String("user_id", req.UserID),
		zap.String("role_id", req.RoleID))

	errors.RespondWithSuccess(c, gin.H{"message": "角色分配成功"}, "分配成功")
}

// RevokeRole 撤销用户角色
// @Summary 撤销用户角色
// @Description 撤销指定用户的角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body RevokeRoleRequest true "角色撤销信息"
// @Success 200 {object} errors.SuccessResponse "撤销成功"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 404 {object} errors.ErrorResponse "角色未分配"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/permissions/revoke-role [post]
func (h *PermissionHandler) RevokeRole(c *gin.Context) {
	var req RevokeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "请求参数错误: "+err.Error())
		return
	}

	// 撤销角色
	err := h.permissionService.RevokeRoleFromUser(
		c.Request.Context(),
		req.UserID,
		valueobject.RoleID(req.RoleID),
	)
	if err != nil {
		logger.Error("Role revocation failed",
			zap.String("user_id", req.UserID),
			zap.String("role_id", req.RoleID),
			zap.Error(err))

		// 检查是否是角色未分配错误
		if isRoleNotAssignedError(err) {
			errors.RespondWithError(c, http.StatusNotFound, "ROLE_NOT_ASSIGNED", "角色未分配给用户")
			return
		}

		errors.RespondWithError(c, http.StatusInternalServerError, "ROLE_REVOCATION_FAILED", "角色撤销失败")
		return
	}

	logger.Info("Role revoked successfully",
		zap.String("user_id", req.UserID),
		zap.String("role_id", req.RoleID))

	errors.RespondWithSuccess(c, gin.H{"message": "角色撤销成功"}, "撤销成功")
}

// GetUserPermissions 获取用户权限
// @Summary 获取用户权限
// @Description 获取指定用户的所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path string true "用户ID"
// @Success 200 {array} aggregate.Permission "用户权限列表"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/permissions/users/{user_id}/permissions [get]
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_USER_ID", "用户ID不能为空")
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user permissions",
			zap.String("user_id", userID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "GET_PERMISSIONS_FAILED", "获取用户权限失败")
		return
	}

	errors.RespondWithSuccess(c, permissions, "获取用户权限成功")
}

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取指定用户的所有角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path string true "用户ID"
// @Success 200 {array} aggregate.Role "用户角色列表"
// @Failure 400 {object} errors.ErrorResponse "请求参数错误"
// @Failure 401 {object} errors.ErrorResponse "未认证"
// @Failure 500 {object} errors.ErrorResponse "服务器内部错误"
// @Router /api/v1/permissions/users/{user_id}/roles [get]
func (h *PermissionHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		errors.RespondWithError(c, http.StatusBadRequest, "INVALID_USER_ID", "用户ID不能为空")
		return
	}

	roles, err := h.permissionService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user roles",
			zap.String("user_id", userID),
			zap.Error(err))
		errors.RespondWithError(c, http.StatusInternalServerError, "GET_ROLES_FAILED", "获取用户角色失败")
		return
	}

	errors.RespondWithSuccess(c, roles, "获取用户角色成功")
}

// 辅助函数：检查是否是角色已分配错误
func isRoleAlreadyAssignedError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "already") && strings.Contains(errMsg, "assigned")
}

// 辅助函数：检查是否是角色未分配错误
func isRoleNotAssignedError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "not") && strings.Contains(errMsg, "assigned")
}
