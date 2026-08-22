package service

import (
	"errors"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type TaskService struct {
	TaskRepo      repository.TaskRepositoryInterface
	MemberRepo    repository.ProjectMemberRepositoryInterface
	ActivityRepo  repository.ActivityLogRepositoryInterface
	EpicRepo      repository.EpicRepositoryInterface
	MilestoneRepo repository.MilestoneRepositoryInterface
	SprintRepo    repository.SprintRepositoryInterface
}

func NewTaskService(
	taskRepo repository.TaskRepositoryInterface,
	memberRepo repository.ProjectMemberRepositoryInterface,
	activityRepo repository.ActivityLogRepositoryInterface,
	epicRepo repository.EpicRepositoryInterface,
	milestoneRepo repository.MilestoneRepositoryInterface,
	sprintRepo repository.SprintRepositoryInterface,
) *TaskService {
	return &TaskService{
		TaskRepo:      taskRepo,
		MemberRepo:    memberRepo,
		ActivityRepo:  activityRepo,
		EpicRepo:      epicRepo,
		MilestoneRepo: milestoneRepo,
		SprintRepo:    sprintRepo,
	}
}

// validateEpicInProject đảm bảo epic được chọn thuộc đúng project của task,
// tránh trường hợp gắn nhầm epic của project khác.
func (s *TaskService) validateEpicInProject(epicID, projectID uint) error {
	epic, err := s.EpicRepo.FindByID(epicID)
	if err != nil {
		return errors.New("epic not found")
	}
	if epic.ProjectID != projectID {
		return errors.New("epic does not belong to this project")
	}
	return nil
}

// validateMilestoneInProject đảm bảo milestone được chọn thuộc đúng project của task.
func (s *TaskService) validateMilestoneInProject(milestoneID, projectID uint) error {
	milestone, err := s.MilestoneRepo.FindByID(milestoneID)
	if err != nil {
		return errors.New("milestone not found")
	}
	if milestone.ProjectID != projectID {
		return errors.New("milestone does not belong to this project")
	}
	return nil
}

// validateSprintInProject đảm bảo sprint được chọn thuộc đúng project của task.
func (s *TaskService) validateSprintInProject(sprintID, projectID uint) error {
	sprint, err := s.SprintRepo.FindByID(sprintID)
	if err != nil {
		return errors.New("sprint not found")
	}
	if sprint.ProjectID != projectID {
		return errors.New("sprint does not belong to this project")
	}
	return nil
}

// Thứ tự tham số phải khớp với cách task_handler.go gọi hàm này:
// h.TaskService.CreateTask(userID, req.Title, req.Description, req.Status, req.Deadline, req.ProjectID, req.AssigneeID, req.EpicID, req.MilestoneID, req.SprintID)
func (s *TaskService) CreateTask(
	userID uint,
	title, description string,
	status model.TaskStatus,
	deadlineStr string,
	projectID, assigneeID, epicID, milestoneID, sprintID *uint,
) (*model.Task, error) {
	if title == "" {
		logger.TaskMissingTitle(userID)
		return nil, errors.New("title is required")
	}

	if status == "" {
		status = model.StatusPending
	} else if !status.IsValid() {
		logger.TaskCreateInvalidStatus(userID, string(status))
		return nil, errors.New("invalid status value")
	}

	if projectID != nil {
		if _, err := s.MemberRepo.FindMember(*projectID, userID); err != nil {
			logger.TaskCreateAccessDenied(userID, *projectID)
			return nil, errors.New("you are not a member of this project")
		}
		if assigneeID != nil {
			if _, err := s.MemberRepo.FindMember(*projectID, *assigneeID); err != nil {
				logger.TaskCreateAssigneeAccessDenied(userID, *projectID, *assigneeID)
				return nil, errors.New("assignee is not a member of this project")
			}
		}
		if epicID != nil {
			if err := s.validateEpicInProject(*epicID, *projectID); err != nil {
				return nil, err
			}
		}
		if milestoneID != nil {
			if err := s.validateMilestoneInProject(*milestoneID, *projectID); err != nil {
				return nil, err
			}
		}
		if sprintID != nil {
			if err := s.validateSprintInProject(*sprintID, *projectID); err != nil {
				return nil, err
			}
		}
	} else {
		if assigneeID != nil {
			logger.TaskCreateAssigneeWithoutProject(userID, *assigneeID)
			return nil, errors.New("cannot assign a task that does not belong to a project")
		}
		if epicID != nil || milestoneID != nil || sprintID != nil {
			return nil, errors.New("cannot assign epic, milestone, or sprint to a task that does not belong to a project")
		}
	}

	var deadline time.Time
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			logger.TaskCreateInvalidDeadline(userID, deadlineStr)
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
		EpicID:      epicID,
		MilestoneID: milestoneID,
		SprintID:    sprintID,
		Deadline:    deadline,
	}

	if err := s.TaskRepo.Create(task); err != nil {
		logger.TaskCreateFailed(userID, title, err)
		return nil, err
	}

	if err := s.ActivityRepo.Create(&model.ActivityLog{
		TaskID: task.TaskID,
		UserID: userID,
		Action: model.ActionCreated,
		Detail: "Task created",
	}); err != nil {
		logger.TaskActivityLogCreateFailed(task.TaskID, userID, err)
	}

	logger.TaskCreated(task.TaskID, userID, title)
	return task, nil
}

func (s *TaskService) GetTasks(userID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	page, limit = normalizePagination(page, limit)
	tasks, total, err := s.TaskRepo.FindAll(userID, status, sort, page, limit)
	if err != nil {
		logger.TaskListFetchFailed(userID, status, err)
		return nil, 0, err
	}
	return tasks, total, nil
}

func (s *TaskService) GetProjectTasks(userID, projectID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.ProjectTaskListAccessDenied(userID, projectID)
		return nil, 0, errors.New("you are not a member of this project")
	}
	page, limit = normalizePagination(page, limit)
	tasks, total, err := s.TaskRepo.FindAllByProject(projectID, status, sort, page, limit)
	if err != nil {
		logger.ProjectTaskListFetchFailed(projectID, err)
		return nil, 0, err
	}
	return tasks, total, nil
}

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
	return task.UserID == userID || (task.AssigneeID != nil && *task.AssigneeID == userID)
}

// Thứ tự tham số phải khớp với cách task_handler.go gọi hàm này:
// h.TaskService.UpdateTask(uint(taskID), userID, req.Title, req.Description, req.Status, req.Deadline, req.AssigneeID, req.EpicID, req.MilestoneID, req.SprintID)
func (s *TaskService) UpdateTask(
	taskID, userID uint,
	title, description string,
	status model.TaskStatus,
	deadlineStr string,
	assigneeID, epicID, milestoneID, sprintID *uint,
) (*model.Task, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.TaskUpdateNotFound(taskID, userID)
		return nil, errors.New("task not found")
	}

	if !s.canModifyTask(task, userID) {
		logger.TaskUpdateAccessDenied(taskID, userID)
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
			logger.TaskUpdateInvalidStatus(taskID, string(status))
			return nil, errors.New("invalid status value")
		}
		task.Status = status
	}
	if deadlineStr != "" {
		parsed, err := time.Parse(time.RFC3339, deadlineStr)
		if err != nil {
			logger.TaskUpdateInvalidDeadline(taskID, deadlineStr)
			return nil, errors.New("invalid deadline format, use RFC3339 e.g. 2026-08-15T17:00:00Z")
		}
		task.Deadline = parsed
	}
	if assigneeID != nil {
		if task.ProjectID == nil {
			logger.TaskUpdateAssigneeWithoutProject(taskID, *assigneeID)
			return nil, errors.New("cannot assign a task that does not belong to a project")
		}
		if _, err := s.MemberRepo.FindMember(*task.ProjectID, *assigneeID); err != nil {
			logger.TaskUpdateAssigneeNotProjectMember(taskID, *task.ProjectID, *assigneeID)
			return nil, errors.New("assignee is not a member of this project")
		}
		task.AssigneeID = assigneeID
	}
	if epicID != nil {
		if task.ProjectID == nil {
			return nil, errors.New("cannot assign epic to a task that does not belong to a project")
		}
		if err := s.validateEpicInProject(*epicID, *task.ProjectID); err != nil {
			return nil, err
		}
		task.EpicID = epicID
	}
	if milestoneID != nil {
		if task.ProjectID == nil {
			return nil, errors.New("cannot assign milestone to a task that does not belong to a project")
		}
		if err := s.validateMilestoneInProject(*milestoneID, *task.ProjectID); err != nil {
			return nil, err
		}
		task.MilestoneID = milestoneID
	}
	if sprintID != nil {
		if task.ProjectID == nil {
			return nil, errors.New("cannot assign sprint to a task that does not belong to a project")
		}
		if err := s.validateSprintInProject(*sprintID, *task.ProjectID); err != nil {
			return nil, err
		}
		task.SprintID = sprintID
	}

	if err := s.TaskRepo.Update(task); err != nil {
		logger.TaskUpdateFailed(taskID, userID, err)
		return nil, err
	}

	if status != "" && status != oldStatus {
		if err := s.ActivityRepo.Create(&model.ActivityLog{
			TaskID: task.TaskID,
			UserID: userID,
			Action: model.ActionStatusChanged,
			Detail: "Status changed from " + string(oldStatus) + " to " + string(status),
		}); err != nil {
			logger.TaskStatusActivityLogFailed(task.TaskID, userID, err)
		}
	}
	logger.TaskUpdated(task.TaskID, userID)
	return task, nil
}

func (s *TaskService) DeleteTask(taskID, userID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.TaskDeleteNotFound(taskID, userID)
		return errors.New("task not found")
	}

	if !s.canModifyTask(task, userID) {
		logger.TaskDeleteAccessDenied(taskID, userID)
		return errors.New("you do not have permission to delete this task")
	}

	if err := s.TaskRepo.Delete(taskID); err != nil {
		logger.TaskDeleteFailed(taskID, userID, err)
		return err
	}
	logger.TaskDeleted(taskID, userID)
	return nil
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