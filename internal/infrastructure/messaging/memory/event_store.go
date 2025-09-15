package memory

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/taskflow/internal/domain/event"
)

// InMemoryEventStore 内存事件存储实现
type InMemoryEventStore struct {
	events     []event.DomainEvent
	eventsByID map[string]event.DomainEvent
	mu         sync.RWMutex
	maxEvents  int
}

// NewInMemoryEventStore 创建内存事件存储
func NewInMemoryEventStore(maxEvents int) *InMemoryEventStore {
	if maxEvents <= 0 {
		maxEvents = 10000 // 默认最大事件数
	}

	return &InMemoryEventStore{
		events:     make([]event.DomainEvent, 0),
		eventsByID: make(map[string]event.DomainEvent),
		maxEvents:  maxEvents,
	}
}

// Save 保存单个事件
func (store *InMemoryEventStore) Save(event event.DomainEvent) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	// 检查事件是否已存在
	if _, exists := store.eventsByID[event.EventID()]; exists {
		return fmt.Errorf("event with ID %s already exists", event.EventID())
	}

	// 添加事件
	store.events = append(store.events, event)
	store.eventsByID[event.EventID()] = event

	// 检查是否超过最大事件数，如果是则删除最旧的事件
	if len(store.events) > store.maxEvents {
		oldestEvent := store.events[0]
		store.events = store.events[1:]
		delete(store.eventsByID, oldestEvent.EventID())
	}

	return nil
}

// SaveBatch 批量保存事件
func (store *InMemoryEventStore) SaveBatch(events []event.DomainEvent) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, event := range events {
		// 检查事件是否已存在
		if _, exists := store.eventsByID[event.EventID()]; exists {
			return fmt.Errorf("event with ID %s already exists", event.EventID())
		}

		// 添加事件
		store.events = append(store.events, event)
		store.eventsByID[event.EventID()] = event
	}

	// 检查是否超过最大事件数
	if len(store.events) > store.maxEvents {
		excessCount := len(store.events) - store.maxEvents
		for i := 0; i < excessCount; i++ {
			oldestEvent := store.events[i]
			delete(store.eventsByID, oldestEvent.EventID())
		}
		store.events = store.events[excessCount:]
	}

	return nil
}

// GetEvents 获取聚合根的事件
func (store *InMemoryEventStore) GetEvents(aggregateID string, fromVersion int) ([]event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var result []event.DomainEvent
	for _, event := range store.events {
		if event.AggregateID() == aggregateID && event.Version() >= fromVersion {
			result = append(result, event)
		}
	}

	// 按版本排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Version() < result[j].Version()
	})

	return result, nil
}

// GetEventsByType 根据事件类型获取事件
func (store *InMemoryEventStore) GetEventsByType(eventType string, limit int) ([]event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var result []event.DomainEvent
	count := 0

	// 从最新的事件开始查找
	for i := len(store.events) - 1; i >= 0 && count < limit; i-- {
		event := store.events[i]
		if event.EventType() == eventType {
			result = append(result, event)
			count++
		}
	}

	// 反转结果，使其按时间顺序排列
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// GetEventByID 根据事件ID获取事件
func (store *InMemoryEventStore) GetEventByID(eventID string) (event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	event, exists := store.eventsByID[eventID]
	if !exists {
		return nil, fmt.Errorf("event with ID %s not found", eventID)
	}

	return event, nil
}

// GetAllEvents 获取所有事件
func (store *InMemoryEventStore) GetAllEvents(limit int, offset int) ([]event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	total := len(store.events)
	if offset >= total {
		return []event.DomainEvent{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	// 返回指定范围的事件（按时间倒序）
	result := make([]event.DomainEvent, 0, end-offset)
	for i := total - 1 - offset; i >= total-end; i-- {
		result = append(result, store.events[i])
	}

	return result, nil
}

// GetEventsByAggregateType 根据聚合类型获取事件
func (store *InMemoryEventStore) GetEventsByAggregateType(aggregateType string, limit int) ([]event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var result []event.DomainEvent
	count := 0

	for i := len(store.events) - 1; i >= 0 && count < limit; i-- {
		event := store.events[i]
		if event.AggregateType() == aggregateType {
			result = append(result, event)
			count++
		}
	}

	return result, nil
}

// GetEventsByTimeRange 根据时间范围获取事件
func (store *InMemoryEventStore) GetEventsByTimeRange(start, end time.Time, limit int) ([]event.DomainEvent, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	var result []event.DomainEvent
	count := 0

	for i := len(store.events) - 1; i >= 0 && count < limit; i-- {
		event := store.events[i]
		eventTime := event.OccurredAt()
		if (eventTime.After(start) || eventTime.Equal(start)) &&
			(eventTime.Before(end) || eventTime.Equal(end)) {
			result = append(result, event)
			count++
		}
	}

	return result, nil
}

// Clear 清空所有事件
func (store *InMemoryEventStore) Clear() error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.events = make([]event.DomainEvent, 0)
	store.eventsByID = make(map[string]event.DomainEvent)

	return nil
}

// Count 获取事件总数
func (store *InMemoryEventStore) Count() int {
	store.mu.RLock()
	defer store.mu.RUnlock()

	return len(store.events)
}

// GetStats 获取存储统计信息
func (store *InMemoryEventStore) GetStats() EventStoreStats {
	store.mu.RLock()
	defer store.mu.RUnlock()

	stats := EventStoreStats{
		TotalEvents:    len(store.events),
		MaxEvents:      store.maxEvents,
		EventTypes:     make(map[string]int),
		AggregateTypes: make(map[string]int),
	}

	for _, event := range store.events {
		stats.EventTypes[event.EventType()]++
		stats.AggregateTypes[event.AggregateType()]++
	}

	if len(store.events) > 0 {
		stats.OldestEvent = store.events[0].OccurredAt()
		stats.NewestEvent = store.events[len(store.events)-1].OccurredAt()
	}

	return stats
}

// EventStoreStats 事件存储统计信息
type EventStoreStats struct {
	TotalEvents    int            `json:"total_events"`
	MaxEvents      int            `json:"max_events"`
	EventTypes     map[string]int `json:"event_types"`
	AggregateTypes map[string]int `json:"aggregate_types"`
	OldestEvent    time.Time      `json:"oldest_event"`
	NewestEvent    time.Time      `json:"newest_event"`
}
