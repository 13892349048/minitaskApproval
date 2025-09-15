package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// Import event types for type casting
type (
	TaskCreatedEvent   = event.TaskCreatedEvent
	TaskCompletedEvent = event.TaskCompletedEvent
	TaskRejectedEvent  = event.TaskRejectedEvent
)

// FixedStatisticsHandler 修复后的统计事件处理器
type FixedStatisticsHandler struct {
	statsRepo StatisticsRepository
	cache     map[string]*TaskStatistics
	mu        sync.RWMutex
}

// TaskStatistics 任务统计数据
type TaskStatistics struct {
	ProjectID         string    `json:"project_id"`
	TotalTasks        int       `json:"total_tasks"`
	CompletedTasks    int       `json:"completed_tasks"`
	PendingTasks      int       `json:"pending_tasks"`
	RejectedTasks     int       `json:"rejected_tasks"`
	AverageCompletion float64   `json:"average_completion_days"`
	LastUpdated       time.Time `json:"last_updated"`
}

// StatisticsRepository 统计仓储接口
type StatisticsRepository interface {
	SaveProjectStats(projectID string, stats *TaskStatistics) error
	GetProjectStats(projectID string) (*TaskStatistics, error)
	UpdateTaskCount(projectID string, increment int) error
	UpdateCompletedCount(projectID string, increment int) error
}

// NewFixedStatisticsHandler 创建修复后的统计处理器
func NewStatisticsHandler(statsRepo StatisticsRepository) *FixedStatisticsHandler {
	return &FixedStatisticsHandler{
		statsRepo: statsRepo,
		cache:     make(map[string]*TaskStatistics),
	}
}

// Handle 处理事件
func (h *FixedStatisticsHandler) Handle(domainevent event.DomainEvent) error {
	switch domainevent.EventType() {
	case "TaskCreated":
		return h.handleTaskCreatedSafe(domainevent)
	case "TaskCompleted":
		return h.handleTaskCompletedSafe(domainevent)
	case "TaskRejected":
		return h.handleTaskRejectedSafe(domainevent)
	default:
		return nil
	}
}

// CanHandle 检查是否可以处理指定事件类型
func (h *FixedStatisticsHandler) CanHandle(eventType string) bool {
	supportedTypes := []string{
		"TaskCreated",
		"TaskCompleted",
		"TaskRejected",
	}

	for _, supportedType := range supportedTypes {
		if eventType == supportedType {
			return true
		}
	}
	return false
}

// EventTypes 返回支持的事件类型
func (h *FixedStatisticsHandler) EventTypes() []string {
	return []string{
		"TaskCreated",
		"TaskCompleted",
		"TaskRejected",
	}
}

// handleTaskCreatedSafe 安全处理TaskCreated事件
func (h *FixedStatisticsHandler) handleTaskCreatedSafe(event event.DomainEvent) error {
	data, err := SafeEventCast[TaskCreatedEvent](event, "TaskCreated")
	if err != nil {
		logger.Error("Failed to cast TaskCreatedEvent", zap.Error(err))
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 更新缓存
	stats := h.getOrCreateStats(data.ProjectID)
	stats.TotalTasks++
	stats.PendingTasks++
	stats.LastUpdated = time.Now()

	// 更新仓储
	if err := h.statsRepo.UpdateTaskCount(data.ProjectID, 1); err != nil {
		log.Printf("Failed to update task count: %v", err)
	}

	log.Printf("Statistics updated for task created in project: %s", data.ProjectID)
	return nil
}

// handleTaskCompletedSafe 安全处理TaskCompleted事件
func (h *FixedStatisticsHandler) handleTaskCompletedSafe(event event.DomainEvent) error {
	data, err := SafeEventCast[TaskCompletedEvent](event, "TaskCompleted")
	if err != nil {
		logger.Error("Failed to cast TaskCompletedEvent", zap.Error(err))
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 这里需要获取任务的项目ID，简化处理
	projectID := "unknown" // 实际应该从任务聚合根获取

	stats := h.getOrCreateStats(projectID)
	stats.CompletedTasks++
	if stats.PendingTasks > 0 {
		stats.PendingTasks--
	}
	stats.LastUpdated = time.Now()

	// 更新仓储
	if err := h.statsRepo.UpdateCompletedCount(projectID, 1); err != nil {
		log.Printf("Failed to update completed count: %v", err)
	}

	log.Printf("Statistics updated for task completed: %s", data.TaskID)
	return nil
}

// handleTaskRejectedSafe 安全处理TaskRejected事件
func (h *FixedStatisticsHandler) handleTaskRejectedSafe(event event.DomainEvent) error {
	data, err := SafeEventCast[TaskRejectedEvent](event, "TaskRejected")
	if err != nil {
		logger.Error("Failed to cast TaskRejectedEvent", zap.Error(err))
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 这里需要获取任务的项目ID，简化处理
	projectID := "unknown"

	stats := h.getOrCreateStats(projectID)
	stats.RejectedTasks++
	stats.LastUpdated = time.Now()

	log.Printf("Statistics updated for task rejected: %s", data.TaskID)
	return nil
}

func (h *FixedStatisticsHandler) getOrCreateStats(projectID string) *TaskStatistics {
	if stats, exists := h.cache[projectID]; exists {
		return stats
	}

	// 尝试从仓储加载
	if stats, err := h.statsRepo.GetProjectStats(projectID); err == nil {
		h.cache[projectID] = stats
		return stats
	}

	// 创建新的统计数据
	stats := &TaskStatistics{
		ProjectID:   projectID,
		LastUpdated: time.Now(),
	}
	h.cache[projectID] = stats
	return stats
}

// GetProjectStatistics 获取项目统计数据
func (h *FixedStatisticsHandler) GetProjectStatistics(projectID string) *TaskStatistics {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.getOrCreateStats(projectID)
}

// FlushCache 刷新缓存到仓储
func (h *FixedStatisticsHandler) FlushCache() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for projectID, stats := range h.cache {
		if err := h.statsRepo.SaveProjectStats(projectID, stats); err != nil {
			log.Printf("Failed to flush stats for project %s: %v", projectID, err)
			return err
		}
	}

	logger.Info("Statistics cache flushed to repository")
	return nil
}
