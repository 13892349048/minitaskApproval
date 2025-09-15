package aggregate

import (
	"testing"
	"time"

	"github.com/taskflow/internal/domain/auth/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRole(t *testing.T) {
	// Arrange
	id := valueobject.RoleID("role-123")
	name := "manager"
	displayName := "Manager"
	description := "Manager role"
	isSystem := false

	// Act
	role := NewRole(id, name, displayName, description, isSystem)

	// Assert
	assert.NotNil(t, role)
	assert.Equal(t, id, role.ID)
	assert.Equal(t, name, role.Name)
	assert.Equal(t, displayName, role.DisplayName)
	assert.Equal(t, description, role.Description)
	assert.Equal(t, isSystem, role.IsSystem)
	assert.Empty(t, role.Permissions)
	assert.False(t, role.CreatedAt.IsZero())
	assert.False(t, role.UpdatedAt.IsZero())
}

func TestRole_AddPermission(t *testing.T) {
	// Arrange
	role := createTestRole()
	permissionID := valueobject.PermissionID("perm-123")
	originalUpdatedAt := role.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	err := role.AddPermission(permissionID)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, role.Permissions, permissionID)
	assert.True(t, role.UpdatedAt.After(originalUpdatedAt))
}

func TestRole_AddPermission_Duplicate(t *testing.T) {
	// Arrange
	role := createTestRole()
	permissionID := valueobject.PermissionID("perm-123")
	
	// Add permission first time
	err := role.AddPermission(permissionID)
	require.NoError(t, err)

	// Act - try to add same permission again
	err = role.AddPermission(permissionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "权限已存在")
	assert.Len(t, role.Permissions, 1) // Should still have only one permission
}

func TestRole_RemovePermission(t *testing.T) {
	// Arrange
	role := createTestRole()
	permissionID := valueobject.PermissionID("perm-123")
	
	// Add permission first
	err := role.AddPermission(permissionID)
	require.NoError(t, err)
	
	originalUpdatedAt := role.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	// Act
	err = role.RemovePermission(permissionID)

	// Assert
	require.NoError(t, err)
	assert.NotContains(t, role.Permissions, permissionID)
	assert.True(t, role.UpdatedAt.After(originalUpdatedAt))
}

func TestRole_RemovePermission_NotFound(t *testing.T) {
	// Arrange
	role := createTestRole()
	permissionID := valueobject.PermissionID("nonexistent-perm")

	// Act
	err := role.RemovePermission(permissionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "权限不存在")
}

func TestRole_RemovePermission_SystemRole(t *testing.T) {
	// Arrange
	role := createTestSystemRole()
	permissionID := valueobject.PermissionID("perm-123")

	// Act
	err := role.RemovePermission(permissionID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "系统角色不能修改权限")
}

func TestRole_HasPermission(t *testing.T) {
	// Arrange
	role := createTestRole()
	permissionID := valueobject.PermissionID("perm-123")
	otherPermissionID := valueobject.PermissionID("other-perm")
	
	// Add one permission
	err := role.AddPermission(permissionID)
	require.NoError(t, err)

	// Act & Assert
	assert.True(t, role.HasPermission(permissionID))
	assert.False(t, role.HasPermission(otherPermissionID))
}

func TestRole_UpdateInfo(t *testing.T) {
	// Arrange
	role := createTestRole()
	newDisplayName := "Updated Manager"
	newDescription := "Updated description"
	originalUpdatedAt := role.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(1 * time.Millisecond)

	// Act
	err := role.UpdateInfo(newDisplayName, newDescription)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, newDisplayName, role.DisplayName)
	assert.Equal(t, newDescription, role.Description)
	assert.True(t, role.UpdatedAt.After(originalUpdatedAt))
}

func TestRole_UpdateInfo_SystemRole(t *testing.T) {
	// Arrange
	role := createTestSystemRole()
	newDisplayName := "Updated System Role"
	newDescription := "Updated description"

	// Act
	err := role.UpdateInfo(newDisplayName, newDescription)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "系统角色不能修改")
}

func TestRole_MultiplePermissions(t *testing.T) {
	// Arrange
	role := createTestRole()
	perm1 := valueobject.PermissionID("perm-1")
	perm2 := valueobject.PermissionID("perm-2")
	perm3 := valueobject.PermissionID("perm-3")

	// Act - Add multiple permissions
	err1 := role.AddPermission(perm1)
	err2 := role.AddPermission(perm2)
	err3 := role.AddPermission(perm3)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	assert.Len(t, role.Permissions, 3)
	assert.Contains(t, role.Permissions, perm1)
	assert.Contains(t, role.Permissions, perm2)
	assert.Contains(t, role.Permissions, perm3)

	// Act - Remove middle permission
	err := role.RemovePermission(perm2)

	// Assert
	require.NoError(t, err)
	assert.Len(t, role.Permissions, 2)
	assert.Contains(t, role.Permissions, perm1)
	assert.NotContains(t, role.Permissions, perm2)
	assert.Contains(t, role.Permissions, perm3)
}

// Helper functions to create test roles
func createTestRole() *Role {
	return NewRole(
		valueobject.RoleID("test-role-123"),
		"test_role",
		"Test Role",
		"Test role description",
		false, // not a system role
	)
}

func createTestSystemRole() *Role {
	return NewRole(
		valueobject.RoleID("system-role-123"),
		"system_admin",
		"System Admin",
		"System administrator role",
		true, // system role
	)
}
