package handlers

import (
	"fmt"
	"reflect"

	"github.com/taskflow/internal/domain/event"
)

// SafeEventCast 安全的事件类型转换
func SafeEventCast[T any](domainEvent event.DomainEvent, expectedType string) (*T, error) {
	if domainEvent.EventType() != expectedType {
		return nil, fmt.Errorf("event type mismatch: expected %s, got %s", expectedType, domainEvent.EventType())
	}

	// 方法1：直接类型断言具体事件类型
	if concreteEvent, ok := domainEvent.(T); ok {
		return &concreteEvent, nil
	}

	// 方法2：通过EventData()获取
	eventData := domainEvent.EventData()
	if eventData == nil {
		return nil, fmt.Errorf("event data is nil for event type %s", expectedType)
	}

	if typedData, ok := eventData.(*T); ok {
		return typedData, nil
	}

	// 方法3：使用反射进行类型检查和转换
	eventValue := reflect.ValueOf(domainEvent)
	if eventValue.Kind() == reflect.Ptr {
		eventValue = eventValue.Elem()
	}

	var zero T
	expectedValue := reflect.ValueOf(&zero).Elem()

	if eventValue.Type().AssignableTo(expectedValue.Type()) {
		result := eventValue.Interface()
		if typedResult, ok := result.(T); ok {
			return &typedResult, nil
		}
	}

	return nil, fmt.Errorf("failed to cast event to type %T for event type %s", zero, expectedType)
}
