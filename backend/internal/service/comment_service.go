package service

import (
	"errors"

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

// canAccessTask kiểm tra user có quyền xem/bình luận task này không
func (s *CommentService) canAccessTask(userID uint, task *model.Task) bool {
	if task.ProjectID == nil {
		return task.UserID == userID
	}
	_, err := s.MemberRepo.FindMember(*task.ProjectID, userID)
	return err == nil
}

func (s *CommentService) CreateComment(userID, taskID uint, content string) (*model.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}

	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if !s.canAccessTask(userID, task) {
		return nil, errors.New("you do not have access to this task")
	}

	comment := &model.Comment{TaskID: taskID, UserID: userID, Content: content}
	if err := s.CommentRepo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetTaskComments(userID, taskID uint) ([]model.Comment, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if !s.canAccessTask(userID, task) {
		return nil, errors.New("you do not have access to this task")
	}
	return s.CommentRepo.FindAllByTask(taskID)
}

func (s *CommentService) DeleteComment(userID, commentID uint) error {
	comment, err := s.CommentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}
	// Chỉ tác giả bình luận mới xóa được - không phụ thuộc vào role trong project
	if comment.UserID != userID {
		return errors.New("you can only delete your own comments")
	}
	return s.CommentRepo.Delete(commentID)
}