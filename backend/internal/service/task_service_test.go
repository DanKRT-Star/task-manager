package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_ActivityLog(t *testing.T) {
	t.Run("create task logs created action", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
		mockActivityRepo.On("Create", mock.MatchedBy(func(log *model.ActivityLog) bool {
			return log.Action == model.ActionCreated && log.Detail == "Task created"
		})).Return(nil)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Build API", "desc", model.StatusPending, "", nil, nil, nil, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, task)
	})

	t.Run("update task status logs status change", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		existing := &model.Task{TaskID: 1, UserID: 1, ProjectID: nil, Status: model.StatusPending}
		mockRepo.On("FindByIDOnly", uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)
		mockActivityRepo.On("Create", mock.MatchedBy(func(log *model.ActivityLog) bool {
			return log.TaskID == 1 && log.UserID == 1 && log.Action == model.ActionStatusChanged && log.Detail == "Status changed from pending to done"
		})).Return(nil)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		updated, err := taskService.UpdateTask(1, 1, "", "", model.StatusDone, "", nil, nil, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, model.StatusDone, updated.Status)
	})
}

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
			mockEpicRepo := new(MockEpicRepository)
			mockMilestoneRepo := new(MockMilestoneRepository)
			mockSprintRepo := new(MockSprintRepository)
			if !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
				mockActivityRepo.On("Create", mock.AnythingOfType("*model.ActivityLog")).Return(nil)
			}

			taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
			task, err := taskService.CreateTask(1, tt.title, "desc", tt.status, tt.deadline, nil, nil, nil, nil, nil)

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

func TestTaskService_CreateTask_EpicAndMilestone(t *testing.T) {
	projectID := uint(10)
	epicID := uint(5)
	milestoneID := uint(7)

	t.Run("success - epic and milestone belong to the same project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockEpicRepo.On("FindByID", epicID).Return(&model.Epic{EpicID: epicID, ProjectID: projectID}, nil)
		mockMilestoneRepo.On("FindByID", milestoneID).Return(&model.Milestone{MilestoneID: milestoneID, ProjectID: projectID}, nil)
		mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
		mockActivityRepo.On("Create", mock.AnythingOfType("*model.ActivityLog")).Return(nil)

		// sprintID không được truyền (nil) trong test này nên SprintRepo không bị gọi tới —
		// không cần stub mockSprintRepo.On("FindByID", ...).
		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Task with epic", "desc", "", "", &projectID, nil, &epicID, &milestoneID, nil)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, epicID, *task.EpicID)
		assert.Equal(t, milestoneID, *task.MilestoneID)
	})

	t.Run("fail - epic belongs to a different project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		otherProjectID := uint(99)
		mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockEpicRepo.On("FindByID", epicID).Return(&model.Epic{EpicID: epicID, ProjectID: otherProjectID}, nil)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Task with wrong epic", "desc", "", "", &projectID, nil, &epicID, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "epic does not belong")
	})

	t.Run("fail - epic on a task without a project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Personal task with epic", "desc", "", "", nil, nil, &epicID, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "does not belong to a project")
	})
}

func TestTaskService_CreateTask_Sprint(t *testing.T) {
	projectID := uint(10)
	sprintID := uint(3)

	t.Run("success - sprint belongs to the same project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockSprintRepo.On("FindByID", sprintID).Return(&model.Sprint{SprintID: sprintID, ProjectID: projectID}, nil)
		mockRepo.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
		mockActivityRepo.On("Create", mock.AnythingOfType("*model.ActivityLog")).Return(nil)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Task with sprint", "desc", "", "", &projectID, nil, nil, nil, &sprintID)

		assert.NoError(t, err)
		assert.NotNil(t, task)
		assert.Equal(t, sprintID, *task.SprintID)
	})

	t.Run("fail - sprint belongs to a different project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		otherProjectID := uint(99)
		mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
		mockSprintRepo.On("FindByID", sprintID).Return(&model.Sprint{SprintID: sprintID, ProjectID: otherProjectID}, nil)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Task with wrong sprint", "desc", "", "", &projectID, nil, nil, nil, &sprintID)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "sprint does not belong")
	})

	t.Run("fail - sprint on a task without a project", func(t *testing.T) {
		mockRepo := new(MockTaskRepository)
		mockMemberRepo := new(MockProjectMemberRepository)
		mockActivityRepo := new(MockActivityLogRepository)
		mockEpicRepo := new(MockEpicRepository)
		mockMilestoneRepo := new(MockMilestoneRepository)
		mockSprintRepo := new(MockSprintRepository)

		taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
		task, err := taskService.CreateTask(1, "Personal task with sprint", "desc", "", "", nil, nil, nil, nil, &sprintID)

		assert.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "does not belong to a project")
	})
}

func TestTaskService_UpdateTask_AuthorizationDenied(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	// Task thuộc user 1, không nằm trong project nào — user 2 không phải chủ sở hữu
	mockRepo.On("FindByIDOnly", uint(1)).Return(&model.Task{TaskID: 1, UserID: 1, ProjectID: nil}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	task, err := taskService.UpdateTask(1, 2, "New title", "", "", "", nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "permission")
}

func TestTaskService_UpdateTask_EpicMustBelongToSameProject(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	projectID := uint(10)
	otherProjectID := uint(99)
	epicID := uint(5)
	task := &model.Task{TaskID: 1, UserID: 1, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
	mockEpicRepo.On("FindByID", epicID).Return(&model.Epic{EpicID: epicID, ProjectID: otherProjectID}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	updated, err := taskService.UpdateTask(1, 1, "", "", "", "", nil, &epicID, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "epic does not belong")
}

func TestTaskService_UpdateTask_SprintMustBelongToSameProject(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	projectID := uint(10)
	otherProjectID := uint(99)
	sprintID := uint(3)
	task := &model.Task{TaskID: 1, UserID: 1, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
	mockSprintRepo.On("FindByID", sprintID).Return(&model.Sprint{SprintID: sprintID, ProjectID: otherProjectID}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	updated, err := taskService.UpdateTask(1, 1, "", "", "", "", nil, nil, nil, &sprintID)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "sprint does not belong")
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
			mockEpicRepo := new(MockEpicRepository)
			mockMilestoneRepo := new(MockMilestoneRepository)
			mockSprintRepo := new(MockSprintRepository)
			mockRepo.On("FindByIDOnly", uint(1)).Return(tt.findResult, tt.findErr)

			if !tt.expectError {
				mockRepo.On("Delete", uint(1)).Return(tt.deleteErr)
			}

			taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
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
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	projectID := uint(10)
	task := &model.Task{TaskID: 1, UserID: 2, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	updated, err := taskService.UpdateTask(1, 2, "Updated title", "", "", "", nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated title", updated.Title)
}

func TestTaskService_UpdateTask_ProjectMemberCannotEditOthersTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	projectID := uint(10)
	// task được tạo bởi user 3, không assign cho ai
	task := &model.Task{TaskID: 1, UserID: 3, ProjectID: &projectID, AssigneeID: nil}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(2)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	updated, err := taskService.UpdateTask(1, 2, "Updated title", "", "", "", nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "permission")
}

func TestTaskService_UpdateTask_ProjectOwnerCanEditAnyTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMemberRepo := new(MockProjectMemberRepository)
	mockActivityRepo := new(MockActivityLogRepository)
	mockEpicRepo := new(MockEpicRepository)
	mockMilestoneRepo := new(MockMilestoneRepository)
	mockSprintRepo := new(MockSprintRepository)

	projectID := uint(10)
	task := &model.Task{TaskID: 1, UserID: 3, ProjectID: &projectID}

	mockRepo.On("FindByIDOnly", uint(1)).Return(task, nil)
	mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleOwner}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)

	taskService := NewTaskService(mockRepo, mockMemberRepo, mockActivityRepo, mockEpicRepo, mockMilestoneRepo, mockSprintRepo)
	updated, err := taskService.UpdateTask(1, 1, "Owner edited this", "", "", "", nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
}