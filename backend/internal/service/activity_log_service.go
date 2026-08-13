package service

import (
	"errors"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type ActivityLogService struct {
	ActivityRepo repository.ActivityLogRepositoryInterface
	TaskRepo     repository.TaskRepositoryInterface
	MemberRepo   repository.ProjectMemberRepositoryInterface
}

func NewActivityLogService(activityRepo repository.ActivityLogRepositoryInterface, taskRepo repository.TaskRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *ActivityLogService {
	return &ActivityLogService{ActivityRepo: activityRepo, TaskRepo: taskRepo, MemberRepo: memberRepo}
}

func (s *ActivityLogService) GetTaskActivity(userID, taskID uint) ([]model.ActivityLog, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.ActivityFetchTaskNotFound(taskID, userID)
		return nil, errors.New("task not found")
	}

	if task.ProjectID == nil {
		if task.UserID != userID {
			logger.ActivityFetchAccessDenied(taskID, userID, 0)
			return nil, errors.New("you do not have access to this task")
		}
	} else {
		if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
			logger.ActivityFetchAccessDenied(taskID, userID, *task.ProjectID)
			return nil, errors.New("you do not have access to this task")
		}
	}

	logs, err := s.ActivityRepo.FindAllByTask(taskID)
	if err != nil {
		logger.ActivityFetchFailed(taskID, err)
		return nil, err
	}
	return logs, nil
}
