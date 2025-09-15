package aggregate

import (
	"testing"
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewPermission(t *testing.T) {
	// Arrange
	id := valueobject.PermissionID("perm-123")
	name := "Create Task"
	resource := valueobject.ResourceTypeTask
	action := valueobject.ActionTypeCreate
	description := "Permission to create tasks"

	// Act
	permission := NewPermission(id, name, resource, action, description)

	// Assert
	assert.NotNil(t, permission)
	assert.Equal(t, id, permission.ID)
	assert.Equal(t, name, permission.Name)
	assert.Equal(t, resource, permission.Resource)
	assert.Equal(t, action, permission.Action)
	assert.Equal(t, description, permission.Description)
	assert.False(t, permission.CreatedAt.IsZero())
	assert.False(t, permission.UpdatedAt.IsZero())
}

func TestPermission_UpdateDescription(t *testing.T) {
	// Arrange
	permission := createTestPermission()
	newDescription := "Updated description"
	originalUpdatedAt := permission.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	permission.UpdateDescription(newDescription)

	// Assert
	assert.Equal(t, newDescription, permission.Description)
	assert.True(t, permission.UpdatedAt.After(originalUpdatedAt))
}

func TestPermission_Matches(t *testing.T) {
	tests := []struct {
		name               string
		permissionResource valueobject.ResourceType
		permissionAction   valueobject.ActionType
		testResource       valueobject.ResourceType
		testAction         valueobject.ActionType
		expected           bool
	}{
		{
			name:               "Exact match",
			permissionResource: valueobject.ResourceTypeTask,
			permissionAction:   valueobject.ActionTypeCreate,
			testResource:       valueobject.ResourceTypeTask,
			testAction:         valueobject.ActionTypeCreate,
			expected:           true,
		},
		{
			name:               "Different resource",
			permissionResource: valueobject.ResourceTypeTask,
			permissionAction:   valueobject.ActionTypeCreate,
			testResource:       valueobject.ResourceTypeProject,
			testAction:         valueobject.ActionTypeCreate,
			expected:           false,
		},
		{
			name:               "Different action",
			permissionResource: valueobject.ResourceTypeTask,
			permissionAction:   valueobject.ActionTypeCreate,
			testResource:       valueobject.ResourceTypeTask,
			testAction:         valueobject.ActionTypeDelete,
			expected:           false,
		},
		{
			name:               "Both different",
			permissionResource: valueobject.ResourceTypeTask,
			permissionAction:   valueobject.ActionTypeCreate,
			testResource:       valueobject.ResourceTypeProject,
			testAction:         valueobject.ActionTypeDelete,
			expected:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			permission := NewPermission(
				valueobject.PermissionID("test-perm"),
				"Test Permission",
				tt.permissionResource,
				tt.permissionAction,
				"Test description",
			)

			// Act
			result := permission.Matches(tt.testResource, tt.testAction)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create a test permission
func createTestPermission() *Permission {
	return NewPermission(
		valueobject.PermissionID("test-perm-123"),
		"Test Permission",
		valueobject.ResourceTypeTask,
		valueobject.ActionTypeCreate,
		"Test permission description",
	)
}
