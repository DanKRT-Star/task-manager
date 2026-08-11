package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		expectError bool
	}{
		{name: "success", projectName: "New Project", expectError: false},
		{name: "fail - empty name", projectName: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockProjectRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockUserRepo := new(MockUserRepository)

			if !tt.expectError {
				mockProjectRepo.On("Create", mock.AnythingOfType("*model.Project")).Return(nil)
				mockMemberRepo.On("AddMember", mock.AnythingOfType("*model.ProjectMember")).Return(nil)
			}

			projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
			project, err := projectService.CreateProject(1, tt.projectName, "desc", nil)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, project)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, project)
				assert.Equal(t, uint(1), project.OwnerID)
			}
		})
	}
}

func TestProjectService_UpdateProject_OnlyOwnerCanUpdate(t *testing.T) {
	tests := []struct {
		name        string
		role        model.ProjectRole
		expectError bool
	}{
		{name: "owner can update", role: model.RoleOwner, expectError: false},
		{name: "member cannot update", role: model.RoleMember, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockProjectRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockUserRepo := new(MockUserRepository)

			mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: tt.role}, nil)

			if !tt.expectError {
				mockProjectRepo.On("FindByID", uint(1)).Return(&model.Project{ProjectID: 1, Name: "Old"}, nil)
				mockProjectRepo.On("Update", mock.AnythingOfType("*model.Project")).Return(nil)
			}

			projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
			project, err := projectService.UpdateProject(1, 1, "New Name", "", nil)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, project)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "New Name", project.Name)
			}
		})
	}
}

func TestProjectService_DeleteProject_OnlyOwnerCanDelete(t *testing.T) {
	tests := []struct {
		name        string
		role        model.ProjectRole
		expectError bool
	}{
		{name: "owner can delete", role: model.RoleOwner, expectError: false},
		{name: "member cannot delete", role: model.RoleMember, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectRepo := new(MockProjectRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockUserRepo := new(MockUserRepository)

			mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: tt.role}, nil)

			if !tt.expectError {
				mockProjectRepo.On("Delete", uint(1)).Return(nil)
			}

			projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
			err := projectService.DeleteProject(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_AddMember(t *testing.T) {
	t.Run("owner can add a new member", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockUserRepo := new(MockUserRepository)

		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockUserRepo.On("FindByEmail", "new@example.com").Return(&model.User{UserID: 2, Email: "new@example.com"}, nil)
		mockMemberRepo.On("FindMember", uint(1), uint(2)).Return(nil, assert.AnError)
		mockMemberRepo.On("AddMember", mock.AnythingOfType("*model.ProjectMember")).Return(nil)

		projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
		member, err := projectService.AddMember(1, 1, "new@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, member)
		assert.Equal(t, model.RoleMember, member.Role)
	})

	t.Run("non-owner cannot add member", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockUserRepo := new(MockUserRepository)

		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)

		projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
		member, err := projectService.AddMember(1, 1, "new@example.com")

		assert.Error(t, err)
		assert.Nil(t, member)
	})

	t.Run("cannot add already existing member", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockUserRepo := new(MockUserRepository)

		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockUserRepo.On("FindByEmail", "existing@example.com").Return(&model.User{UserID: 2}, nil)
		mockMemberRepo.On("FindMember", uint(1), uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)

		projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
		member, err := projectService.AddMember(1, 1, "existing@example.com")

		assert.Error(t, err)
		assert.Nil(t, member)
		assert.Contains(t, err.Error(), "already a member")
	})
}

func TestProjectService_RemoveMember(t *testing.T) {
	t.Run("owner can remove a regular member", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockUserRepo := new(MockUserRepository)

		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockMemberRepo.On("FindMember", uint(1), uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
		mockMemberRepo.On("RemoveMember", uint(1), uint(2)).Return(nil)

		projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
		err := projectService.RemoveMember(1, 1, 2)

		assert.NoError(t, err)
	})

	t.Run("cannot remove the project owner", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockUserRepo := new(MockUserRepository)

		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)

		projectService := NewProjectService(mockProjectRepo, mockMemberRepo, mockUserRepo)
		err := projectService.RemoveMember(1, 1, 1)

		assert.Error(t, err)
	})
}