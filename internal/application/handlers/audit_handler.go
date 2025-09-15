package handlers

import (
	"encoding/json"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/domain/shared"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// AuditHandler 审计事件处理器
type AuditHandler struct {
	auditRepo AuditRepository
}

// AuditLog 审计日志
type AuditLog struct {
	ID            string                 `json:"id"`
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	EventData     map[string]interface{} `json:"event_data"`
	OccurredAt    time.Time              `json:"occurred_at"`
	CreatedAt     time.Time              `json:"created_at"`
}

// AuditRepository 审计仓储接口
type AuditRepository interface {
	Save(log *AuditLog) error
	FindByAggregateID(aggregateID string, limit int) ([]*AuditLog, error)
	FindByEventType(eventType string, limit int) ([]*AuditLog, error)
	FindByTimeRange(start, end time.Time, limit int) ([]*AuditLog, error)
}

// NewAuditHandler 创建审计处理器
func NewAuditHandler(auditRepo AuditRepository) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

// Handle 处理事件
func (h *AuditHandler) Handle(event event.DomainEvent) error {
	// 将事件数据转换为map
	eventDataBytes, err := json.Marshal(event.EventData())
	if err != nil {
		logger.Error("Failed to marshal event data", zap.Error(err))
		return err
	}

	var eventDataMap map[string]interface{}
	if err := json.Unmarshal(eventDataBytes, &eventDataMap); err != nil {
		logger.Error("Failed to unmarshal event data to map", zap.Error(err))
		return err
	}

	// 创建审计日志
	auditLog := &AuditLog{
		ID:            shared.GenerateUUID(),
		EventID:       event.EventID(),
		EventType:     event.EventType(),
		AggregateID:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		EventData:     eventDataMap,
		OccurredAt:    event.OccurredAt(),
		CreatedAt:     time.Now(),
	}

	// 保存审计日志
	if err := h.auditRepo.Save(auditLog); err != nil {
		logger.Error("Failed to save audit log", zap.Error(err))
		return err
	}

	logger.Info("Audit log saved for event", zap.String("event_id", event.EventID()))
	return nil
}

// CanHandle 检查是否可以处理指定事件类型
func (h *AuditHandler) CanHandle(eventType string) bool {
	// 审计处理器处理所有事件类型
	return true
}

// EventTypes 返回支持的事件类型（审计所有事件）
func (h *AuditHandler) EventTypes() []string {
	return []string{
		"TaskCreated",
		"TaskAssigned",
		"TaskStatusChanged",
		"ParticipantAdded",
		"ParticipantRemoved",
		"WorkSubmitted",
		"WorkReviewed",
		"TaskCompletionSubmitted",
		"TaskCompleted",
		"TaskRejected",
		"ExtensionRequested",
		"ExtensionApproved",
		"ExtensionRejected",
		"NextExecutionPrepared",
		"AllParticipantsCompleted",
	}
}
