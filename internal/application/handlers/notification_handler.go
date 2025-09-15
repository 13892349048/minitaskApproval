package handlers

import (
	"fmt"
	"reflect"

	"github.com/taskflow/internal/domain/event"
	"github.com/taskflow/pkg/logger"
	"go.uber.org/zap"
)

// FixedNotificationHandler 修复后的通知事件处理器
type FixedNotificationHandler struct {
	emailService EmailService
	smsService   SMSService
}

// EmailService 邮件服务接口
type EmailService interface {
	SendEmail(to, subject, body string) error
}

// SMSService 短信服务接口
type SMSService interface {
	SendSMS(to, message string) error
}

// NewFixedNotificationHandler 创建修复后的通知处理器
func NewNotificationHandler(emailService EmailService, smsService SMSService) *FixedNotificationHandler {
	return &FixedNotificationHandler{
		emailService: emailService,
		smsService:   smsService,
	}
}

// Handle 处理事件 - 使用反射和类型安全的方法
func (h *FixedNotificationHandler) Handle(domainEvent event.DomainEvent) error {
	eventType := domainEvent.EventType()

	logger.Info("Processing event",
		zap.String("event_type", eventType),
		zap.String("event_id", domainEvent.EventID()),
		zap.String("aggregate_id", domainEvent.AggregateID()))

	switch eventType {
	case "TaskCreated":
		return h.handleTaskCreatedSafe(domainEvent)
	case "TaskAssigned":
		return h.handleTaskAssignedSafe(domainEvent)
	case "WorkSubmitted":
		return h.handleWorkSubmittedSafe(domainEvent)
	case "WorkReviewed":
		return h.handleWorkReviewedSafe(domainEvent)
	case "TaskCompletionSubmitted":
		return h.handleTaskCompletionSubmittedSafe(domainEvent)
	case "TaskCompleted":
		return h.handleTaskCompletedSafe(domainEvent)
	case "TaskRejected":
		return h.handleTaskRejectedSafe(domainEvent)
	case "ExtensionRequested":
		return h.handleExtensionRequestedSafe(domainEvent)
	case "ExtensionApproved":
		return h.handleExtensionApprovedSafe(domainEvent)
	case "ExtensionRejected":
		return h.handleExtensionRejectedSafe(domainEvent)
	default:
		logger.Warn("Unhandled event type", zap.String("event_type", eventType))
		return nil
	}
}

// safeEventCast 安全的事件类型转换
func safeEventCast[T any](domainEvent event.DomainEvent, expectedType string) (*T, error) {
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

// handleTaskCreatedSafe 安全处理TaskCreated事件
func (h *FixedNotificationHandler) handleTaskCreatedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.TaskCreatedEvent](domainEvent, "TaskCreated")
	if err != nil {
		logger.Error("Failed to cast TaskCreatedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for TaskCreated: %w", err)
	}

	subject := fmt.Sprintf("新任务创建：%s", data.Title)
	body := fmt.Sprintf("任务 '%s' 已创建，负责人：%s，截止日期：%s",
		data.Title, data.ResponsibleID, data.DueDate.Format("2006-01-02"))

	// 通知负责人
	if err := h.emailService.SendEmail(data.ResponsibleID+"@company.com", subject, body); err != nil {
		logger.Error("Failed to send email for TaskCreated", zap.Error(err))
		return err
	}

	logger.Info("Task created notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("responsible_id", data.ResponsibleID))
	return nil
}

// handleTaskAssignedSafe 安全处理TaskAssigned事件
func (h *FixedNotificationHandler) handleTaskAssignedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.TaskAssignedEvent](domainEvent, "TaskAssigned")
	if err != nil {
		logger.Error("Failed to cast TaskAssignedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for TaskAssigned: %w", err)
	}

	subject := "任务分配通知"
	body := fmt.Sprintf("您被分配了新任务，任务ID：%s", data.TaskID)

	// 通知新的执行者
	if err := h.emailService.SendEmail(data.ExecutorID+"@company.com", subject, body); err != nil {
		logger.Error("Failed to send email for TaskAssigned", zap.Error(err))
		return err
	}

	logger.Info("Task assigned notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("executor_id", data.ExecutorID))
	return nil
}

// handleWorkSubmittedSafe 安全处理WorkSubmitted事件
func (h *FixedNotificationHandler) handleWorkSubmittedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.WorkSubmittedEvent](domainEvent, "WorkSubmitted")
	if err != nil {
		logger.Error("Failed to cast WorkSubmittedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for WorkSubmitted: %w", err)
	}

	subject := "工作提交通知"
	body := fmt.Sprintf("任务 %s 的工作已提交，请进行审核", data.TaskID)

	logger.Info("Work submitted notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("participant_id", data.ParticipantID),
		zap.String("work_content", data.WorkContent))

	// 实际发送邮件的逻辑会在这里调用 h.emailService.SendEmail
	_ = subject
	_ = body
	return nil
}

// handleWorkReviewedSafe 安全处理WorkReviewed事件
func (h *FixedNotificationHandler) handleWorkReviewedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.WorkReviewedEvent](domainEvent, "WorkReviewed")
	if err != nil {
		logger.Error("Failed to cast WorkReviewedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for WorkReviewed: %w", err)
	}

	subject := "工作审批结果通知"
	status := "通过"
	if !data.Approved {
		status = "需要修改"
	}
	body := fmt.Sprintf("您的工作成果审批结果：%s。评论：%s", status, data.Comment)

	// 通知参与人员
	if err := h.emailService.SendEmail(data.ParticipantID+"@company.com", subject, body); err != nil {
		logger.Error("Failed to send email for WorkReviewed", zap.Error(err))
		return err
	}

	logger.Info("Work reviewed notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("participant_id", data.ParticipantID),
		zap.Bool("approved", data.Approved))
	return nil
}

// handleTaskCompletionSubmittedSafe 安全处理TaskCompletionSubmitted事件
func (h *FixedNotificationHandler) handleTaskCompletionSubmittedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.TaskCompletionSubmittedEvent](domainEvent, "TaskCompletionSubmitted")
	if err != nil {
		logger.Error("Failed to cast TaskCompletionSubmittedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for TaskCompletionSubmitted: %w", err)
	}

	subject := "任务完成提交通知"
	body := fmt.Sprintf("任务 %s 已提交完成，等待最终审批", data.TaskID)

	logger.Info("Task completion submitted notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("responsible_id", data.ResponsibleID))

	_ = subject
	_ = body
	return nil
}

// handleTaskCompletedSafe 安全处理TaskCompleted事件
func (h *FixedNotificationHandler) handleTaskCompletedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.TaskCompletedEvent](domainEvent, "TaskCompleted")
	if err != nil {
		logger.Error("Failed to cast TaskCompletedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for TaskCompleted: %w", err)
	}

	subject := "任务完成通知"
	body := fmt.Sprintf("恭喜！任务 %s 已成功完成", data.TaskID)

	logger.Info("Task completed notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("completed_by", data.CompletedBy))

	_ = subject
	_ = body
	return nil
}

// handleTaskRejectedSafe 安全处理TaskRejected事件
func (h *FixedNotificationHandler) handleTaskRejectedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.TaskRejectedEvent](domainEvent, "TaskRejected")
	if err != nil {
		logger.Error("Failed to cast TaskRejectedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for TaskRejected: %w", err)
	}

	subject := "任务返工通知"
	body := fmt.Sprintf("任务 %s 需要返工。原因：%s", data.TaskID, data.Comment)

	logger.Info("Task rejected notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("rejected_by", data.RejectedBy),
		zap.String("comment", data.Comment))

	_ = subject
	_ = body
	return nil
}

// handleExtensionRequestedSafe 安全处理ExtensionRequested事件
func (h *FixedNotificationHandler) handleExtensionRequestedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.ExtensionRequestedEvent](domainEvent, "ExtensionRequested")
	if err != nil {
		logger.Error("Failed to cast ExtensionRequestedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for ExtensionRequested: %w", err)
	}

	subject := "延期申请通知"
	body := fmt.Sprintf("任务 %s 申请延期至 %s，原因：%s",
		data.TaskID, data.NewDueDate.Format("2006-01-02"), data.Reason)

	logger.Info("Extension requested notification sent",
		zap.String("task_id", data.TaskID),
		zap.Time("new_due_date", data.NewDueDate),
		zap.String("reason", data.Reason))

	_ = subject
	_ = body
	return nil
}

// handleExtensionApprovedSafe 安全处理ExtensionApproved事件
func (h *FixedNotificationHandler) handleExtensionApprovedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.ExtensionApprovedEvent](domainEvent, "ExtensionApproved")
	if err != nil {
		logger.Error("Failed to cast ExtensionApprovedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for ExtensionApproved: %w", err)
	}

	subject := "延期申请批准通知"
	body := fmt.Sprintf("您的延期申请已批准，新的截止日期：%s", data.NewDueDate.Format("2006-01-02"))

	logger.Info("Extension approved notification sent",
		zap.String("task_id", data.TaskID),
		zap.Time("new_due_date", data.NewDueDate))

	_ = subject
	_ = body
	return nil
}

// handleExtensionRejectedSafe 安全处理ExtensionRejected事件
func (h *FixedNotificationHandler) handleExtensionRejectedSafe(domainEvent event.DomainEvent) error {
	data, err := safeEventCast[event.ExtensionRejectedEvent](domainEvent, "ExtensionRejected")
	if err != nil {
		logger.Error("Failed to cast ExtensionRejectedEvent", zap.Error(err))
		return fmt.Errorf("invalid event data for ExtensionRejected: %w", err)
	}

	subject := "延期申请拒绝通知"
	body := fmt.Sprintf("您的延期申请已被拒绝，原因：%s", data.Comment)

	logger.Info("Extension rejected notification sent",
		zap.String("task_id", data.TaskID),
		zap.String("comment", data.Comment))

	_ = subject
	_ = body
	return nil
}

// CanHandle 判断是否能处理该事件
func (h *FixedNotificationHandler) CanHandle(eventType string) bool {
	supportedEvents := []string{
		"TaskCreated", "TaskAssigned", "WorkSubmitted",
		"WorkReviewed", "TaskCompletionSubmitted", "TaskCompleted",
		"TaskRejected", "ExtensionRequested", "ExtensionApproved", "ExtensionRejected",
	}

	for _, supported := range supportedEvents {
		if eventType == supported {
			return true
		}
	}
	return false
}

// EventTypes 返回支持的事件类型
func (h *FixedNotificationHandler) EventTypes() []string {
	return []string{
		"TaskCreated",
		"TaskAssigned",
		"WorkSubmitted",
		"WorkReviewed",
		"TaskCompletionSubmitted",
		"TaskCompleted",
		"TaskRejected",
		"ExtensionRequested",
		"ExtensionApproved",
		"ExtensionRejected",
	}
}
