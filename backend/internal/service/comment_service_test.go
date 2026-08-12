package service

import (
	"testing"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentService_CreateComment(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		isMember    bool
		taskExists  bool
		expectError bool
	}{
		{name: "success", content: "Good job", isMember: true, taskExists: true, expectError: false},
		{name: "fail - empty content", content: "", isMember: true, taskExists: true, expectError: true},
		{name: "fail - task not found", content: "No task", isMember: true, taskExists: false, expectError: true},
		{name: "fail - no access", content: "Cannot access", isMember: false, taskExists: true, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.taskExists {
				projectID := uint(10)
				task := &model.Task{TaskID: 5, ProjectID: &projectID, UserID: 1}
				mockTaskRepo.On("FindByIDOnly", uint(5)).Return(task, nil)
				if tt.isMember {
					mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
				} else {
					mockMemberRepo.On("FindMember", projectID, uint(1)).Return(nil, assert.AnError)
				}
				if !tt.expectError {
					mockCommentRepo.On("Create", mock.AnythingOfType("*model.Comment")).Return(nil)
				}
			} else {
				mockTaskRepo.On("FindByIDOnly", uint(5)).Return(nil, assert.AnError)
			}

			service := NewCommentService(mockCommentRepo, mockTaskRepo, mockMemberRepo)
			comment, err := service.CreateComment(1, 5, tt.content)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
				assert.Equal(t, uint(5), comment.TaskID)
				assert.Equal(t, uint(1), comment.UserID)
			}
		})
	}
}

func TestCommentService_GetTaskComments(t *testing.T) {
	tests := []struct {
		name        string
		isMember    bool
		taskExists  bool
		expectError bool
	}{
		{name: "success", isMember: true, taskExists: true, expectError: false},
		{name: "fail - task not found", isMember: true, taskExists: false, expectError: true},
		{name: "fail - no access", isMember: false, taskExists: true, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			if tt.taskExists {
				projectID := uint(20)
				task := &model.Task{TaskID: 8, ProjectID: &projectID, UserID: 1}
				mockTaskRepo.On("FindByIDOnly", uint(8)).Return(task, nil)
				if tt.isMember {
					mockMemberRepo.On("FindMember", projectID, uint(1)).Return(&model.ProjectMember{Role: model.RoleMember}, nil)
					mockCommentRepo.On("FindAllByTask", uint(8)).Return([]model.Comment{{CommentID: 1, TaskID: 8, UserID: 1, Content: "hello"}}, nil)
				} else {
					mockMemberRepo.On("FindMember", projectID, uint(1)).Return(nil, assert.AnError)
				}
			} else {
				mockTaskRepo.On("FindByIDOnly", uint(8)).Return(nil, assert.AnError)
			}

			service := NewCommentService(mockCommentRepo, mockTaskRepo, mockMemberRepo)
			comments, err := service.GetTaskComments(1, 8)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, comments)
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, 1)
			}
		})
	}
}

func TestCommentService_DeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		ownerID     uint
		commentUser uint
		expectError bool
	}{
		{name: "owner can delete", ownerID: 1, commentUser: 1, expectError: false},
		{name: "other user cannot delete", ownerID: 2, commentUser: 1, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommentRepo := new(MockCommentRepository)
			mockTaskRepo := new(MockTaskRepository)
			mockMemberRepo := new(MockProjectMemberRepository)

			mockCommentRepo.On("FindByID", uint(1)).Return(&model.Comment{CommentID: 1, TaskID: 2, UserID: tt.commentUser}, nil)
			if !tt.expectError {
				mockCommentRepo.On("Delete", uint(1)).Return(nil)
			}

			service := NewCommentService(mockCommentRepo, mockTaskRepo, mockMemberRepo)
			err := service.DeleteComment(tt.ownerID, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
