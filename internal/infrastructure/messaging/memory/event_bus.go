package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// InMemoryEventBus 内存事件总线实现
type InMemoryEventBus struct {
	handlers   map[string][]event.EventHandler
	mu         sync.RWMutex
	eventStore event.EventStore
	running    bool
	stopChan   chan struct{}
	eventChan  chan event.DomainEvent
	bufferSize int
	maxRetries int
	retryDelay time.Duration
}

// EventBusConfig 事件总线配置
type EventBusConfig struct {
	BufferSize int
	MaxRetries int
	RetryDelay time.Duration
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus(config EventBusConfig, eventStore event.EventStore) *InMemoryEventBus {
	if config.BufferSize <= 0 {
		config.BufferSize = 1000
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = time.Second
	}

	return &InMemoryEventBus{
		handlers:   make(map[string][]event.EventHandler),
		eventStore: eventStore,
		stopChan:   make(chan struct{}),
		eventChan:  make(chan event.DomainEvent, config.BufferSize),
		bufferSize: config.BufferSize,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// Start 启动事件总线
func (bus *InMemoryEventBus) Start() error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.running {
		return fmt.Errorf("event bus is already running")
	}
	bus.running = true
	go bus.processEvents()

	logger.Info("InMemoryEventBus started")
	return nil
}

// Stop 停止事件总线
func (bus *InMemoryEventBus) Stop() error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if !bus.running {
		return fmt.Errorf("event bus is not running")
	}

	bus.running = false

	// 安全关闭stopChan，避免重复关闭
	select {
	case <-bus.stopChan:
		// 已经关闭
	default:
		close(bus.stopChan)
	}

	// 等待处理完剩余事件
	time.Sleep(100 * time.Millisecond)

	logger.Info("InMemoryEventBus stopped")
	return nil
}

// Publish 发布单个事件
func (bus *InMemoryEventBus) Publish(event event.DomainEvent) error {
	bus.mu.RLock()
	running := bus.running
	bus.mu.RUnlock()

	if !running {
		return fmt.Errorf("event bus is not running")
	}

	select {
	case bus.eventChan <- event:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout publishing event: %s", event.EventType())
	}
}

// PublishBatch 批量发布事件
func (bus *InMemoryEventBus) PublishBatch(events []event.DomainEvent) error {
	for _, event := range events {
		if err := bus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventID(), err)
		}
	}
	return nil
}

// Subscribe 订阅事件
func (bus *InMemoryEventBus) Subscribe(eventType string, handler event.EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = make([]event.EventHandler, 0)
	}

	// 检查是否已经订阅（使用类型和地址比较）
	for _, h := range bus.handlers[eventType] {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			return fmt.Errorf("handler already subscribed to event type: %s", eventType)
		}
	}

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	logger.Info("Handler subscribed to event type", zap.String("event_type", eventType))
	return nil
}

// Unsubscribe 取消订阅事件
func (bus *InMemoryEventBus) Unsubscribe(eventType string, handler event.EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers, exists := bus.handlers[eventType]
	if !exists {
		return fmt.Errorf("no handlers found for event type: %s", eventType)
	}

	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			logger.Info("Handler unsubscribed from event type", zap.String("event_type", eventType))
			return nil
		}
	}

	return fmt.Errorf("handler not found for event type: %s", eventType)
}

// processEvents 处理事件的主循环
func (bus *InMemoryEventBus) processEvents() {
	for {
		select {
		case event := <-bus.eventChan:
			bus.handleEvent(event)
		case <-bus.stopChan:
			// 安全处理剩余事件，避免竞态条件
			bus.drainRemainingEvents()
			return
		}
	}
}

// drainRemainingEvents 安全地处理剩余事件
func (bus *InMemoryEventBus) drainRemainingEvents() {
	for {
		select {
		case event := <-bus.eventChan:
			bus.handleEvent(event)
		default:
			// 没有更多事件，退出
			return
		}
	}
}

// handleEvent 处理单个事件
func (bus *InMemoryEventBus) handleEvent(domainEvent event.DomainEvent) {
	// 先保存事件到存储
	if bus.eventStore != nil {
		if err := bus.eventStore.Save(domainEvent); err != nil {
			logger.Error("Failed to save event",
				zap.String("event_id", domainEvent.EventID()),
				zap.Error(err))
			// 事件保存失败时，记录错误但继续处理
			// 在生产环境中可能需要更严格的错误处理策略
		}
	}

	// 获取事件处理器
	bus.mu.RLock()
	handlers := bus.handlers[domainEvent.EventType()]
	bus.mu.RUnlock()

	if len(handlers) == 0 {
		logger.Warn("No handlers found for event type",
			zap.String("event_type", domainEvent.EventType()),
			zap.String("event_id", domainEvent.EventID()))
		return
	}

	// 过滤能处理该事件的处理器
	validHandlers := make([]event.EventHandler, 0, len(handlers))
	for _, handler := range handlers {
		if handler.CanHandle(domainEvent.EventType()) {
			validHandlers = append(validHandlers, handler)
		}
	}

	if len(validHandlers) == 0 {
		logger.Warn("No valid handlers can process event type",
			zap.String("event_type", domainEvent.EventType()),
			zap.String("event_id", domainEvent.EventID()))
		return
	}

	// 并发处理事件
	var wg sync.WaitGroup
	for _, handler := range validHandlers {
		wg.Add(1)
		go func(h event.EventHandler) {
			defer wg.Done()
			bus.handleEventWithRetry(domainEvent, h)
		}(handler)
	}
	wg.Wait()
}

// handleEventWithRetry 带重试的事件处理
func (bus *InMemoryEventBus) handleEventWithRetry(domainEvent event.DomainEvent, handler event.EventHandler) {
	var err error
	handlerType := fmt.Sprintf("%T", handler)

	for i := 0; i <= bus.maxRetries; i++ {
		err = handler.Handle(domainEvent)
		if err == nil {
			logger.Info("Event handled successfully",
				zap.String("event_id", domainEvent.EventID()),
				zap.String("event_type", domainEvent.EventType()),
				zap.String("handler_type", handlerType),
				zap.Int("attempt", i+1))
			return
		}

		if i < bus.maxRetries {
			logger.Warn("Event handling failed, retrying",
				zap.String("event_id", domainEvent.EventID()),
				zap.String("event_type", domainEvent.EventType()),
				zap.String("handler_type", handlerType),
				zap.Int("attempt", i+1),
				zap.Int("max_retries", bus.maxRetries),
				zap.Error(err))

			// 指数退避重试延迟
			retryDelay := time.Duration(i+1) * bus.retryDelay
			time.Sleep(retryDelay)
		}
	}

	logger.Error("Event handling failed after all retries",
		zap.String("event_id", domainEvent.EventID()),
		zap.String("event_type", domainEvent.EventType()),
		zap.String("handler_type", handlerType),
		zap.Int("total_attempts", bus.maxRetries+1),
		zap.Error(err))

	// 可以在这里添加死信队列或告警机制
	// TODO: 实现死信队列处理失败事件
}

// GetStats 获取事件总线统计信息
func (bus *InMemoryEventBus) GetStats() EventBusStats {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	stats := EventBusStats{
		Running:       bus.running,
		BufferSize:    bus.bufferSize,
		PendingEvents: len(bus.eventChan),
		HandlerCount:  0,
	}

	for eventType, handlers := range bus.handlers {
		stats.HandlerCount += len(handlers)
		if stats.EventTypes == nil {
			stats.EventTypes = make(map[string]int)
		}
		stats.EventTypes[eventType] = len(handlers)
	}

	return stats
}

// EventBusStats 事件总线统计信息
type EventBusStats struct {
	Running       bool           `json:"running"`
	BufferSize    int            `json:"buffer_size"`
	PendingEvents int            `json:"pending_events"`
	HandlerCount  int            `json:"handler_count"`
	EventTypes    map[string]int `json:"event_types"`
}

// AsyncEventHandler 异步事件处理器包装器
type AsyncEventHandler struct {
	handler event.EventHandler
	timeout time.Duration
}

// NewAsyncEventHandler 创建异步事件处理器
func NewAsyncEventHandler(handler event.EventHandler, timeout time.Duration) *AsyncEventHandler {
	return &AsyncEventHandler{
		handler: handler,
		timeout: timeout,
	}
}

// Handle 异步处理事件
func (h *AsyncEventHandler) Handle(event event.DomainEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- h.handler.Handle(event)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("event handling timeout: %s", event.EventType())
	}
}

// CanHandle 检查是否可以处理指定事件类型
func (h *AsyncEventHandler) CanHandle(eventType string) bool {
	return h.handler.CanHandle(eventType)
}

// EventTypes 返回支持的事件类型
func (h *AsyncEventHandler) EventTypes() []string {
	return h.handler.EventTypes()
}
