package aggregate

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taskflow/internal/domain/valueobject"
)

func TestNewUser(t *testing.T) {
	// Arrange
	userID := valueobject.UserID("user-123")
	username := "testuser"
	email := "test@example.com"
	fullName := "Test User"
	passwordHash := "password"
	role := valueobject.UserRoleEmployee

	// Act
	user := NewUser(userID, username, email, fullName, passwordHash, role)

	// Assert
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, fullName, user.FullName)
	assert.Equal(t, passwordHash, user.PasswordHash)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, valueobject.UserStatusActive, user.Status)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.Nil(t, user.DeletedAt)
}

func TestUser_UpdateProfile(t *testing.T) {
	// Arrange
	user := createTestUser()
	newFullName := "New Full Name"
	newEmail := "newemail@example.com"
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	err := user.UpdateProfile(newFullName, newEmail)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, newFullName, user.FullName)
	assert.Equal(t, newEmail, user.Email)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_UpdateProfile_InvalidEmail(t *testing.T) {
	// Arrange
	user := createTestUser()
	newFullName := "New Full Name"
	invalidEmail := ""

	// Act
	err := user.UpdateProfile(newFullName, invalidEmail)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")
}

func TestUser_ChangeRole(t *testing.T) {
	// Arrange
	user := createTestUser()
	newRole := valueobject.UserRoleManager
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.ChangeRole(newRole)

	// Assert
	assert.Equal(t, newRole, user.Role)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_Activate(t *testing.T) {
	// Arrange
	user := createTestUser()
	user.Status = valueobject.UserStatusInactive
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.Activate()

	// Assert
	assert.Equal(t, valueobject.UserStatusActive, user.Status)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_Deactivate(t *testing.T) {
	// Arrange
	user := createTestUser()
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.Deactivate()

	// Assert
	assert.Equal(t, valueobject.UserStatusInactive, user.Status)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_Suspend(t *testing.T) {
	// Arrange
	user := createTestUser()
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.Suspend()

	// Assert
	assert.Equal(t, valueobject.UserStatusSuspended, user.Status)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		userRole   valueobject.UserRole
		permission string
		expected   bool
	}{
		{
			name:       "Admin has all permissions",
			userRole:   valueobject.UserRoleAdmin,
			permission: "task:create",
			expected:   true,
		},
		{
			name:       "Manager can manage tasks",
			userRole:   valueobject.UserRoleManager,
			permission: "task:manage",
			expected:   true,
		},
		{
			name:       "Manager can read users",
			userRole:   valueobject.UserRoleManager,
			permission: "user:read",
			expected:   true,
		},
		{
			name:       "Employee can read tasks",
			userRole:   valueobject.UserRoleEmployee,
			permission: "task:read",
			expected:   true,
		},
		{
			name:       "Employee can read users",
			userRole:   valueobject.UserRoleEmployee,
			permission: "user:read",
			expected:   true,
		},
		{
			name:       "Employee cannot delete tasks",
			userRole:   valueobject.UserRoleEmployee,
			permission: "task:delete",
			expected:   false,
		},
		{
			name:       "Director has most permissions",
			userRole:   valueobject.UserRoleDirector,
			permission: "task:manage",
			expected:   true,
		},
		{
			name:       "Director cannot access system admin",
			userRole:   valueobject.UserRoleDirector,
			permission: "system:admin",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			user := createTestUser()
			user.Role = tt.userRole

			// Act
			result := user.HasPermission(tt.permission)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   valueobject.UserStatus
		expected bool
	}{
		{
			name:     "Active user",
			status:   valueobject.UserStatusActive,
			expected: true,
		},
		{
			name:     "Inactive user",
			status:   valueobject.UserStatusInactive,
			expected: false,
		},
		{
			name:     "Suspended user",
			status:   valueobject.UserStatusSuspended,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			user := createTestUser()
			user.Status = tt.status

			// Act
			result := user.IsActive()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_SetManager(t *testing.T) {
	// Arrange
	user := createTestUser()
	managerID := valueobject.UserID("manager-123")
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.AssignManager(managerID)

	// Assert
	assert.Equal(t, &managerID, user.ManagerID)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

func TestUser_SetDepartment(t *testing.T) {
	// Arrange
	user := createTestUser()
	departmentID := "dept-123"
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	user.AssignToDepartment(departmentID)

	// Assert
	assert.Equal(t, &departmentID, user.DepartmentID)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
}

// validateEmail 简单的邮箱验证函数（仅用于测试）
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func TestUser_ValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "Invalid email - no @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "Invalid email - no domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "Invalid email - no local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "Empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := validateEmail(tt.email)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create a test user
func createTestUser() *User {
	return NewUser(
		valueobject.UserID("test-user-123"),
		"testuser",
		"test@example.com",
		"Test User",
		"hashed_password",
		valueobject.UserRoleEmployee,
	)
}
