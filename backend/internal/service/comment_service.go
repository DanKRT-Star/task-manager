package service

import (
	"errors"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type CommentService struct {
	CommentRepo repository.CommentRepositoryInterface
	TaskRepo    repository.TaskRepositoryInterface
	MemberRepo  repository.ProjectMemberRepositoryInterface
}

func NewCommentService(commentRepo repository.CommentRepositoryInterface, taskRepo repository.TaskRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *CommentService {
	return &CommentService{CommentRepo: commentRepo, TaskRepo: taskRepo, MemberRepo: memberRepo}
}

func (s *CommentService) canAccessTask(userID uint, task *model.Task) bool {
	if task.ProjectID == nil {
		return task.UserID == userID
	}
	_, err := s.MemberRepo.FindMember(*task.ProjectID, userID)
	return err == nil
}

func (s *CommentService) CreateComment(userID, taskID uint, content string) (*model.Comment, error) {
	if content == "" {
		logger.CommentEmptyContent(userID, taskID)
		return nil, errors.New("content is required")
	}

	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.CommentTaskNotFound(taskID, userID)
		return nil, errors.New("task not found")
	}
	if !s.canAccessTask(userID, task) {
		logger.CommentTaskAccessDenied(userID, taskID)
		return nil, errors.New("you do not have access to this task")
	}

	comment := &model.Comment{TaskID: taskID, UserID: userID, Content: content}
	if err := s.CommentRepo.Create(comment); err != nil {
		logger.CommentCreateFailed(taskID, userID, err)
		return nil, err
	}
	logger.CommentCreated(comment.CommentID, taskID, userID)
	return comment, nil
}

func (s *CommentService) GetTaskComments(userID, taskID uint) ([]model.Comment, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.CommentTaskNotFound(taskID, userID)
		return nil, errors.New("task not found")
	}
	if !s.canAccessTask(userID, task) {
		logger.CommentTaskAccessDenied(userID, taskID)
		return nil, errors.New("you do not have access to this task")
	}
	comments, err := s.CommentRepo.FindAllByTask(taskID)
	if err != nil {
		logger.CommentFetchFailed(taskID, err)
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) DeleteComment(userID, commentID uint) error {
	comment, err := s.CommentRepo.FindByID(commentID)
	if err != nil {
		logger.CommentNotFound(commentID, userID)
		return errors.New("comment not found")
	}
	if comment.UserID != userID {
		logger.CommentDeleteForbidden(commentID, userID, comment.UserID)
		return errors.New("you can only delete your own comments")
	}
	if err := s.CommentRepo.Delete(commentID); err != nil {
		logger.CommentDeleteFailed(commentID, userID, err)
		return err
	}
	logger.CommentDeleted(commentID, userID)
	return nil
}