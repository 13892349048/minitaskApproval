package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/taskflow/internal/application/dto"
	"github.com/taskflow/internal/application/service"
	"github.com/taskflow/internal/domain/valueobject"
	"go.uber.org/zap"
)

// TaskHandler 任务HTTP处理器
type TaskHandler struct {
	taskService *service.TaskAppService
	logger      *zap.Logger
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService *service.TaskAppService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// 验证请求
	if err := h.validateCreateTaskRequest(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// 调用应用服务
	resp, err := h.taskService.CreateTask(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create task", zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create task", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusCreated, resp)
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	resp, err := h.taskService.GetTask(r.Context(), taskID)
	if err != nil {
		h.logger.Error("Failed to get task", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusNotFound, "Task not found", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, resp)
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.ID = taskID
	resp, err := h.taskService.UpdateTask(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to update task", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update task", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, resp)
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	err := h.taskService.DeleteTask(r.Context(), valueobject.TaskID(taskID))
	if err != nil {
		h.logger.Error("Failed to delete task", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete task", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusNoContent, nil)
}

// ListTasks 获取任务列表
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 构建查询条件
	criteria := dto.TaskSearchCriteria{}

	if projectID := query.Get("project_id"); projectID != "" {
		pid := valueobject.ProjectID(projectID)
		criteria.ProjectID = &pid
	}

	if creatorID := query.Get("creator_id"); creatorID != "" {
		cid := valueobject.UserID(creatorID)
		criteria.CreatorID = &cid
	}

	if responsibleID := query.Get("responsible_id"); responsibleID != "" {
		rid := valueobject.UserID(responsibleID)
		criteria.ResponsibleID = &rid
	}

	if status := query.Get("status"); status != "" {
		s := valueobject.TaskStatus(status)
		criteria.Status = &s
	}

	if priority := query.Get("priority"); priority != "" {
		p := valueobject.TaskPriority(priority)
		criteria.Priority = &p
	}

	if taskType := query.Get("type"); taskType != "" {
		tt := valueobject.TaskType(taskType)
		criteria.TaskType = &tt
	}

	if title := query.Get("title"); title != "" {
		criteria.Title = &title
	}

	if description := query.Get("description"); description != "" {
		criteria.Description = &description
	}

	// 分页参数
	page := 1
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := query.Get("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	req := dto.ListTasksRequest{
		Criteria: criteria,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.taskService.ListTasks(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to list tasks", zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to list tasks", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, resp)
}

// AssignTask 分配任务
func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	var req dto.AssignTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.TaskID = taskID
	err := h.taskService.AssignTask(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to assign task", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to assign task", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Task assigned successfully"})
}

// UpdateTaskStatus 更新任务状态
func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	var req dto.UpdateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.TaskID = taskID
	err := h.taskService.UpdateTaskStatus(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to update task status", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update task status", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Task status updated successfully"})
}

// AddTaskParticipant 添加任务参与者
func (h *TaskHandler) AddTaskParticipant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	if taskID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID is required", nil)
		return
	}

	var req dto.AddTaskParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	req.TaskID = taskID
	err := h.taskService.AddTaskParticipant(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to add task participant", zap.String("taskID", taskID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to add task participant", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Participant added successfully"})
}

// RemoveTaskParticipant 移除任务参与者
func (h *TaskHandler) RemoveTaskParticipant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]
	participantID := vars["participant_id"]

	if taskID == "" || participantID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Task ID and participant ID are required", nil)
		return
	}

	req := dto.RemoveTaskParticipantRequest{
		TaskID:        taskID,
		ParticipantID: participantID,
	}

	err := h.taskService.RemoveTaskParticipant(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to remove task participant", zap.String("taskID", taskID), zap.String("participantID", participantID), zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to remove task participant", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, map[string]string{"message": "Participant removed successfully"})
}

// GetTaskStatistics 获取任务统计信息
func (h *TaskHandler) GetTaskStatistics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var projectID *valueobject.ProjectID
	if pid := query.Get("project_id"); pid != "" {
		p := valueobject.ProjectID(pid)
		projectID = &p
	}

	resp, err := h.taskService.GetTaskStatistics(r.Context(), projectID)
	if err != nil {
		h.logger.Error("Failed to get task statistics", zap.Error(err))
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get task statistics", err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, resp)
}

// validateCreateTaskRequest 验证创建任务请求
func (h *TaskHandler) validateCreateTaskRequest(req *dto.CreateTaskRequest) error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}
	if req.CreatorID == "" {
		return fmt.Errorf("creator_id is required")
	}
	if req.ResponsibleID == "" {
		return fmt.Errorf("responsible_id is required")
	}
	if req.Priority == "" {
		return fmt.Errorf("priority is required")
	}
	if req.TaskType == "" {
		return fmt.Errorf("task_type is required")
	}
	return nil
}

// writeSuccessResponse 写入成功响应
func (h *TaskHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}

	if data != nil {
		json.NewEncoder(w).Encode(response)
	}
}

// writeErrorResponse 写入错误响应
func (h *TaskHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"success": false,
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}
