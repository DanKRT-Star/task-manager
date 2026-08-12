package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMilestoneService_CreateMilestone(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		isMember    bool
		expectError bool
	}{
		{name: "success", title: "v1.0 Release", isMember: true, expectError: false},
		{name: "fail - empty title", title: "", isMember: true, expectError: true},
		{name: "fail - not a project member", title: "v1.0", isMember: false, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMilestoneRepo := new(MockMilestoneRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.isMember {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
			} else {
				mockMemberRepo.On("FindMember", uint(1), uint(1)).Return(nil, assert.AnError)
			}

			if tt.isMember && tt.title != "" {
				mockMilestoneRepo.On("Create", mock.AnythingOfType("*model.Milestone")).Return(nil)
			}

			milestoneService := NewMilestoneService(mockMilestoneRepo, mockMemberRepo)
			milestone, err := milestoneService.CreateMilestone(1, 1, tt.title, "desc", nil)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, milestone)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, milestone)
			}
		})
	}
}

func TestMilestoneService_DeleteMilestone_OnlyOwnerCanDelete(t *testing.T) {
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
			mockMilestoneRepo := new(MockMilestoneRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			mockMilestoneRepo.On("FindByID", uint(1)).Return(&model.Milestone{MilestoneID: 1, ProjectID: 10}, nil)
			mockMemberRepo.On("FindMember", uint(10), uint(1)).Return(&model.ProjectMember{Role: tt.role}, nil)

			if !tt.expectError {
				mockMilestoneRepo.On("Delete", uint(1)).Return(nil)
			}

			milestoneService := NewMilestoneService(mockMilestoneRepo, mockMemberRepo)
			err := milestoneService.DeleteMilestone(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}