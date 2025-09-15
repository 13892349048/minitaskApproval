package events

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/taskflow/internal/domain/event"
)

/*
异步优先的事件架构：

默认异步处理 - 避免阻塞主业务流程
工作池模式 - 多协程并发处理事件
同步选项 - 关键业务可选择同步处理
错误隔离 - 单个处理器失败不影响其他处理器
优雅关闭 - 支持安全停止和资源清理
*/

// EventBusImpl 事件总线实现 - 支持同步和异步发布
type EventBusImpl struct {
	handlers map[string][]event.EventHandler
	mutex    sync.RWMutex

	// 异步处理配置
	workerCount int
	eventQueue  chan event.DomainEvent
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewEventBus 创建事件总线
func NewEventBus(workerCount int) *EventBusImpl {
	ctx, cancel := context.WithCancel(context.Background())

	bus := &EventBusImpl{
		handlers:    make(map[string][]event.EventHandler),
		workerCount: workerCount,
		eventQueue:  make(chan event.DomainEvent, 1000), // 缓冲1000个事件
		ctx:         ctx,
		cancel:      cancel,
	}

	// 启动异步工作协程
	bus.startWorkers()

	return bus
}

// Publish 发布单个事件 - 异步处理
func (b *EventBusImpl) Publish(event event.DomainEvent) error {
	select {
	case b.eventQueue <- event:
		return nil
	case <-b.ctx.Done():
		return fmt.Errorf("event bus is shutting down")
	default:
		return fmt.Errorf("event queue is full")
	}
}

// PublishSync 同步发布事件 - 用于关键业务流程
func (b *EventBusImpl) PublishSync(event event.DomainEvent) error {
	b.mutex.RLock()
	handlers := b.handlers[event.EventType()]
	b.mutex.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			log.Printf("Error handling event %s: %v", event.EventType(), err)
			// 可以根据需要决定是否继续处理其他handler
		}
	}

	return nil
}

// PublishBatch 批量发布事件
func (b *EventBusImpl) PublishBatch(events []event.DomainEvent) error {
	for _, event := range events {
		if err := b.Publish(event); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅事件
func (b *EventBusImpl) Subscribe(eventType string, handler event.EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

// Unsubscribe 取消订阅
func (b *EventBusImpl) Unsubscribe(eventType string, handler event.EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// Start 启动事件总线
func (b *EventBusImpl) Start() error {
	// 事件总线在创建时已经启动
	return nil
}

// Stop 停止事件总线
func (b *EventBusImpl) Stop() error {
	b.cancel()
	close(b.eventQueue)
	return nil
}

// 私有方法

func (b *EventBusImpl) startWorkers() {
	for i := 0; i < b.workerCount; i++ {
		go b.worker()
	}
}

func (b *EventBusImpl) worker() {
	for {
		select {
		case event, ok := <-b.eventQueue:
			if !ok {
				return // 队列已关闭
			}
			b.handleEvent(event)
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *EventBusImpl) handleEvent(event event.DomainEvent) {
	b.mutex.RLock()
	handlers := b.handlers[event.EventType()]
	b.mutex.RUnlock()

	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in event handler for %s: %v", event.EventType(), r)
				}
			}()

			if err := handler.Handle(event); err != nil {
				log.Printf("Error handling event %s: %v", event.EventType(), err)
			}
		}()
	}
}

// 确保实现了接口
var _ event.EventBus = (*EventBusImpl)(nil)
