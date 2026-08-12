package service

import (
	"errors"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type TaskService struct {
	TaskRepo     repository.TaskRepositoryInterface
	MemberRepo   repository.ProjectMemberRepositoryInterface
	ActivityRepo repository.ActivityLogRepositoryInterface
}

func NewTaskService(taskRepo repository.TaskRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface, activityRepo repository.ActivityLogRepositoryInterface) *TaskService {
	return &TaskService{TaskRepo: taskRepo, MemberRepo: memberRepo, ActivityRepo: activityRepo}
}

func (s *TaskService) CreateTask(userID uint, title, description string, status model.TaskStatus, deadlineStr string, projectID, assigneeID *uint) (*model.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if status == "" {
		status = model.StatusPending
	} else if !status.IsValid() {
		return nil, errors.New("invalid status value")
	}

	// If attaching to a project, requester must be a member (any role)
	if projectID != nil {
		if _, err := s.MemberRepo.FindMember(*projectID, userID); err != nil {
			return nil, errors.New("you are not a member of this project")
		}
		// If assigning to someone, that person must also be a project member
		if assigneeID != nil {
			if _, err := s.MemberRepo.FindMember(*projectID, *assigneeID); err != nil {
				return nil, errors.New("assignee is not a member of this project")
			}
		}
	} else if assigneeID != nil {
		return nil, errors.New("cannot assign a task that does not belong to a project")
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
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
		Deadline:    deadline,
	}

	if err := s.TaskRepo.Create(task); err != nil {
		return nil, err
	}

	s.ActivityRepo.Create(&model.ActivityLog{
		TaskID: task.TaskID,
		UserID: userID,
		Action: model.ActionCreated,
		Detail: "Task created",
	})

	return task, nil
}

func (s *TaskService) GetTasks(userID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	page, limit = normalizePagination(page, limit)
	return s.TaskRepo.FindAll(userID, status, sort, page, limit)
}

func (s *TaskService) GetProjectTasks(userID, projectID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, 0, errors.New("you are not a member of this project")
	}
	page, limit = normalizePagination(page, limit)
	return s.TaskRepo.FindAllByProject(projectID, status, sort, page, limit)
}

// canModifyTask centralizes the authorization rule for update/delete.
func (s *TaskService) canModifyTask(task *model.Task, userID uint) bool {
	if task.ProjectID == nil {
		return task.UserID == userID
	}

	member, err := s.MemberRepo.FindMember(*task.ProjectID, userID)
	if err != nil {
		return false
	}
	if member.Role == model.RoleOwner {
		return true
	}
	// member: only their own created task, or a task assigned to them
	return task.UserID == userID || (task.AssigneeID != nil && *task.AssigneeID == userID)
}

func (s *TaskService) UpdateTask(taskID, userID uint, title, description string, status model.TaskStatus, deadlineStr string, assigneeID *uint) (*model.Task, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	if !s.canModifyTask(task, userID) {
		return nil, errors.New("you do not have permission to modify this task")
	}

	oldStatus := task.Status
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
	if assigneeID != nil {
		if task.ProjectID == nil {
			return nil, errors.New("cannot assign a task that does not belong to a project")
		}
		if _, err := s.MemberRepo.FindMember(*task.ProjectID, *assigneeID); err != nil {
			return nil, errors.New("assignee is not a member of this project")
		}
		task.AssigneeID = assigneeID
	}

	if err := s.TaskRepo.Update(task); err != nil {
		return nil, err
	}

	if status != "" && status != oldStatus {
		s.ActivityRepo.Create(&model.ActivityLog{
			TaskID: task.TaskID,
			UserID: userID,
			Action: model.ActionStatusChanged,
			Detail: "Status changed from " + string(oldStatus) + " to " + string(status),
		})
	}
	return task, nil
}

func (s *TaskService) DeleteTask(taskID, userID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if !s.canModifyTask(task, userID) {
		return errors.New("you do not have permission to delete this task")
	}

	return s.TaskRepo.Delete(taskID)
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}
