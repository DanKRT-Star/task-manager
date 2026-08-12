package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestActivityLogService_GetTaskActivity(t *testing.T) {
	tests := []struct {
		name        string
		taskExists  bool
		projectID   *uint
		taskOwnerID uint
		requesterID uint
		isMember    bool
		expectError bool
	}{
		{
			name:        "success - personal task owner can view",
			taskExists:  true,
			projectID:   nil,
			taskOwnerID: 1,
			requesterID: 1,
			expectError: false,
		},
		{
			name:        "fail - personal task, not the owner",
			taskExists:  true,
			projectID:   nil,
			taskOwnerID: 1,
			requesterID: 2,
			expectError: true,
		},
		{
			name:        "success - project task, requester is a member",
			taskExists:  true,
			projectID:   func() *uint { v := uint(10); return &v }(),
			taskOwnerID: 1,
			requesterID: 2,
			isMember:    true,
			expectError: false,
		},
		{
			name:        "fail - project task, requester is not a member",
			taskExists:  true,
			projectID:   func() *uint { v := uint(10); return &v }(),
			taskOwnerID: 1,
			requesterID: 3,
			isMember:    false,
			expectError: true,
		},
		{
			name:        "fail - task not found",
			taskExists:  false,
			requesterID: 1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivityRepo := new(MockActivityLogRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.taskExists {
				task := &model.Task{TaskID: 5, ProjectID: tt.projectID, UserID: tt.taskOwnerID}
				mockTaskRepo.On("FindByIDOnly", uint(5)).Return(task, nil)

				if tt.projectID != nil {
					if tt.isMember {
						mockMemberRepo.On("FindMember", *tt.projectID, tt.requesterID).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
					} else {
						mockMemberRepo.On("FindMember", *tt.projectID, tt.requesterID).Return(nil, assert.AnError)
					}
				}

				if !tt.expectError {
					mockActivityRepo.On("FindAllByTask", uint(5)).Return([]model.ActivityLog{
						{ActivityLogID: 1, TaskID: 5, UserID: tt.taskOwnerID, Action: model.ActionCreated, Detail: "Task created"},
					}, nil)
				}
			} else {
				mockTaskRepo.On("FindByIDOnly", uint(5)).Return(nil, assert.AnError)
			}

			service := NewActivityLogService(mockActivityRepo, mockTaskRepo, mockMemberRepo)
			logs, err := service.GetTaskActivity(tt.requesterID, 5)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, logs)
			} else {
				assert.NoError(t, err)
				assert.Len(t, logs, 1)
			}
		})
	}
}