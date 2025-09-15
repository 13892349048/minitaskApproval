package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/internal/application/service"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectAppService *service.ProjectAppService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(projectAppService *service.ProjectAppService) *ProjectHandler {
	return &ProjectHandler{
		projectAppService: projectAppService,
	}
}

// ListProjects 获取项目列表
// @Summary 获取项目列表
// @Description 分页获取项目列表，支持搜索和过滤
// @Tags projects
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "项目状态" Enums(draft,active,paused,completed,cancelled)
// @Param type query string false "项目类型" Enums(master,sub)
// @Param owner_id query string false "所有者ID"
// @Param manager_id query string false "管理者ID"
// @Param search query string false "搜索关键词"
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" default(desc)
// @Success 200 {object} service.ProjectListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	var req service.ProjectListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	response, err := h.projectAppService.ListProjects(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateProject 创建项目
// @Summary 创建新项目
// @Description 创建新的项目
// @Tags projects
// @Accept json
// @Produce json
// @Param request body service.CreateProjectRequest true "创建项目请求"
// @Success 201 {object} service.ProjectResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.projectAppService.CreateProject(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProject 获取项目详情
// @Summary 获取项目详情
// @Description 根据ID获取项目详细信息
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} service.ProjectResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	response, err := h.projectAppService.GetProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProject 更新项目
// @Summary 更新项目信息
// @Description 更新项目基本信息
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body service.UpdateProjectRequest true "更新项目请求"
// @Success 200 {object} service.ProjectResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = projectID
	err := h.projectAppService.UpdateProject(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的项目信息
	response, err := h.projectAppService.GetProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProject 删除项目
// @Summary 删除项目
// @Description 软删除项目
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	// 从JWT或session中获取用户ID
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.DeleteProject(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProjectMembers 获取项目成员
// @Summary 获取项目成员列表
// @Description 获取指定项目的所有成员
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {array} service.ProjectMemberResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/members [get]
func (h *ProjectHandler) GetProjectMembers(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	project, err := h.projectAppService.GetProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project.Members)
}

// AddProjectMember 添加项目成员
// @Summary 添加项目成员
// @Description 向项目添加新成员
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body service.AddMemberRequest true "添加成员请求"
// @Success 201 {object} service.ProjectMemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/members [post]
func (h *ProjectHandler) AddProjectMember(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	var req service.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT或session中获取操作者ID
	operatorID := c.GetString("user_id")
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.AddMember(c.Request.Context(), projectID, req.UserID, req.Role, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "member added successfully"})
}

// RemoveProjectMember 移除项目成员
// @Summary 移除项目成员
// @Description 从项目中移除成员
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param user_id path string true "用户ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/members/{user_id} [delete]
func (h *ProjectHandler) RemoveProjectMember(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("user_id")
	
	if projectID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID and user ID are required"})
		return
	}

	// 从JWT或session中获取操作者ID
	operatorID := c.GetString("user_id")
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.RemoveMember(c.Request.Context(), projectID, userID, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateMemberRole 更新成员角色
// @Summary 更新项目成员角色
// @Description 更新项目成员的角色
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param user_id path string true "用户ID"
// @Param request body service.UpdateMemberRoleRequest true "更新角色请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/members/{user_id}/role [put]
func (h *ProjectHandler) UpdateMemberRole(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("user_id")
	
	if projectID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID and user ID are required"})
		return
	}

	var req service.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT或session中获取操作者ID
	operatorID := c.GetString("user_id")
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.UpdateMemberRole(c.Request.Context(), projectID, userID, req.Role, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member role updated successfully"})
}

// AssignManager 分配项目管理者
// @Summary 分配项目管理者
// @Description 为项目分配管理者
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body service.AssignManagerRequest true "分配管理者请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/manager [put]
func (h *ProjectHandler) AssignManager(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	var req service.AssignManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT或session中获取操作者ID
	operatorID := c.GetString("user_id")
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.AssignManager(c.Request.Context(), projectID, req.ManagerID, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "manager assigned successfully"})
}

// ChangeProjectStatus 更改项目状态
// @Summary 更改项目状态
// @Description 更改项目状态（激活、暂停、完成、取消）
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body service.ChangeStatusRequest true "状态更改请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/status [put]
func (h *ProjectHandler) ChangeProjectStatus(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	var req service.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT或session中获取操作者ID
	operatorID := c.GetString("user_id")
	if operatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.projectAppService.ChangeStatus(c.Request.Context(), projectID, operatorID, req.Status, req.Reason)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project status updated successfully"})
}

// GetSubProjects 获取子项目
// @Summary 获取子项目列表
// @Description 获取指定项目的所有子项目
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {array} service.ProjectResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/subprojects [get]
func (h *ProjectHandler) GetSubProjects(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	hierarchy, err := h.projectAppService.GetProjectHierarchy(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hierarchy.Children)
}

// CreateSubProject 创建子项目
// @Summary 创建子项目
// @Description 在指定项目下创建子项目
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "父项目ID"
// @Param request body service.CreateSubProjectRequest true "创建子项目请求"
// @Success 201 {object} service.ProjectResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/subprojects [post]
func (h *ProjectHandler) CreateSubProject(c *gin.Context) {
	parentID := c.Param("id")
	if parentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent project ID is required"})
		return
	}

	var req service.CreateSubProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT或session中获取创建者ID
	creatorID := c.GetString("user_id")
	if creatorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	response, err := h.projectAppService.CreateSubProject(c.Request.Context(), parentID, req.Name, req.Description, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProjectHierarchy 获取项目层级结构
// @Summary 获取项目层级结构
// @Description 获取项目的完整层级结构，包括父项目和子项目
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} service.ProjectHierarchyResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/projects/{id}/hierarchy [get]
func (h *ProjectHandler) GetProjectHierarchy(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project ID is required"})
		return
	}

	response, err := h.projectAppService.GetProjectHierarchy(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Legacy functions for backward compatibility
func ListProjects(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.ListProjects instead"})
}

func CreateProject(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.CreateProject instead"})
}

func GetProject(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.GetProject instead"})
}

func UpdateProject(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.UpdateProject instead"})
}

func DeleteProject(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.DeleteProject instead"})
}

func GetProjectMembers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.GetProjectMembers instead"})
}

func AddProjectMember(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.AddProjectMember instead"})
}

func RemoveProjectMember(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.RemoveProjectMember instead"})
}

func GetSubProjects(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.GetSubProjects instead"})
}

func CreateSubProject(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Please use ProjectHandler.CreateSubProject instead"})
}
