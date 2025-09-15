package events

import (
	"log"
	"time"

	"github.com/taskflow/internal/application/handlers"
	"github.com/taskflow/pkg/logger"
)

// MockEmailService 模拟邮件服务
type MockEmailService struct{}

func (s *MockEmailService) SendEmail(to, subject, body string) error {
	logger.Logger.Sugar().Infof("Mock Email sent to %s: %s", to, subject)
	return nil
}

// MockSMSService 模拟短信服务
type MockSMSService struct{}

func (s *MockSMSService) SendSMS(to, message string) error {
	logger.Logger.Sugar().Infof("Mock SMS sent to %s: %s", to, message)
	return nil
}

// MockAuditRepository 模拟审计仓储
type MockAuditRepository struct{}

func (r *MockAuditRepository) Save(log *handlers.AuditLog) error {
	// 在实际实现中，这里会保存到数据库
	logger.Logger.Sugar().Infof("Mock Audit log saved: %s - %s", log.EventType, log.AggregateID)
	return nil
}

func (r *MockAuditRepository) FindByAggregateID(aggregateID string, limit int) ([]*handlers.AuditLog, error) {
	return []*handlers.AuditLog{}, nil
}

func (r *MockAuditRepository) FindByEventType(eventType string, limit int) ([]*handlers.AuditLog, error) {
	return []*handlers.AuditLog{}, nil
}

func (r *MockAuditRepository) FindByTimeRange(start, end time.Time, limit int) ([]*handlers.AuditLog, error) {
	return []*handlers.AuditLog{}, nil
}

// MockStatisticsRepository 模拟统计仓储
type MockStatisticsRepository struct{}

func (r *MockStatisticsRepository) SaveProjectStats(projectID string, stats *handlers.TaskStatistics) error {
	log.Printf("Mock Statistics saved for project %s: %+v", projectID, stats)
	return nil
}

func (r *MockStatisticsRepository) GetProjectStats(projectID string) (*handlers.TaskStatistics, error) {
	return &handlers.TaskStatistics{
		ProjectID:   projectID,
		LastUpdated: time.Now(),
	}, nil
}

func (r *MockStatisticsRepository) UpdateTaskCount(projectID string, increment int) error {
	log.Printf("Mock Task count updated for project %s: +%d", projectID, increment)
	return nil
}

func (r *MockStatisticsRepository) UpdateCompletedCount(projectID string, increment int) error {
	log.Printf("Mock Completed count updated for project %s: +%d", projectID, increment)
	return nil
}
