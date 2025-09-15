package event

import (
	"time"
)

// 任务相关事件定义

// TaskCreatedEvent 任务创建事件
type TaskCreatedEvent struct {
	*BaseEvent
	TaskID        string    `json:"task_id"`
	Title         string    `json:"title"`
	ProjectID     string    `json:"project_id"`
	CreatorID     string    `json:"creator_id"`
	ResponsibleID string    `json:"responsible_id"`
	TaskType      string    `json:"task_type"`
	Priority      string    `json:"priority"`
	DueDate       time.Time `json:"due_date"`
}

func NewTaskCreatedEvent(taskID, title, projectID, creatorID, responsibleID, taskType, priority string, dueDate time.Time) *TaskCreatedEvent {
	event := &TaskCreatedEvent{
		TaskID:        taskID,
		Title:         title,
		ProjectID:     projectID,
		CreatorID:     creatorID,
		ResponsibleID: responsibleID,
		TaskType:      taskType,
		Priority:      priority,
		DueDate:       dueDate,
	}

	event.BaseEvent = NewBaseEvent("TaskCreated", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskCreatedEvent) EventData() interface{} {
	return e
}

// TaskAssignedEvent 任务分配事件
type TaskAssignedEvent struct {
	*BaseEvent
	TaskID             string  `json:"task_id"`
	ProjectID          string  `json:"project_id"`
	ExecutorID         string  `json:"executor_id"`
	AssignerID         string  `json:"assigner_id"`
	PreviousExecutorID *string `json:"previous_executor_id,omitempty"`
}

func NewTaskAssignedEvent(taskID, projectID, executorID, assignerID string, previousExecutorID *string) *TaskAssignedEvent {
	event := &TaskAssignedEvent{
		TaskID:             taskID,
		ProjectID:          projectID,
		ExecutorID:         executorID,
		AssignerID:         assignerID,
		PreviousExecutorID: previousExecutorID,
	}

	event.BaseEvent = NewBaseEvent("TaskAssigned", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskAssignedEvent) EventData() interface{} {
	return e
}

// TaskPriorityChangedEvent 任务优先级变更事件
type TaskPriorityChangedEvent struct {
	*BaseEvent
	TaskID      string `json:"task_id"`
	OldPriority string `json:"old_priority"`
	NewPriority string `json:"new_priority"`
	ChangedBy   string `json:"changed_by"`
}

func NewTaskPriorityChangedEvent(taskID, oldPriority, newPriority, changedBy string) *TaskPriorityChangedEvent {
	event := &TaskPriorityChangedEvent{
		TaskID:      taskID,
		OldPriority: oldPriority,
		NewPriority: newPriority,
		ChangedBy:   changedBy,
	}

	event.BaseEvent = NewBaseEvent("TaskPriorityChanged", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskPriorityChangedEvent) EventData() interface{} {
	return e
}

// TaskStatusChangedEvent 任务状态变更事件
type TaskStatusChangedEvent struct {
	*BaseEvent
	TaskID       string `json:"task_id"`
	OldStatus    string `json:"old_status"`
	NewStatus    string `json:"new_status"`
	ChangedBy    string `json:"changed_by"`
	ChangeReason string `json:"change_reason,omitempty"`
}

func NewTaskStatusChangedEvent(taskID, oldStatus, newStatus, changedBy, changeReason string) *TaskStatusChangedEvent {
	// eventData := map[string]interface{}{
	// 	"task_id":       taskID,
	// 	"old_status":    oldStatus,
	// 	"new_status":    newStatus,
	// 	"changed_by":    changedBy,
	// 	"change_reason": changeReason,
	// }

	event := &TaskStatusChangedEvent{
		TaskID:       taskID,
		OldStatus:    oldStatus,
		NewStatus:    newStatus,
		ChangedBy:    changedBy,
		ChangeReason: changeReason,
	}

	event.BaseEvent = NewBaseEvent("TaskStatusChanged", taskID, "Task")
	return event
}

// ParticipantAddedEvent 参与者添加事件
type ParticipantAddedEvent struct {
	*BaseEvent
	TaskID        string `json:"task_id"`
	ParticipantID string `json:"participant_id"`
	AddedBy       string `json:"added_by"`
	Role          string `json:"role"`
}

func NewParticipantAddedEvent(taskID, participantID, addedBy, role string) *ParticipantAddedEvent {
	// eventData := map[string]interface{}{
	// 	"task_id":        taskID,
	// 	"participant_id": participantID,
	// 	"added_by":       addedBy,
	// 	"role":           role,
	// }

	event := &ParticipantAddedEvent{
		TaskID:        taskID,
		ParticipantID: participantID,
		AddedBy:       addedBy,
		Role:          role,
	}

	event.BaseEvent = NewBaseEvent("ParticipantAdded", taskID, "Task")
	return event
}

// ParticipantRemovedEvent 参与者移除事件
type ParticipantRemovedEvent struct {
	*BaseEvent
	TaskID        string `json:"task_id"`
	ParticipantID string `json:"participant_id"`
	RemovedBy     string `json:"removed_by"`
	Reason        string `json:"reason,omitempty"`
}

func NewParticipantRemovedEvent(taskID, participantID, removedBy, reason string) *ParticipantRemovedEvent {
	// eventData := map[string]interface{}{
	// 	"task_id":        taskID,
	// 	"participant_id": participantID,
	// 	"removed_by":     removedBy,
	// 	"reason":         reason,
	// }

	event := &ParticipantRemovedEvent{
		TaskID:        taskID,
		ParticipantID: participantID,
		RemovedBy:     removedBy,
		Reason:        reason,
	}

	event.BaseEvent = NewBaseEvent("ParticipantRemoved", taskID, "Task")
	return event
}

// WorkSubmittedEvent 工作提交事件
type WorkSubmittedEvent struct {
	*BaseEvent
	TaskID        string   `json:"task_id"`
	ParticipantID string   `json:"participant_id"`
	WorkContent   string   `json:"work_content"`
	Attachments   []string `json:"attachments,omitempty"`
}

func NewWorkSubmittedEvent(taskID, participantID, workContent string, attachments []string) *WorkSubmittedEvent {
	event := &WorkSubmittedEvent{
		TaskID:        taskID,
		ParticipantID: participantID,
		WorkContent:   workContent,
		Attachments:   attachments,
	}

	event.BaseEvent = NewBaseEvent("WorkSubmitted", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *WorkSubmittedEvent) EventData() interface{} {
	return e
}

// WorkReviewedEvent 工作审核事件
type WorkReviewedEvent struct {
	*BaseEvent
	TaskID        string `json:"task_id"`
	ParticipantID string `json:"participant_id"`
	ReviewerID    string `json:"reviewer_id"`
	Approved      bool   `json:"approved"`
	Comment       string `json:"comment,omitempty"`
}

func NewWorkReviewedEvent(taskID, participantID, reviewerID string, approved bool, comment string) *WorkReviewedEvent {
	event := &WorkReviewedEvent{
		TaskID:        taskID,
		ParticipantID: participantID,
		ReviewerID:    reviewerID,
		Approved:      approved,
		Comment:       comment,
	}

	event.BaseEvent = NewBaseEvent("WorkReviewed", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *WorkReviewedEvent) EventData() interface{} {
	return e
}

// TaskCompletionSubmittedEvent 任务完成提交事件
type TaskCompletionSubmittedEvent struct {
	*BaseEvent
	TaskID        string `json:"task_id"`
	ResponsibleID string `json:"responsible_id"`
	Summary       string `json:"summary"`
}

func NewTaskCompletionSubmittedEvent(taskID, responsibleID, summary string) *TaskCompletionSubmittedEvent {
	event := &TaskCompletionSubmittedEvent{
		TaskID:        taskID,
		ResponsibleID: responsibleID,
		Summary:       summary,
	}

	event.BaseEvent = NewBaseEvent("TaskCompletionSubmitted", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskCompletionSubmittedEvent) EventData() interface{} {
	return e
}

// TaskCompletedEvent 任务完成事件
type TaskCompletedEvent struct {
	*BaseEvent
	TaskID      string    `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
	CompletedBy string    `json:"completed_by"`
}

func NewTaskCompletedEvent(taskID, completedBy string) *TaskCompletedEvent {
	completedAt := time.Now()
	event := &TaskCompletedEvent{
		TaskID:      taskID,
		CompletedAt: completedAt,
		CompletedBy: completedBy,
	}

	event.BaseEvent = NewBaseEvent("TaskCompleted", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskCompletedEvent) EventData() interface{} {
	return e
}

// TaskRejectedEvent 任务拒绝事件
type TaskRejectedEvent struct {
	*BaseEvent
	TaskID     string `json:"task_id"`
	RejectedBy string `json:"rejected_by"`
	Comment    string `json:"comment"`
}

func NewTaskRejectedEvent(taskID, rejectedBy, comment string) *TaskRejectedEvent {
	event := &TaskRejectedEvent{
		TaskID:     taskID,
		RejectedBy: rejectedBy,
		Comment:    comment,
	}

	event.BaseEvent = NewBaseEvent("TaskRejected", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *TaskRejectedEvent) EventData() interface{} {
	return e
}

// ExtensionRequestedEvent 延期申请事件
type ExtensionRequestedEvent struct {
	*BaseEvent
	TaskID      string    `json:"task_id"`
	RequestID   string    `json:"request_id"`
	RequesterID string    `json:"requester_id"`
	NewDueDate  time.Time `json:"new_due_date"`
	Reason      string    `json:"reason"`
}

func NewExtensionRequestedEvent(taskID, requestID, requesterID string, newDueDate time.Time, reason string) *ExtensionRequestedEvent {
	event := &ExtensionRequestedEvent{
		TaskID:      taskID,
		RequestID:   requestID,
		RequesterID: requesterID,
		NewDueDate:  newDueDate,
		Reason:      reason,
	}

	event.BaseEvent = NewBaseEvent("ExtensionRequested", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *ExtensionRequestedEvent) EventData() interface{} {
	return e
}

// ExtensionApprovedEvent 延期批准事件
type ExtensionApprovedEvent struct {
	*BaseEvent
	TaskID     string    `json:"task_id"`
	RequestID  string    `json:"request_id"`
	ReviewerID string    `json:"reviewer_id"`
	NewDueDate time.Time `json:"new_due_date"`
}

func NewExtensionApprovedEvent(taskID, requestID, reviewerID string, newDueDate time.Time) *ExtensionApprovedEvent {
	event := &ExtensionApprovedEvent{
		TaskID:     taskID,
		RequestID:  requestID,
		ReviewerID: reviewerID,
		NewDueDate: newDueDate,
	}

	event.BaseEvent = NewBaseEvent("ExtensionApproved", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *ExtensionApprovedEvent) EventData() interface{} {
	return e
}

// ExtensionRejectedEvent 延期拒绝事件
type ExtensionRejectedEvent struct {
	*BaseEvent
	TaskID     string `json:"task_id"`
	RequestID  string `json:"request_id"`
	ReviewerID string `json:"reviewer_id"`
	Comment    string `json:"comment"`
}

func NewExtensionRejectedEvent(taskID, requestID, reviewerID, comment string) *ExtensionRejectedEvent {
	event := &ExtensionRejectedEvent{
		TaskID:     taskID,
		RequestID:  requestID,
		ReviewerID: reviewerID,
		Comment:    comment,
	}

	event.BaseEvent = NewBaseEvent("ExtensionRejected", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *ExtensionRejectedEvent) EventData() interface{} {
	return e
}

// NextExecutionPreparedEvent 下次执行准备事件（重复任务）
type NextExecutionPreparedEvent struct {
	*BaseEvent
	TaskID        string    `json:"task_id"`
	ExecutionID   string    `json:"execution_id"`
	ExecutionDate time.Time `json:"execution_date"`
}

func NewNextExecutionPreparedEvent(taskID, executionID string, executionDate time.Time) *NextExecutionPreparedEvent {
	// eventData := map[string]interface{}{
	// 	"task_id":        taskID,
	// 	"execution_id":   executionID,
	// 	"execution_date": executionDate,
	// }

	event := &NextExecutionPreparedEvent{
		TaskID:        taskID,
		ExecutionID:   executionID,
		ExecutionDate: executionDate,
	}

	event.BaseEvent = NewBaseEvent("NextExecutionPrepared", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *NextExecutionPreparedEvent) EventData() interface{} {
	return e
}

// AllParticipantsCompletedEvent 所有参与者完成事件
type AllParticipantsCompletedEvent struct {
	*BaseEvent
	TaskID          string   `json:"task_id"`
	ParticipantIDs  []string `json:"participant_ids"`
	CompletionCount int      `json:"completion_count"`
}

func NewAllParticipantsCompletedEvent(taskID string, participantIDs []string, completionCount int) *AllParticipantsCompletedEvent {
	// eventData := map[string]interface{}{
	// 	"task_id":          taskID,
	// 	"participant_ids":  participantIDs,
	// 	"completion_count": completionCount,
	// }

	event := &AllParticipantsCompletedEvent{
		TaskID:          taskID,
		ParticipantIDs:  participantIDs,
		CompletionCount: completionCount,
	}

	event.BaseEvent = NewBaseEvent("AllParticipantsCompleted", taskID, "Task")
	return event
}

// EventData 实现 DomainEvent 接口
func (e *AllParticipantsCompletedEvent) EventData() interface{} {
	return e
}
