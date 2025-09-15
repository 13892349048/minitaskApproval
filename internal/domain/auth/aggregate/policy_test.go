package aggregate

import (
	"testing"
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNewPolicy(t *testing.T) {
	// Arrange
	id := valueobject.PolicyID("policy-123")
	name := "Task Owner Policy"
	description := "Users can only modify their own tasks"
	resource := valueobject.ResourceTypeTask
	action := valueobject.ActionTypeUpdate
	effect := valueobject.PolicyEffectAllow
	conditions := valueobject.PolicyConditions{
		"resource.owner_id": "${user.id}",
	}
	priority := 100

	// Act
	policy := NewPolicy(id, name, description, resource, action, effect, conditions, priority)

	// Assert
	assert.NotNil(t, policy)
	assert.Equal(t, id, policy.ID)
	assert.Equal(t, name, policy.Name)
	assert.Equal(t, description, policy.Description)
	assert.Equal(t, resource, policy.Resource)
	assert.Equal(t, action, policy.Action)
	assert.Equal(t, effect, policy.Effect)
	assert.Equal(t, conditions, policy.Conditions)
	assert.Equal(t, priority, policy.Priority)
	assert.True(t, policy.IsActive)
	assert.False(t, policy.CreatedAt.IsZero())
	assert.False(t, policy.UpdatedAt.IsZero())
}

func TestPolicy_UpdatePolicy(t *testing.T) {
	// Arrange
	policy := createTestPolicy()
	newName := "Updated Policy"
	newDescription := "Updated description"
	newEffect := valueobject.PolicyEffectDeny
	newConditions := valueobject.PolicyConditions{
		"user.department": "IT",
	}
	newPriority := 200
	originalUpdatedAt := policy.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	policy.UpdatePolicy(newName, newDescription, newEffect, newConditions, newPriority)

	// Assert
	assert.Equal(t, newName, policy.Name)
	assert.Equal(t, newDescription, policy.Description)
	assert.Equal(t, newEffect, policy.Effect)
	assert.Equal(t, newConditions, policy.Conditions)
	assert.Equal(t, newPriority, policy.Priority)
	assert.True(t, policy.UpdatedAt.After(originalUpdatedAt))
}

func TestPolicy_Activate(t *testing.T) {
	// Arrange
	policy := createTestPolicy()
	policy.IsActive = false
	originalUpdatedAt := policy.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	policy.Activate()

	// Assert
	assert.True(t, policy.IsActive)
	assert.True(t, policy.UpdatedAt.After(originalUpdatedAt))
}

func TestPolicy_Deactivate(t *testing.T) {
	// Arrange
	policy := createTestPolicy()
	originalUpdatedAt := policy.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	policy.Deactivate()

	// Assert
	assert.False(t, policy.IsActive)
	assert.True(t, policy.UpdatedAt.After(originalUpdatedAt))
}

func TestPolicy_Matches(t *testing.T) {
	tests := []struct {
		name           string
		policyResource valueobject.ResourceType
		policyAction   valueobject.ActionType
		policyActive   bool
		testResource   valueobject.ResourceType
		testAction     valueobject.ActionType
		expected       bool
	}{
		{
			name:           "Active policy matches",
			policyResource: valueobject.ResourceTypeTask,
			policyAction:   valueobject.ActionTypeUpdate,
			policyActive:   true,
			testResource:   valueobject.ResourceTypeTask,
			testAction:     valueobject.ActionTypeUpdate,
			expected:       true,
		},
		{
			name:           "Inactive policy doesn't match",
			policyResource: valueobject.ResourceTypeTask,
			policyAction:   valueobject.ActionTypeUpdate,
			policyActive:   false,
			testResource:   valueobject.ResourceTypeTask,
			testAction:     valueobject.ActionTypeUpdate,
			expected:       false,
		},
		{
			name:           "Different resource doesn't match",
			policyResource: valueobject.ResourceTypeTask,
			policyAction:   valueobject.ActionTypeUpdate,
			policyActive:   true,
			testResource:   valueobject.ResourceTypeProject,
			testAction:     valueobject.ActionTypeUpdate,
			expected:       false,
		},
		{
			name:           "Different action doesn't match",
			policyResource: valueobject.ResourceTypeTask,
			policyAction:   valueobject.ActionTypeUpdate,
			policyActive:   true,
			testResource:   valueobject.ResourceTypeTask,
			testAction:     valueobject.ActionTypeDelete,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			policy := NewPolicy(
				valueobject.PolicyID("test-policy"),
				"Test Policy",
				"Test description",
				tt.policyResource,
				tt.policyAction,
				valueobject.PolicyEffectAllow,
				valueobject.PolicyConditions{},
				100,
			)
			policy.IsActive = tt.policyActive

			// Act
			result := policy.Matches(tt.testResource, tt.testAction)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPolicy_EffectTypes(t *testing.T) {
	tests := []struct {
		name   string
		effect valueobject.PolicyEffect
	}{
		{
			name:   "Allow effect",
			effect: valueobject.PolicyEffectAllow,
		},
		{
			name:   "Deny effect",
			effect: valueobject.PolicyEffectDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			policy := NewPolicy(
				valueobject.PolicyID("test-policy"),
				"Test Policy",
				"Test description",
				valueobject.ResourceTypeTask,
				valueobject.ActionTypeUpdate,
				tt.effect,
				valueobject.PolicyConditions{},
				100,
			)

			// Assert
			assert.Equal(t, tt.effect, policy.Effect)
		})
	}
}

func TestPolicy_ComplexConditions(t *testing.T) {
	// Arrange
	complexConditions := valueobject.PolicyConditions{
		"user.department":    "Engineering",
		"resource.owner_id":  "${user.id}",
		"time.hour":          []interface{}{9, 10, 11, 12, 13, 14, 15, 16, 17},
		"resource.priority":  "high",
	}

	// Act
	policy := NewPolicy(
		valueobject.PolicyID("complex-policy"),
		"Complex Policy",
		"Policy with multiple conditions",
		valueobject.ResourceTypeTask,
		valueobject.ActionTypeUpdate,
		valueobject.PolicyEffectAllow,
		complexConditions,
		150,
	)

	// Assert
	assert.Equal(t, complexConditions, policy.Conditions)
	assert.Len(t, policy.Conditions, 4)
	assert.Equal(t, "Engineering", policy.Conditions["user.department"])
	assert.Equal(t, "${user.id}", policy.Conditions["resource.owner_id"])
}

// Helper function to create a test policy
func createTestPolicy() *Policy {
	return NewPolicy(
		valueobject.PolicyID("test-policy-123"),
		"Test Policy",
		"Test policy description",
		valueobject.ResourceTypeTask,
		valueobject.ActionTypeUpdate,
		valueobject.PolicyEffectAllow,
		valueobject.PolicyConditions{
			"resource.owner_id": "${user.id}",
		},
		100,
	)
}
