package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		status      model.TaskStatus
		deadline    string
		expectError bool
	}{
		{
			name:        "success - default status",
			title:       "Học Golang",
			status:      "",
			deadline:    "2026-08-15T17:00:00Z",
			expectError: false,
		},
		{
			name:        "success - explicit valid status",
			title:       "Học Fiber",
			status:      model.StatusInProgress,
			deadline:    "2026-08-15T17:00:00Z",
			expectError: false,
		},
		{
			name:        "fail - empty title",
			title:       "",
			status:      "",
			deadline:    "2026-08-15T17:00:00Z",
			expectError: true,
		},
		{
			name:        "fail - invalid status",
			title:       "Task X",
			status:      model.TaskStatus("not_real"),
			deadline:    "2026-08-15T17:00:00Z",
			expectError: true,
		},
		{
			name:        "fail - invalid deadline format",
			title:       "Task Y",
			status:      "",
			deadline:    "15/08/2026",
			expectError: true,
		},
		{
			name:        "success - empty deadline allowed",
			title:       "Task Z",
			status:      "",
			deadline:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			if !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
			}

			taskService := NewTaskService(mockRepo)
			task, err := taskService.CreateTask(1, tt.title, "desc", tt.status, tt.deadline)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				if tt.status == "" {
					assert.Equal(t, model.StatusPending, task.Status)
				}
			}
		})
	}
}

func TestTaskService_UpdateTask_AuthorizationDenied(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	// Giả lập: FindByID không tìm thấy task thuộc userID này (task của người khác)
	mockRepo.On("FindByID", uint(1), uint(2)).Return(nil, assert.AnError)

	taskService := NewTaskService(mockRepo)
	task, err := taskService.UpdateTask(1, 2, "New title", "", "", "")

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "access denied")
}

func TestTaskService_DeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		findResult  *model.Task
		findErr     error
		deleteErr   error
		expectError bool
	}{
		{
			name:        "success",
			findResult:  &model.Task{TaskID: 1, UserID: 1},
			findErr:     nil,
			deleteErr:   nil,
			expectError: false,
		},
		{
			name:        "fail - task not found or not owned",
			findResult:  nil,
			findErr:     assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("FindByID", uint(1), uint(1)).Return(tt.findResult, tt.findErr)

			if !tt.expectError {
				mockRepo.On("Delete", uint(1), uint(1)).Return(tt.deleteErr)
			}

			taskService := NewTaskService(mockRepo)
			err := taskService.DeleteTask(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}