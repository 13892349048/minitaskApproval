package main

import (
	"fmt"
	"time"

	"github.com/taskflow/internal/domain/event"
)

func main1111() {
	// Test TaskCreatedEvent assertion
	taskCreatedEvent := event.NewTaskCreatedEvent(
		"task-123",
		"Test Task",
		"project-456",
		"creator-789",
		"responsible-101",
		"development",
		"high",
		time.Now().Add(24*time.Hour),
	)
	// Test the event data assertion
	data, ok := taskCreatedEvent.EventData().(*event.TaskCreatedEvent)
	if !ok {
		fmt.Printf("❌ TaskCreatedEvent assertion failed\n")
		return
	}
	fmt.Printf("✅ TaskCreatedEvent assertion successful: %s\n", data.Title)

	// Test TaskAssignedEvent assertion
	previousExecutor := "old-assignee-789"
	taskAssignedEvent := event.NewTaskAssignedEvent(
		"task-123",
		"project-456",
		"new-assignee-101",
		"assigner-202",
		&previousExecutor,
	)

	assignedData, ok := taskAssignedEvent.EventData().(*event.TaskAssignedEvent)
	if !ok {
		fmt.Printf("❌ TaskAssignedEvent assertion failed\n")
		return
	}
	fmt.Printf("✅ TaskAssignedEvent assertion successful: %s -> %s\n",
		*assignedData.PreviousExecutorID, assignedData.ExecutorID)

	// Test WorkSubmittedEvent assertion
	workSubmittedEvent := event.NewWorkSubmittedEvent(
		"task-123",
		"submitter-303",
		"Work completed successfully",
		[]string{"file1.txt", "file2.pdf"},
	)

	workData, ok := workSubmittedEvent.EventData().(*event.WorkSubmittedEvent)
	if !ok {
		fmt.Printf("❌ WorkSubmittedEvent assertion failed\n")
		return
	}
	fmt.Printf("✅ WorkSubmittedEvent assertion successful: %s\n", workData.WorkContent)

	fmt.Printf("\n🎉 All event type assertions working correctly!\n")
}
