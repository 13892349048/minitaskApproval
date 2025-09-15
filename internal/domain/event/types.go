package event

import (
	"context"
	"time"
)

// ID 通用ID类型
type ID string

func (id ID) String() string {
	return string(id)
}

func (id ID) IsEmpty() bool {
	return string(id) == ""
}

// AggregateRoot 聚合根接口
type AggregateRoot interface {
	// ID 获取聚合根ID
	ID() string

	// Events 获取未提交的领域事件
	Events() []DomainEvent

	// ClearEvents 清空事件列表
	ClearEvents()

	// Version 获取版本号（用于乐观锁）
	Version() int
}

// BaseAggregate 基础聚合根
type BaseAggregate struct {
	id        string
	version   int
	events    []DomainEvent
	createdAt time.Time
	updatedAt time.Time
}

// NewBaseAggregate 创建基础聚合根
func NewBaseAggregate(id string) BaseAggregate {
	now := time.Now()
	return BaseAggregate{
		id:        id,
		version:   1,
		events:    make([]DomainEvent, 0),
		createdAt: now,
		updatedAt: now,
	}
}

// ID 实现 AggregateRoot 接口
func (a *BaseAggregate) ID() string {
	return a.id
}

// Events 实现 AggregateRoot 接口
func (a *BaseAggregate) Events() []DomainEvent {
	return a.events
}

// ClearEvents 实现 AggregateRoot 接口
func (a *BaseAggregate) ClearEvents() {
	a.events = make([]DomainEvent, 0)
}

// Version 实现 AggregateRoot 接口
func (a *BaseAggregate) Version() int {
	return a.version
}

// CreatedAt 获取创建时间
func (a *BaseAggregate) CreatedAt() time.Time {
	return a.createdAt
}

// UpdatedAt 获取更新时间
func (a *BaseAggregate) UpdatedAt() time.Time {
	return a.updatedAt
}

// AddEvent 添加领域事件
func (a *BaseAggregate) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
	a.updatedAt = time.Now()
}

// IncrementVersion 增加版本号
func (a *BaseAggregate) IncrementVersion() {
	a.version++
	a.updatedAt = time.Now()
}

// Repository 仓储接口
type Repository interface {
	// NextID 生成下一个ID
	NextID() string
}

// Specification 规约接口
type Specification interface {
	// IsSatisfiedBy 检查对象是否满足规约
	IsSatisfiedBy(candidate interface{}) bool
}

// ValueObject 值对象接口
type ValueObject interface {
	// Equals 比较两个值对象是否相等
	Equals(other ValueObject) bool
}

// DomainService 领域服务接口标记
type DomainService interface {
	// 标记接口，用于标识领域服务
}

// TransactionManager 事务管理器接口
type TransactionManager interface {
	// WithTransaction 在事务中执行操作
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	WithTransactionResult(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
}

// Status 通用状态类型
type Status string

// Priority 优先级类型
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// IsValid 检查优先级是否有效
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (p Priority) String() string {
	return string(p)
}
