package aggregate

import (
	"testing"

	"github.com/taskflow/internal/domain/valueobject"
)

func TestNewProject(t *testing.T) {
	// Arrange
	id := valueobject.ProjectID("test-project-1")
	name := "Test Project"
	description := "Test Description"
	projectType := valueobject.ProjectTypeMaster
	ownerID := valueobject.UserID("user-1")

	// Act
	project := NewProject(id, name, description, projectType, ownerID)

	// Assert
	if project.ID != id {
		t.Errorf("Expected ID %s, got %s", id, project.ID)
	}
	if project.Name != name {
		t.Errorf("Expected Name %s, got %s", name, project.Name)
	}
	if project.Status != valueobject.ProjectStatusDraft {
		t.Errorf("Expected Status %s, got %s", valueobject.ProjectStatusDraft, project.Status)
	}
	if project.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, project.OwnerID)
	}
	if len(project.Events) == 0 {
		t.Error("Expected at least one domain event")
	}
}

func TestProject_UpdateBasicInfo(t *testing.T) {
	// Arrange
	project := createTestProject()
	newName := "Updated Project Name"
	newDescription := "Updated Description"

	// Act
	err := project.UpdateBasicInfo(newName, newDescription)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.Name != newName {
		t.Errorf("Expected Name %s, got %s", newName, project.Name)
	}
	if project.Description != newDescription {
		t.Errorf("Expected Description %s, got %s", newDescription, project.Description)
	}
}

func TestProject_UpdateBasicInfo_EmptyName(t *testing.T) {
	// Arrange
	project := createTestProject()

	// Act
	err := project.UpdateBasicInfo("", "New Description")

	// Assert
	if err == nil {
		t.Error("Expected error for empty name")
	}
}

func TestProject_AssignManager(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	assignedBy := project.OwnerID

	// Act
	err := project.AssignManager(managerID, assignedBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.ManagerID == nil || *project.ManagerID != managerID {
		t.Errorf("Expected ManagerID %s, got %v", managerID, project.ManagerID)
	}
	// Manager should be automatically added to members
	if !project.isMember(managerID) {
		t.Error("Manager should be automatically added to members")
	}
}

func TestProject_AssignManager_NotOwner(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	assignedBy := valueobject.UserID("not-owner")

	// Act
	err := project.AssignManager(managerID, assignedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when non-owner tries to assign manager")
	}
}

func TestProject_AssignManager_OwnerAsManager(t *testing.T) {
	// Arrange
	project := createTestProject()
	assignedBy := project.OwnerID

	// Act
	err := project.AssignManager(project.OwnerID, assignedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when trying to assign owner as manager")
	}
}

func TestProject_AddMember(t *testing.T) {
	// Arrange
	project := createTestProject()
	userID := valueobject.UserID("member-1")
	role := valueobject.ProjectRoleDeveloper
	addedBy := project.OwnerID

	// Act
	err := project.AddMember(userID, role, addedBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !project.isMember(userID) {
		t.Error("User should be added to members")
	}
	memberRole := project.GetMemberRole(userID)
	if memberRole == nil || *memberRole != role {
		t.Errorf("Expected role %s, got %v", role, memberRole)
	}
}

func TestProject_AddMember_AlreadyMember(t *testing.T) {
	// Arrange
	project := createTestProject()
	userID := valueobject.UserID("member-1")
	role := valueobject.ProjectRoleDeveloper
	addedBy := project.OwnerID

	// Add member first time
	project.AddMember(userID, role, addedBy)

	// Act - try to add again
	err := project.AddMember(userID, role, addedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when adding existing member")
	}
}

func TestProject_AddMember_InsufficientPermission(t *testing.T) {
	// Arrange
	project := createTestProject()
	userID := valueobject.UserID("member-1")
	role := valueobject.ProjectRoleDeveloper
	addedBy := valueobject.UserID("unauthorized-user")

	// Act
	err := project.AddMember(userID, role, addedBy)

	// Assert
	if err == nil {
		t.Error("Expected error for insufficient permission")
	}
}

func TestProject_RemoveMember(t *testing.T) {
	// Arrange
	project := createTestProject()
	userID := valueobject.UserID("member-1")
	role := valueobject.ProjectRoleDeveloper
	addedBy := project.OwnerID

	// Add member first
	project.AddMember(userID, role, addedBy)

	// Act
	err := project.RemoveMember(userID, project.OwnerID)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.isMember(userID) {
		t.Error("User should be removed from members")
	}
}

func TestProject_RemoveMember_CannotRemoveOwner(t *testing.T) {
	// Arrange
	project := createTestProject()

	// Act
	err := project.RemoveMember(project.OwnerID, project.OwnerID)

	// Assert
	if err == nil {
		t.Error("Expected error when trying to remove owner")
	}
}

func TestProject_RemoveMember_CannotRemoveManager(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	project.AssignManager(managerID, project.OwnerID)

	// Act
	err := project.RemoveMember(managerID, project.OwnerID)

	// Assert
	if err == nil {
		t.Error("Expected error when trying to remove manager")
	}
}

func TestProject_CreateSubProject(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Status = valueobject.ProjectStatusActive // Must be active to create sub projects

	subProjectID := valueobject.ProjectID("sub-project-1")
	subName := "Sub Project"
	subDescription := "Sub Description"
	createdBy := project.OwnerID

	// Act
	subProject, err := project.CreateSubProject(subProjectID, subName, subDescription, createdBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if subProject == nil {
		t.Fatal("Expected sub project to be created")
	}

	// Check sub project properties
	subProjectConcrete := subProject.(*Project)
	if subProjectConcrete.ProjectType != valueobject.ProjectTypeSub {
		t.Errorf("Expected sub project type %s, got %s", valueobject.ProjectTypeSub, subProjectConcrete.ProjectType)
	}
	if subProjectConcrete.ParentID == nil || *subProjectConcrete.ParentID != project.ID {
		t.Errorf("Expected parent ID %s, got %v", project.ID, subProjectConcrete.ParentID)
	}

	// Check parent project children list
	found := false
	for _, childID := range project.Children {
		if childID == subProjectID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Sub project ID should be added to parent's children list")
	}
}

func TestProject_CreateSubProject_NotMasterProject(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.ProjectType = valueobject.ProjectTypeSub // Sub projects cannot have children
	project.Status = valueobject.ProjectStatusActive

	subProjectID := valueobject.ProjectID("sub-project-1")
	subName := "Sub Project"
	subDescription := "Sub Description"
	createdBy := project.OwnerID

	// Act
	_, err := project.CreateSubProject(subProjectID, subName, subDescription, createdBy)

	// Assert
	if err == nil {
		t.Error("Expected error when non-master project tries to create sub project")
	}
}

func TestProject_CreateSubProject_NotActive(t *testing.T) {
	// Arrange
	project := createTestProject()
	// Project is in draft status by default

	subProjectID := valueobject.ProjectID("sub-project-1")
	subName := "Sub Project"
	subDescription := "Sub Description"
	createdBy := project.OwnerID

	// Act
	_, err := project.CreateSubProject(subProjectID, subName, subDescription, createdBy)

	// Assert
	if err == nil {
		t.Error("Expected error when inactive project tries to create sub project")
	}
}

func TestProject_Activate(t *testing.T) {
	// Arrange
	project := createTestProject()
	activatedBy := project.OwnerID

	// Act
	err := project.Activate(activatedBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.Status != valueobject.ProjectStatusActive {
		t.Errorf("Expected status %s, got %s", valueobject.ProjectStatusActive, project.Status)
	}
	if project.StartDate.IsZero() {
		t.Error("Start date should be set when activating")
	}
}

func TestProject_Activate_AlreadyActive(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Status = valueobject.ProjectStatusActive
	activatedBy := project.OwnerID

	// Act
	err := project.Activate(activatedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when activating already active project")
	}
}

func TestProject_Complete(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Status = valueobject.ProjectStatusActive
	project.TaskCount = 5
	project.CompletedTasks = 5 // All tasks completed
	completedBy := project.OwnerID

	// Act
	err := project.Complete(completedBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.Status != valueobject.ProjectStatusCompleted {
		t.Errorf("Expected status %s, got %s", valueobject.ProjectStatusCompleted, project.Status)
	}
	if project.EndDate == nil {
		t.Error("End date should be set when completing")
	}
}

func TestProject_Complete_WithPendingTasks(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Status = valueobject.ProjectStatusActive
	project.TaskCount = 5
	project.CompletedTasks = 3 // Still has pending tasks
	completedBy := project.OwnerID

	// Act
	err := project.Complete(completedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when completing project with pending tasks")
	}
}

func TestProject_Cancel(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Status = valueobject.ProjectStatusActive
	cancelledBy := project.OwnerID
	reason := "Project no longer needed"

	// Act
	err := project.Cancel(cancelledBy, reason)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.Status != valueobject.ProjectStatusCancelled {
		t.Errorf("Expected status %s, got %s", valueobject.ProjectStatusCancelled, project.Status)
	}
	if project.EndDate == nil {
		t.Error("End date should be set when cancelling")
	}
}

func TestProject_Delete(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.TaskCount = 5
	project.CompletedTasks = 5 // All tasks completed
	deletedBy := project.OwnerID

	// Act
	err := project.Delete(deletedBy)

	// Assert
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if project.DeletedAt == nil {
		t.Error("DeletedAt should be set when deleting")
	}
}

func TestProject_Delete_WithPendingTasks(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.TaskCount = 5
	project.CompletedTasks = 3 // Still has pending tasks
	deletedBy := project.OwnerID

	// Act
	err := project.Delete(deletedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when deleting project with pending tasks")
	}
}

func TestProject_Delete_WithSubProjects(t *testing.T) {
	// Arrange
	project := createTestProject()
	project.Children = []valueobject.ProjectID{"sub-1", "sub-2"} // Has sub projects
	deletedBy := project.OwnerID

	// Act
	err := project.Delete(deletedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when deleting project with sub projects")
	}
}

func TestProject_Delete_NotOwner(t *testing.T) {
	// Arrange
	project := createTestProject()
	deletedBy := valueobject.UserID("not-owner")

	// Act
	err := project.Delete(deletedBy)

	// Assert
	if err == nil {
		t.Error("Expected error when non-owner tries to delete project")
	}
}

func TestProject_CanUserAccess(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	memberID := valueobject.UserID("member-1")
	outsiderID := valueobject.UserID("outsider")

	// Setup
	project.AssignManager(managerID, project.OwnerID)
	project.AddMember(memberID, valueobject.ProjectRoleDeveloper, project.OwnerID)

	// Test owner access
	if !project.CanUserAccess(project.OwnerID) {
		t.Error("Owner should have access")
	}

	// Test manager access
	if !project.CanUserAccess(managerID) {
		t.Error("Manager should have access")
	}

	// Test member access
	if !project.CanUserAccess(memberID) {
		t.Error("Member should have access")
	}

	// Test outsider access
	if project.CanUserAccess(outsiderID) {
		t.Error("Outsider should not have access")
	}
}

func TestProject_GetMemberRole(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	memberID := valueobject.UserID("member-1")
	outsiderID := valueobject.UserID("outsider")

	// Setup
	project.AssignManager(managerID, project.OwnerID)
	project.AddMember(memberID, valueobject.ProjectRoleDeveloper, project.OwnerID)

	// Test owner role
	ownerRole := project.GetMemberRole(project.OwnerID)
	if ownerRole == nil || *ownerRole != valueobject.ProjectRoleManager {
		t.Errorf("Expected owner role %s, got %v", valueobject.ProjectRoleManager, ownerRole)
	}

	// Test manager role
	managerRole := project.GetMemberRole(managerID)
	if managerRole == nil || *managerRole != valueobject.ProjectRoleManager {
		t.Errorf("Expected manager role %s, got %v", valueobject.ProjectRoleManager, managerRole)
	}

	// Test member role
	memberRole := project.GetMemberRole(memberID)
	if memberRole == nil || *memberRole != valueobject.ProjectRoleDeveloper {
		t.Errorf("Expected member role %s, got %v", valueobject.ProjectRoleDeveloper, memberRole)
	}

	// Test outsider role
	outsiderRole := project.GetMemberRole(outsiderID)
	if outsiderRole != nil {
		t.Errorf("Expected nil role for outsider, got %v", outsiderRole)
	}
}

func TestProject_GetMemberIDs(t *testing.T) {
	// Arrange
	project := createTestProject()
	managerID := valueobject.UserID("manager-1")
	memberID := valueobject.UserID("member-1")

	// Setup
	project.AssignManager(managerID, project.OwnerID)
	project.AddMember(memberID, valueobject.ProjectRoleDeveloper, project.OwnerID)

	// Act
	memberIDs := project.GetMemberIDs()

	// Assert
	expectedCount := 3 // owner + manager + member
	if len(memberIDs) != expectedCount {
		t.Errorf("Expected %d member IDs, got %d", expectedCount, len(memberIDs))
	}

	// Check if all expected IDs are present
	expectedIDs := []string{string(project.OwnerID), string(managerID), string(memberID)}
	for _, expectedID := range expectedIDs {
		found := false
		for _, memberID := range memberIDs {
			if memberID == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected member ID %s not found in result", expectedID)
		}
	}
}

func TestProject_UpdateTaskStatistics(t *testing.T) {
	// Arrange
	project := createTestProject()
	totalTasks := 10
	completedTasks := 7

	// Act
	project.UpdateTaskStatistics(totalTasks, completedTasks)

	// Assert
	if project.TaskCount != totalTasks {
		t.Errorf("Expected TaskCount %d, got %d", totalTasks, project.TaskCount)
	}
	if project.CompletedTasks != completedTasks {
		t.Errorf("Expected CompletedTasks %d, got %d", completedTasks, project.CompletedTasks)
	}
}

// Helper function to create a test project
func createTestProject() *Project {
	return NewProject(
		valueobject.ProjectID("test-project"),
		"Test Project",
		"Test Description",
		valueobject.ProjectTypeMaster,
		valueobject.UserID("owner-1"),
	)
}
