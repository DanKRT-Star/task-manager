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
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)
			if !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
				mockActivityRepo.On("Create", mock.AnythingOfType("*model.ActivityLog")).Return(nil)
			}

			taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
			task, err := taskService.CreateTask(1, tt.title, "desc", tt.status, tt.deadline, nil, nil)

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
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)

	// Task thuộc user 1, không nằm trong project nào — user 2 không phải chủ sở hữu
	mockRepo.On("FindByIDOnly", uint(1)).Return(&model.Task{TaskID: 1, UserID: 1, ProjectID: nil}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
	task, err := taskService.UpdateTask(1, 2, "New title", "", "", "", nil)

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "permission")
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
			findResult:  &model.Task{TaskID: 1, UserID: 1, ProjectID: nil},
			findErr:     nil,
			deleteErr:   nil,
			expectError: false,
		},
		{
			name:        "fail - task not found",
			findResult:  nil,
			findErr:     assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)
			mockActivityRepo := new(MockActivityLogRepository)
			mockRepo.On("FindByIDOnly", uint(1)).Return(tt.findResult, tt.findErr)

			if !tt.expectError {
				mockRepo.On("Delete", uint(1)).Return(tt.deleteErr)
			}

			taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
			err := taskService.DeleteTask(1, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskService_UpdateTask_ProjectMemberCanEditOwnTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)

	projectID := uint(10)
	task := &model.Task{TaskID: 1, UserID: 2, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
	updated, err := taskService.UpdateTask(1, 2, "Updated title", "", "", "", nil)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated title", updated.Title)
}

func TestTaskService_UpdateTask_ProjectMemberCannotEditOthersTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)

	projectID := uint(10)
	// task được tạo bởi user 3, không assign cho ai
	task := &model.Task{TaskID: 1, UserID: 3, ProjectID: &projectID, AssigneeID: nil}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
	updated, err := taskService.UpdateTask(1, 2, "Updated title", "", "", "", nil)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "permission")
}

func TestTaskService_UpdateTask_ProjectOwnerCanEditAnyTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)

	projectID := uint(10)
	task := &model.Task{TaskID: 1, UserID: 3, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo)
	updated, err := taskService.UpdateTask(1, 1, "Owner edited this", "", "", "", nil)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
}
