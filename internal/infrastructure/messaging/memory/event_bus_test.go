package memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/pkg/logger"
)

// setupLogger 初始化测试用的logger
func setupLogger(t *testing.T) {
	err := logger.InitLogger(&logger.Config{
		Level:  "info",
		Format: "console",
		Output: "console",
	})
	if err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}
}

// MockEventHandler 模拟事件处理器
type MockEventHandler struct {
	handledEvents []event.DomainEvent
	mu            sync.Mutex
	shouldError   bool
}

func NewMockEventHandler() *MockEventHandler {
	return &MockEventHandler{
		handledEvents: make([]event.DomainEvent, 0),
	}
}

func (h *MockEventHandler) Handle(event event.DomainEvent) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.shouldError {
		return fmt.Errorf("mock error")
	}

	h.handledEvents = append(h.handledEvents, event)
	return nil
}

func (h *MockEventHandler) CanHandle(eventType string) bool {
	for _, t := range h.EventTypes() {
		if t == eventType {
			return true
		}
	}
	return false
}

func (h *MockEventHandler) EventTypes() []string {
	return []string{"TaskCreated", "TaskCompleted"}
}

func (h *MockEventHandler) GetHandledEvents() []event.DomainEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]event.DomainEvent{}, h.handledEvents...)
}

func (h *MockEventHandler) SetShouldError(shouldError bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.shouldError = shouldError
}

func TestInMemoryEventBus_StartStop(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	// 测试启动
	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	// 测试重复启动
	err = bus.Start()
	if err == nil {
		t.Fatal("Expected error when starting already running bus")
	}

	// 测试停止
	err = bus.Stop()
	if err != nil {
		t.Fatalf("Failed to stop event bus: %v", err)
	}

	// 测试重复停止
	err = bus.Stop()
	if err == nil {
		t.Fatal("Expected error when stopping already stopped bus")
	}
}

func TestInMemoryEventBus_PublishSubscribe(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	// 启动事件总线
	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	// 创建模拟处理器
	handler := NewMockEventHandler()

	// 订阅事件
	err = bus.Subscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 创建测试事件
	testEvent := event.NewTaskCreatedEvent(
		"task-1",
		"Test Task",
		"project-1",
		"user-1",
		"user-2",
		"single",
		"normal",
		time.Now().Add(24*time.Hour),
	)

	// 发布事件
	err = bus.Publish(testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待事件处理
	time.Sleep(time.Millisecond * 100)

	// 验证事件被处理
	handledEvents := handler.GetHandledEvents()
	if len(handledEvents) != 1 {
		t.Fatalf("Expected 1 handled event, got %d", len(handledEvents))
	}

	if handledEvents[0].EventType() != "TaskCreated" {
		t.Fatalf("Expected TaskCreated event, got %s", handledEvents[0].EventType())
	}
}

func TestInMemoryEventBus_MultipleHandlers(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	// 创建多个处理器
	handler1 := NewMockEventHandler()
	handler2 := NewMockEventHandler()

	// 订阅相同事件类型
	err = bus.Subscribe("TaskCreated", handler1)
	if err != nil {
		t.Fatalf("Failed to subscribe handler1: %v", err)
	}

	err = bus.Subscribe("TaskCreated", handler2)
	if err != nil {
		t.Fatalf("Failed to subscribe handler2: %v", err)
	}

	// 发布事件
	testEvent := event.NewTaskCreatedEvent(
		"task-1",
		"Test Task",
		"project-1",
		"user-1",
		"user-2",
		"single",
		"normal",
		time.Now().Add(24*time.Hour),
	)

	err = bus.Publish(testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待事件处理
	time.Sleep(time.Millisecond * 100)

	// 验证两个处理器都收到事件
	events1 := handler1.GetHandledEvents()
	events2 := handler2.GetHandledEvents()

	if len(events1) != 1 {
		t.Fatalf("Handler1 expected 1 event, got %d", len(events1))
	}

	if len(events2) != 1 {
		t.Fatalf("Handler2 expected 1 event, got %d", len(events2))
	}
}

func TestInMemoryEventBus_RetryOnError(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 2,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	// 创建会出错的处理器
	handler := NewMockEventHandler()
	handler.SetShouldError(true)

	err = bus.Subscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 发布事件
	testEvent := event.NewTaskCreatedEvent(
		"task-1",
		"Test Task",
		"project-1",
		"user-1",
		"user-2",
		"single",
		"normal",
		time.Now().Add(24*time.Hour),
	)

	err = bus.Publish(testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待重试完成
	time.Sleep(time.Millisecond * 200)

	// 验证事件被尝试处理但失败
	handledEvents := handler.GetHandledEvents()
	if len(handledEvents) != 0 {
		t.Fatalf("Expected 0 handled events due to errors, got %d", len(handledEvents))
	}
}

func TestInMemoryEventBus_BatchPublish(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	handler := NewMockEventHandler()
	err = bus.Subscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 创建多个事件
	events := []event.DomainEvent{
		event.NewTaskCreatedEvent("task-1", "Task 1", "project-1", "user-1", "user-2", "single", "normal", time.Now()),
		event.NewTaskCreatedEvent("task-2", "Task 2", "project-1", "user-1", "user-2", "single", "normal", time.Now()),
		event.NewTaskCreatedEvent("task-3", "Task 3", "project-1", "user-1", "user-2", "single", "normal", time.Now()),
	}

	// 批量发布
	err = bus.PublishBatch(events)
	if err != nil {
		t.Fatalf("Failed to publish batch: %v", err)
	}

	// 等待处理
	time.Sleep(time.Millisecond * 100)

	// 验证所有事件被处理
	handledEvents := handler.GetHandledEvents()
	if len(handledEvents) != 3 {
		t.Fatalf("Expected 3 handled events, got %d", len(handledEvents))
	}
}

func TestInMemoryEventBus_Unsubscribe(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	handler := NewMockEventHandler()

	// 订阅
	err = bus.Subscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 取消订阅
	err = bus.Unsubscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}

	// 发布事件
	testEvent := event.NewTaskCreatedEvent(
		"task-1",
		"Test Task",
		"project-1",
		"user-1",
		"user-2",
		"single",
		"normal",
		time.Now().Add(24*time.Hour),
	)

	err = bus.Publish(testEvent)
	if err != nil {
		t.Fatalf("Failed to publish event: %v", err)
	}

	// 等待处理
	time.Sleep(time.Millisecond * 100)

	// 验证事件未被处理
	handledEvents := handler.GetHandledEvents()
	if len(handledEvents) != 0 {
		t.Fatalf("Expected 0 handled events after unsubscribe, got %d", len(handledEvents))
	}
}

func TestInMemoryEventBus_Stats(t *testing.T) {
	setupLogger(t)

	config := EventBusConfig{
		BufferSize: 10,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	eventStore := NewInMemoryEventStore(100)
	bus := NewInMemoryEventBus(config, eventStore)

	// 获取初始统计
	stats := bus.GetStats()
	if stats.Running {
		t.Fatal("Expected bus to not be running initially")
	}

	// 启动并订阅
	err := bus.Start()
	if err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}
	defer bus.Stop()

	handler := NewMockEventHandler()
	err = bus.Subscribe("TaskCreated", handler)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// 获取更新后的统计
	stats = bus.GetStats()
	if !stats.Running {
		t.Fatal("Expected bus to be running")
	}

	if stats.BufferSize != 10 {
		t.Fatalf("Expected buffer size 10, got %d", stats.BufferSize)
	}

	if stats.HandlerCount != 1 {
		t.Fatalf("Expected 1 handler, got %d", stats.HandlerCount)
	}

	if stats.EventTypes["TaskCreated"] != 1 {
		t.Fatalf("Expected 1 TaskCreated handler, got %d", stats.EventTypes["TaskCreated"])
	}
}
