package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEpicService_CreateEpic(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		isMember    bool
		expectError bool
	}{
		{name: "success", title: "Auth System", isMember: true, expectError: false},
		{name: "fail - empty title", title: "", isMember: true, expectError: true},
		{name: "fail - not a project member", title: "Some Epic", isMember: false, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEpicRepo := new(MockEpicRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.isMember {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
			} else {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(nil, assert.AnError)
			}

			if tt.isMember && tt.title != "" {
				mockEpicRepo.On("Create", mock.AnythingOfType("*model.Epic")).Return(nil)
			}

			epicService := NewEpicService(mockEpicRepo, mockMemberRepo)
			epic, err := epicService.CreateEpic(1, 1, tt.title, "desc")

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, epic)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, epic)
			}
		})
	}
}

func TestEpicService_DeleteEpic_OnlyOwnerCanDelete(t *testing.T) {
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
			mockEpicRepo := new(MockEpicRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			mockEpicRepo.On("FindByID", uint(1)).Return(&model.Epic{EpicID: 1, ProjectID: 10}, nil)
			mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: tt.role}, nil)

			if !tt.expectError {
				mockEpicRepo.On("Delete", uint(1)).Return(nil)
			}

			epicService := NewEpicService(mockEpicRepo, mockMemberRepo)
			err := epicService.DeleteEpic(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEpicService_UpdateEpic(t *testing.T) {
	t.Run("member can update epic title", func(t *testing.T) {
		mockEpicRepo := new(MockEpicRepository)
		mockMemberRepo := new(MockProjectMemberRepository)

		mockEpicRepo.On("FindByID", uint(1)).Return(&model.Epic{EpicID: 1, ProjectID: 10, Title: "Old"}, nil)
		mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
		mockEpicRepo.On("Update", mock.AnythingOfType("*model.Epic")).Return(nil)

		epicService := NewEpicService(mockEpicRepo, mockMemberRepo)
		epic, err := epicService.UpdateEpic(1, 1, "New Title", "")

		assert.NoError(t, err)
		assert.Equal(t, "New Title", epic.Title)
	})

	t.Run("non-member cannot update epic", func(t *testing.T) {
		mockEpicRepo := new(MockEpicRepository)
		mockMemberRepo := new(MockProjectMemberRepository)

		mockEpicRepo.On("FindByID", uint(1)).Return(&model.Epic{EpicID: 1, ProjectID: 10}, nil)
		mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(nil, assert.AnError)

		epicService := NewEpicService(mockEpicRepo, mockMemberRepo)
		epic, err := epicService.UpdateEpic(1, 1, "New Title", "")

		assert.Error(t, err)
		assert.Nil(t, epic)
	})
}