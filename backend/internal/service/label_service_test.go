package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLabelService_CreateLabel(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		projectID   uint
		labelName   string
		color       string
		isMember    bool
		expectError bool
	}{
		{name: "success", userID: 1, projectID: 10, labelName: "bug", color: "#ff0000", isMember: true, expectError: false},
		{name: "fail - empty name", userID: 1, projectID: 10, labelName: "", color: "#ff0000", isMember: true, expectError: true},
		{name: "fail - not a member", userID: 2, projectID: 10, labelName: "urgent", color: "#00ff00", isMember: false, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(MockLabelRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)

			if tt.isMember {
				mockMemberRepo.On("FindMember", tt.projectID, tt.userID).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
				if !tt.expectError {
					mockLabelRepo.On("Create", mock.AnythingOfType("*model.Label")).Return(nil)
				}
			} else {
				mockMemberRepo.On("FindMember", tt.projectID, tt.userID).Return(nil, assert.AnError)
			}

			service := NewLabelService(mockLabelRepo, mockTaskRepo, mockMemberRepo, mockActivityRepo)
			label, err := service.CreateLabel(tt.userID, tt.projectID, tt.labelName, tt.color)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, label)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, label)
				assert.Equal(t, tt.labelName, label.Name)
			}
		})
	}
}

func TestLabelService_GetProjectLabels(t *testing.T) {
	tests := []struct {
		name        string
		isMember    bool
		expectError bool
	}{
		{name: "success", isMember: true, expectError: false},
		{name: "fail - not a member", isMember: false, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(MockLabelRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)

			if tt.isMember {
				mockMemberRepo.On("FindMember", uint(11), uint(5)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
				mockLabelRepo.On("FindAllByProject", uint(11)).Return([]model.Label{{LabelID: 1, ProjectID: 11, Name: "bug"}}, nil)
			} else {
				mockMemberRepo.On("FindMember", uint(11), uint(5)).Return(nil, assert.AnError)
			}

			service := NewLabelService(mockLabelRepo, mockTaskRepo, mockMemberRepo, mockActivityRepo)
			labels, err := service.GetProjectLabels(5, 11)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, labels)
			} else {
				assert.NoError(t, err)
				assert.Len(t, labels, 1)
			}
		})
	}
}

func TestLabelService_DeleteLabel_OnlyOwnerCanDelete(t *testing.T) {
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
			mockLabelRepo := new(MockLabelRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)

			mockLabelRepo.On("FindByID", uint(7)).Return(&model.Label{LabelID: 7, ProjectID: 99}, nil)
			mockMemberRepo.On("FindMember", uint(99), uint(3)).Return(&model.ProjectMember{Role: tt.role}, nil)
			if !tt.expectError {
				mockLabelRepo.On("Delete", uint(7)).Return(nil)
			}

			service := NewLabelService(mockLabelRepo, mockTaskRepo, mockMemberRepo, mockActivityRepo)
			err := service.DeleteLabel(3, 7)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLabelService_AttachAndDetachLabel(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		taskID       uint
		labelID      uint
		taskProject  *uint
		memberRole   model.ProjectRole
		labelProject uint
		isTaskFound  bool
		isLabelFound bool
		expectError  bool
	}{
		{name: "attach success", userID: 1, taskID: 5, labelID: 7, taskProject: func() *uint { v := uint(12); return &v }(), memberRole: model.RoleMember, labelProject: 12, isTaskFound: true, isLabelFound: true, expectError: false},
		{name: "fail - task not found", userID: 1, taskID: 5, labelID: 7, taskProject: nil, memberRole: model.RoleMember, labelProject: 0, isTaskFound: false, isLabelFound: false, expectError: true},
		{name: "fail - label mismatch", userID: 1, taskID: 5, labelID: 7, taskProject: func() *uint { v := uint(12); return &v }(), memberRole: model.RoleMember, labelProject: 13, isTaskFound: true, isLabelFound: true, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(MockLabelRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)

			if tt.isTaskFound {
				mockTaskRepo.On("FindByIDOnly", tt.taskID).Return(&model.Task{TaskID: tt.taskID, ProjectID: tt.taskProject}, nil)
				if tt.taskProject != nil {
					mockMemberRepo.On("FindMember", *tt.taskProject, tt.userID).Return(&model.ProjectMember{Role: tt.memberRole}, nil)
					if tt.isLabelFound && tt.labelProject == *tt.taskProject && !tt.expectError {
						mockLabelRepo.On("FindByID", tt.labelID).Return(&model.Label{LabelID: tt.labelID, ProjectID: *tt.taskProject}, nil)
						mockLabelRepo.On("AttachToTask", tt.taskID, tt.labelID).Return(nil)
						mockActivityRepo.On("Create", mock.AnythingOfType("*model.ActivityLog")).Return(nil)
					} else if tt.isLabelFound {
						mockLabelRepo.On("FindByID", tt.labelID).Return(&model.Label{LabelID: tt.labelID, ProjectID: tt.labelProject}, nil)
					}
				}
			} else {
				mockTaskRepo.On("FindByIDOnly", tt.taskID).Return(nil, assert.AnError)
			}

			service := NewLabelService(mockLabelRepo, mockTaskRepo, mockMemberRepo, mockActivityRepo)
			err := service.AttachLabel(tt.userID, tt.taskID, tt.labelID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLabelService_GetTaskLabels(t *testing.T) {
	tests := []struct {
		name        string
		projectID   *uint
		isMember    bool
		expectError bool
	}{
		{name: "success", projectID: func() *uint { v := uint(15); return &v }(), isMember: true, expectError: false},
		{name: "fail - no access", projectID: func() *uint { v := uint(15); return &v }(), isMember: false, expectError: true},
		{name: "task without project", projectID: nil, isMember: true, expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLabelRepo := new(MockLabelRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)

			if tt.projectID == nil {
				mockTaskRepo.On("FindByIDOnly", uint(9)).Return(&model.Task{TaskID: 9, ProjectID: nil}, nil)
				mockLabelRepo.On("FindLabelsByTask", uint(9)).Return([]model.Label{}, nil)
			} else {
				mockTaskRepo.On("FindByIDOnly", uint(9)).Return(&model.Task{TaskID: 9, ProjectID: tt.projectID}, nil)
				if tt.isMember {
					mockMemberRepo.On("FindMember", *tt.projectID, uint(4)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
					mockLabelRepo.On("FindLabelsByTask", uint(9)).Return([]model.Label{{LabelID: 3, ProjectID: *tt.projectID, Name: "done"}}, nil)
				} else {
					mockMemberRepo.On("FindMember", *tt.projectID, uint(4)).Return(nil, assert.AnError)
				}
			}

			service := NewLabelService(mockLabelRepo, mockTaskRepo, mockMemberRepo, mockActivityRepo)
			labels, err := service.GetTaskLabels(4, 9)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, labels)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}