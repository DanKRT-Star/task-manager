package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSprintService_CreateSprint(t *testing.T) {
	tests := []struct {
		name        string
		sprintName  string
		isMember    bool
		expectError bool
	}{
		{name: "success", sprintName: "Sprint 1", isMember: true, expectError: false},
		{name: "fail - empty name", sprintName: "", isMember: true, expectError: true},
		{name: "fail - not a project member", sprintName: "Sprint X", isMember: false, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSprintRepo := new(MockSprintRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.isMember {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
			} else {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(nil, assert.AnError)
			}

			if tt.isMember && tt.sprintName != "" {
				mockSprintRepo.On("Create", mock.AnythingOfType("*model.Sprint")).Return(nil)
			}

			sprintService := NewSprintService(mockSprintRepo, mockMemberRepo)
			sprint, err := sprintService.CreateSprint(1, 1, tt.sprintName, nil, nil)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, sprint)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, sprint)
				assert.Equal(t, model.SprintPlanned, sprint.Status)
			}
		})
	}
}

func TestSprintService_UpdateSprint(t *testing.T) {
	t.Run("member can update sprint status to active", func(t *testing.T) {
		mockSprintRepo := new(MockSprintRepository)
		mockMemberRepo := new(MockProjectMemberRepository)

		mockSprintRepo.On("FindByID", uint(1)).Return(&model.Sprint{SprintID: 1, ProjectID: 10, Status: model.SprintPlanned}, nil)
		mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
		mockSprintRepo.On("Update", mock.AnythingOfType("*model.Sprint")).Return(nil)

		sprintService := NewSprintService(mockSprintRepo, mockMemberRepo)
		sprint, err := sprintService.UpdateSprint(1, 1, "", model.SprintActive, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, model.SprintActive, sprint.Status)
	})

	t.Run("fail - invalid status value", func(t *testing.T) {
		mockSprintRepo := new(MockSprintRepository)
		mockMemberRepo := new(MockProjectMemberRepository)

		mockSprintRepo.On("FindByID", uint(1)).Return(&model.Sprint{SprintID: 1, ProjectID: 10}, nil)
		mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)

		sprintService := NewSprintService(mockSprintRepo, mockMemberRepo)
		sprint, err := sprintService.UpdateSprint(1, 1, "", model.SprintStatus("bogus"), nil, nil)

		assert.Error(t, err)
		assert.Nil(t, sprint)
	})

	t.Run("non-member cannot update sprint", func(t *testing.T) {
		mockSprintRepo := new(MockSprintRepository)
		mockMemberRepo := new(MockProjectMemberRepository)

		mockSprintRepo.On("FindByID", uint(1)).Return(&model.Sprint{SprintID: 1, ProjectID: 10}, nil)
		mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(nil, assert.AnError)

		sprintService := NewSprintService(mockSprintRepo, mockMemberRepo)
		sprint, err := sprintService.UpdateSprint(1, 1, "New name", "", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, sprint)
	})
}

func TestSprintService_DeleteSprint_OnlyOwnerCanDelete(t *testing.T) {
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
			mockSprintRepo := new(MockSprintRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			mockSprintRepo.On("FindByID", uint(1)).Return(&model.Sprint{SprintID: 1, ProjectID: 10}, nil)
			mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: tt.role}, nil)

			if !tt.expectError {
				mockSprintRepo.On("Delete", uint(1)).Return(nil)
			}

			sprintService := NewSprintService(mockSprintRepo, mockMemberRepo)
			err := sprintService.DeleteSprint(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}