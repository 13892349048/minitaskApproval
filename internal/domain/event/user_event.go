package event

import (
	"time"

	"github.com/taskflow/internal/domain/valueobject"
)

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	ID           string               `json:"id"`
	UserID       valueobject.UserID   `json:"user_id"`
	Email        string               `json:"email"`
	Username     string               `json:"username"`
	Role         valueobject.UserRole `json:"role"`
	CreatedBy    valueobject.UserID   `json:"created_by"`
	OccurredOn   time.Time            `json:"occurred_on"`
	EventVersion int                  `json:"event_version"`
}

func (e UserCreatedEvent) EventID() string                 { return e.ID }
func (e UserCreatedEvent) EventType() string               { return "user.created" }
func (e UserCreatedEvent) AggregateID() valueobject.UserID { return e.UserID }
func (e UserCreatedEvent) OccurredAt() time.Time           { return e.OccurredOn }
func (e UserCreatedEvent) Version() int                    { return e.EventVersion }

// UserRoleChangedEvent 用户角色变更事件
type UserRoleChangedEvent struct {
	ID           string               `json:"id"`
	UserID       valueobject.UserID   `json:"user_id"`
	OldRole      valueobject.UserRole `json:"old_role"`
	NewRole      valueobject.UserRole `json:"new_role"`
	ChangedBy    valueobject.UserID   `json:"changed_by"`
	Reason       string               `json:"reason,omitempty"`
	OccurredOn   time.Time            `json:"occurred_on"`
	EventVersion int                  `json:"event_version"`
}

func (e UserRoleChangedEvent) EventID() string        { return e.ID }
func (e UserRoleChangedEvent) EventType() string      { return "user.role_changed" }
func (e UserRoleChangedEvent) AggregateID() string    { return string(e.UserID) }
func (e UserRoleChangedEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e UserRoleChangedEvent) Version() int           { return e.EventVersion }
func (e UserRoleChangedEvent) EventData() interface{} { return e }
func (e UserRoleChangedEvent) AggregateType() string  { return "user" }

// UserDeactivatedEvent 用户停用事件
type UserDeactivatedEvent struct {
	ID                  string             `json:"id"`
	UserID              valueobject.UserID `json:"user_id"`
	DeactivatedBy       valueobject.UserID `json:"deactivated_by"`
	Reason              string             `json:"reason,omitempty"`
	TasksTransferred    bool               `json:"tasks_transferred"`
	ProjectsTransferred bool               `json:"projects_transferred"`
	OccurredOn          time.Time          `json:"occurred_on"`
	EventVersion        int                `json:"event_version"`
}

func (e UserDeactivatedEvent) EventID() string        { return e.ID }
func (e UserDeactivatedEvent) EventType() string      { return "user.deactivated" }
func (e UserDeactivatedEvent) AggregateID() string    { return string(e.UserID) }
func (e UserDeactivatedEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e UserDeactivatedEvent) Version() int           { return e.EventVersion }
func (e UserDeactivatedEvent) EventData() interface{} { return e }
func (e UserDeactivatedEvent) AggregateType() string  { return "user" }

// UserDepartmentTransferredEvent 用户部门转移事件
type UserDepartmentTransferredEvent struct {
	ID               string                   `json:"id"`
	UserID           valueobject.UserID       `json:"user_id"`
	FromDepartmentID valueobject.DepartmentID `json:"from_department_id"`
	ToDepartmentID   valueobject.DepartmentID `json:"to_department_id"`
	OldManagerID     *valueobject.UserID      `json:"old_manager_id,omitempty"`
	NewManagerID     valueobject.UserID       `json:"new_manager_id"`
	TransferredBy    valueobject.UserID       `json:"transferred_by"`
	Reason           string                   `json:"reason,omitempty"`
	OccurredOn       time.Time                `json:"occurred_on"`
	EventVersion     int                      `json:"event_version"`
}

func (e UserDepartmentTransferredEvent) EventID() string        { return e.ID }
func (e UserDepartmentTransferredEvent) EventType() string      { return "user.department_transferred" }
func (e UserDepartmentTransferredEvent) AggregateID() string    { return string(e.UserID) }
func (e UserDepartmentTransferredEvent) OccurredAt() time.Time  { return e.OccurredOn }
func (e UserDepartmentTransferredEvent) Version() int           { return e.EventVersion }
func (e UserDepartmentTransferredEvent) EventData() interface{} { return e }
func (e UserDepartmentTransferredEvent) AggregateType() string  { return "user" }
