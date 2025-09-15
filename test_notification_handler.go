package main

import (
	"fmt"
	"time"

	"github.com/taskflow/internal/application/handlers"
	"github.com/taskflow/internal/domain/event"
)

// MockEmailService 模拟邮件服务
type MockEmailService struct{}

func (m *MockEmailService) SendEmail(to, subject, body string) error {
	fmt.Printf("📧 Email sent to: %s\n   Subject: %s\n   Body: %s\n\n", to, subject, body)
	return nil
}

// MockSMSService 模拟短信服务
type MockSMSService struct{}

func (m *MockSMSService) SendSMS(to, message string) error {
	fmt.Printf("📱 SMS sent to: %s\n   Message: %s\n\n", to, message)
	return nil
}

func main() {
	fmt.Println("🧪 Testing NotificationHandler Event Processing")
	fmt.Println("=")

	// 创建模拟服务
	emailService := &MockEmailService{}
	smsService := &MockSMSService{}

	// 创建通知处理器
	handler := handlers.NewNotificationHandler(emailService, smsService)

	testCount := 0
	successCount := 0

	// 测试 TaskCreatedEvent
	testCount++
	fmt.Printf("\n%d. Testing TaskCreatedEvent handling...\n", testCount)
	taskCreatedEvent := event.NewTaskCreatedEvent(
		"task-123",
		"Test Task Creation",
		"project-456",
		"creator-789",
		"responsible-101",
		"development",
		"high",
		time.Now().Add(24*time.Hour),
	)

	if err := handler.Handle(taskCreatedEvent); err != nil {
		fmt.Printf("❌ TaskCreatedEvent handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ TaskCreatedEvent handled successfully\n")
		successCount++
	}

	// 测试 TaskAssignedEvent
	testCount++
	fmt.Printf("\n%d. Testing TaskAssignedEvent handling...\n", testCount)
	previousExecutor := "old-assignee-789"
	taskAssignedEvent := event.NewTaskAssignedEvent(
		"task-123",
		"project-456",
		"new-assignee-101",
		"assigner-202",
		&previousExecutor,
	)

	if err := handler.Handle(taskAssignedEvent); err != nil {
		fmt.Printf("❌ TaskAssignedEvent handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ TaskAssignedEvent handled successfully\n")
		successCount++
	}

	// 测试 WorkSubmittedEvent
	testCount++
	fmt.Printf("\n%d. Testing WorkSubmittedEvent handling...\n", testCount)
	workSubmittedEvent := event.NewWorkSubmittedEvent(
		"task-123",
		"submitter-303",
		"Work completed successfully",
		[]string{"file1.txt", "file2.pdf"},
	)

	if err := handler.Handle(workSubmittedEvent); err != nil {
		fmt.Printf("❌ WorkSubmittedEvent handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ WorkSubmittedEvent handled successfully\n")
		successCount++
	}

	// 测试 TaskCompletedEvent
	testCount++
	fmt.Printf("\n%d. Testing TaskCompletedEvent handling...\n", testCount)
	taskCompletedEvent := event.NewTaskCompletedEvent(
		"task-123",
		"completer-404",
	)

	if err := handler.Handle(taskCompletedEvent); err != nil {
		fmt.Printf("❌ TaskCompletedEvent handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ TaskCompletedEvent handled successfully\n")
		successCount++
	}

	// 测试 ExtensionRequestedEvent
	testCount++
	fmt.Printf("\n%d. Testing ExtensionRequestedEvent handling...\n", testCount)
	extensionEvent := event.NewExtensionRequestedEvent(
		"task-123",
		"request-505",
		"requester-606",
		time.Now().Add(7*24*time.Hour),
		"Need more time for thorough testing",
	)

	if err := handler.Handle(extensionEvent); err != nil {
		fmt.Printf("❌ ExtensionRequestedEvent handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ ExtensionRequestedEvent handled successfully\n")
		successCount++
	}

	// 测试未支持的事件类型
	testCount++
	fmt.Printf("\n%d. Testing unsupported event type...\n", testCount)
	unsupportedEvent := event.NewTaskStatusChangedEvent(
		"task-123",
		"in_progress",
		"completed",
		"status-changer-707",
		"none",
	)

	if err := handler.Handle(unsupportedEvent); err != nil {
		fmt.Printf("❌ Unsupported event handling failed: %v\n", err)
	} else {
		fmt.Printf("✅ Unsupported event handled gracefully (should log warning)\n")
		successCount++
	}

	// 总结
	fmt.Printf("\n")
	fmt.Printf("\n📊 NOTIFICATION HANDLER TEST SUMMARY:")
	fmt.Printf("\n   Total Tests: %d", testCount)
	fmt.Printf("\n   Passed: %d", successCount)
	fmt.Printf("\n   Failed: %d", testCount-successCount)

	if successCount == testCount {
		fmt.Printf("\n🎉 ALL NOTIFICATION HANDLER TESTS PASSED!")
		fmt.Printf("\n✨ Event type assertions and handling working correctly.")
	} else {
		fmt.Printf("\n⚠️  Some notification handler tests failed.")
		fmt.Printf("\n🔧 Please check the event handling logic.")
	}

	fmt.Printf("\n" + "\n")
}
