package event

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
// 统一的事件接口，所有领域事件必须实现此接口
type DomainEvent interface {
	// EventID 事件唯一标识
	EventID() string

	// EventType 事件类型（如：user.created, project.updated）
	EventType() string

	// AggregateID 聚合根ID
	AggregateID() string

	// AggregateType 聚合根类型
	AggregateType() string

	// OccurredAt 事件发生时间
	OccurredAt() time.Time

	// EventData 事件数据
	EventData() interface{}

	// Version 事件版本
	Version() int
}

// BaseEvent 基础事件结构
// 提供事件的通用字段和方法实现
type BaseEvent struct {
	ID                string    `json:"event_id"`
	Type              string    `json:"event_type"`
	AggregateRootID   string    `json:"aggregate_id"`
	AggregateRootType string    `json:"aggregate_type"`
	Timestamp         time.Time `json:"occurred_at"`
	EventVersion      int       `json:"version"`
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType, aggregateID, aggregateType string) *BaseEvent {
	return &BaseEvent{
		ID:                GenerateEventID(),
		Type:              eventType,
		AggregateRootID:   aggregateID,
		AggregateRootType: aggregateType,
		Timestamp:         time.Now(),
		EventVersion:      1,
	}
}

// EventID 实现 DomainEvent 接口
func (e BaseEvent) EventID() string {
	return e.ID
}

// EventType 实现 DomainEvent 接口
func (e BaseEvent) EventType() string {
	return e.Type
}

// AggregateID 实现 DomainEvent 接口
func (e BaseEvent) AggregateID() string {
	return e.AggregateRootID
}

// AggregateType 实现 DomainEvent 接口
func (e BaseEvent) AggregateType() string {
	return e.AggregateRootType
}

// OccurredAt 实现 DomainEvent 接口
func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// Version 实现 DomainEvent 接口
func (e BaseEvent) Version() int {
	return e.EventVersion
}

// EventData 需要由具体事件实现
// 基础事件返回nil，具体事件应该重写此方法
func (e BaseEvent) EventData() interface{} {
	return nil
}

// GenerateEventID 生成事件ID
func GenerateEventID() string {
	return uuid.New().String()
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(event DomainEvent) error

	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler) error

	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(event DomainEvent) error

	// CanHandle 判断是否能处理该事件
	CanHandle(eventType string) bool

	// EventTypes 返回支持的事件类型列表
	EventTypes() []string
}

// EventStore 事件存储接口
type EventStore interface {
	// Save 保存事件
	Save(event DomainEvent) error

	// GetEvents 获取聚合的所有事件
	GetEvents(aggregateID string, fromVersion int) ([]DomainEvent, error)

	// GetEventsByType 根据类型获取事件
	GetEventsByType(eventType string, limit int) ([]DomainEvent, error)
}
