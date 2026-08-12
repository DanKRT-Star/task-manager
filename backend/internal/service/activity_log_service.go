package service

import (
	"errors"

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
		return nil, errors.New("task not found")
	}

	// Cùng rule truy cập với Comment: ai xem được task thì xem được lịch sử của nó
	if task.ProjectID == nil {
		if task.UserID != userID {
			return nil, errors.New("you do not have access to this task")
		}
	} else {
		if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
			return nil, errors.New("you do not have access to this task")
		}
	}

	return s.ActivityRepo.FindAllByTask(taskID)
}