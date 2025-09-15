package events

import (
	"time"

	"github.com/taskflow/internal/application/handlers"
	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/internal/infrastructure/messaging/memory"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// EventBusManager 事件总线管理器
type EventBusManager struct {
	eventBus   *memory.InMemoryEventBus
	eventStore event.EventStore
	handlers   []event.EventHandler
}

// NewEventBusManager 创建事件总线管理器
func NewEventBusManager() *EventBusManager {
	// 创建内存事件存储
	eventStore := memory.NewInMemoryEventStore(10000)

	// 配置事件总线
	config := memory.EventBusConfig{
		BufferSize: 1000,
		MaxRetries: 3,
		RetryDelay: time.Second,
	}

	// 创建事件总线
	eventBus := memory.NewInMemoryEventBus(config, eventStore)

	return &EventBusManager{
		eventBus:   eventBus,
		eventStore: eventStore,
		handlers:   make([]event.EventHandler, 0),
	}
}

// RegisterHandlers 注册事件处理器
func (m *EventBusManager) RegisterHandlers() error {
	// 创建通知处理器
	notificationHandler := handlers.NewNotificationHandler(
		&MockEmailService{},
		&MockSMSService{},
	)

	// 创建审计处理器
	auditHandler := handlers.NewAuditHandler(
		&MockAuditRepository{},
	)

	// 创建统计处理器
	statisticsHandler := handlers.NewStatisticsHandler(
		&MockStatisticsRepository{},
	)

	// 注册处理器
	handlers := []event.EventHandler{
		notificationHandler,
		auditHandler,
		statisticsHandler,
	}

	// 定义事件类型到处理器的映射
	eventTypeMapping := map[string][]event.EventHandler{
		"TaskCreated":              {notificationHandler, auditHandler, statisticsHandler},
		"TaskAssigned":             {notificationHandler, auditHandler},
		"TaskStatusChanged":        {notificationHandler, auditHandler},
		"TaskCompleted":            {notificationHandler, auditHandler, statisticsHandler},
		"TaskRejected":             {notificationHandler, auditHandler, statisticsHandler},
		"ParticipantAdded":         {notificationHandler, auditHandler},
		"ParticipantRemoved":       {notificationHandler, auditHandler},
		"WorkSubmitted":            {notificationHandler, auditHandler},
		"WorkReviewed":             {notificationHandler, auditHandler},
		"TaskCompletionSubmitted":  {notificationHandler, auditHandler},
		"ExtensionRequested":       {notificationHandler, auditHandler},
		"ExtensionApproved":        {notificationHandler, auditHandler},
		"ExtensionRejected":        {notificationHandler, auditHandler},
		"NextExecutionPrepared":    {auditHandler},
		"AllParticipantsCompleted": {auditHandler},
	}

	// 注册事件处理器
	for eventType, eventHandlers := range eventTypeMapping {
		for _, handler := range eventHandlers {
			if handler.CanHandle(eventType) {
				if err := m.eventBus.Subscribe(eventType, handler); err != nil {
					logger.Error("Failed to subscribe handler for event type",
						zap.String("event_type", eventType),
						zap.Error(err))
					return err
				}
			}
		}
	}

	// 保存处理器引用
	m.handlers = handlers

	logger.Info("Registered event handlers", zap.Int("handler_count", len(handlers)))
	return nil
}

// Start 启动事件总线
func (m *EventBusManager) Start() error {
	if err := m.RegisterHandlers(); err != nil {
		return err
	}

	if err := m.eventBus.Start(); err != nil {
		return err
	}

	logger.Info("Event bus started successfully")
	return nil
}

// Stop 停止事件总线
func (m *EventBusManager) Stop() error {
	if err := m.eventBus.Stop(); err != nil {
		return err
	}

	logger.Info("Event bus stopped successfully")
	return nil
}

// GetEventBus 获取事件总线
func (m *EventBusManager) GetEventBus() event.EventBus {
	return m.eventBus
}

// GetEventStore 获取事件存储
func (m *EventBusManager) GetEventStore() event.EventStore {
	return m.eventStore
}
