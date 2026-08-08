package service

import (
	"errors"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type TaskService struct {
	TaskRepo repository.TaskRepositoryInterface
}

func NewTaskService(taskRepo repository.TaskRepositoryInterface) *TaskService {
	return &TaskService{TaskRepo: taskRepo}
}

func (s *TaskService) CreateTask(userID uint, title, description string, status model.TaskStatus, deadlineStr string) (*model.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if status == "" {
		status = model.StatusPending
	} else if !status.IsValid() {
		return nil, errors.New("invalid status value")
	}

	var deadline time.Time
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			return nil, errors.New("invalid deadline format, use RFC3339 e.g. 2026-08-15T17:00:00Z")
		}
		deadline = parsed
	}

	task := &model.Task{
		Title:       title,
		Description: description,
		Status:      status,
		UserID:      userID,
		Deadline:    deadline,
	}

	if err := s.TaskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTasks(userID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.TaskRepo.FindAll(userID, status, sort, page, limit)
}

func (s *TaskService) UpdateTask(taskID, userID uint, title, description string, status model.TaskStatus, deadlineStr string) (*model.Task, error) {
	task, err := s.TaskRepo.FindByID(taskID, userID)
	if err != nil {
		return nil, errors.New("task not found or access denied")
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if status != "" {
		if !status.IsValid() {
			return nil, errors.New("invalid status value")
		}
		task.Status = status
	}
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			return nil, errors.New("invalid deadline format, use RFC3339 e.g. 2026-08-15T17:00:00Z")
		}
		task.Deadline = parsed
	}

	if err := s.TaskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(taskID, userID uint) error {
	_, err := s.TaskRepo.FindByID(taskID, userID)
	if err != nil {
		return errors.New("task not found or access denied")
	}

	return s.TaskRepo.Delete(taskID, userID)
}