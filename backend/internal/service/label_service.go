package service

import (
	"errors"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type LabelService struct {
	LabelRepo  repository.LabelRepositoryInterface
	TaskRepo   repository.TaskRepositoryInterface
	MemberRepo repository.ProjectMemberRepositoryInterface
}

func NewLabelService(labelRepo repository.LabelRepositoryInterface, taskRepo repository.TaskRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *LabelService {
	return &LabelService{LabelRepo: labelRepo, TaskRepo: taskRepo, MemberRepo: memberRepo}
}

func (s *LabelService) CreateLabel(userID, projectID uint, name, color string) (*model.Label, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.LabelProjectAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	if name == "" {
		logger.LabelCreateMissingName(userID, projectID)
		return nil, errors.New("name is required")
	}
	if color == "" {
		color = "#6b7280"
	}

	label := &model.Label{ProjectID: projectID, Name: name, Color: color}
	if err := s.LabelRepo.Create(label); err != nil {
		logger.LabelCreateFailed(projectID, userID, name, err)
		return nil, err
	}
	logger.LabelCreated(label.LabelID, projectID, userID)
	return label, nil
}

func (s *LabelService) GetProjectLabels(userID, projectID uint) ([]model.Label, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.LabelProjectAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	labels, err := s.LabelRepo.FindAllByProject(projectID)
	if err != nil {
		logger.LabelFetchFailed(projectID, err)
		return nil, err
	}
	return labels, nil
}

func (s *LabelService) DeleteLabel(userID, labelID uint) error {
	label, err := s.LabelRepo.FindByID(labelID)
	if err != nil {
		logger.LabelNotFound(labelID, userID)
		return errors.New("label not found")
	}
	member, err := s.MemberRepo.FindMember(label.ProjectID, userID)
	if err != nil {
		logger.LabelProjectAccessDenied(userID, label.ProjectID)
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		logger.LabelOwnerRequired(userID, label.ProjectID, labelID)
		return errors.New("only the project owner can delete a label")
	}
	if err := s.LabelRepo.Delete(labelID); err != nil {
		logger.LabelDeleteFailed(labelID, userID, err)
		return err
	}
	logger.LabelDeleted(labelID, userID)
	return nil
}

func (s *LabelService) AttachLabel(userID, taskID, labelID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.LabelTaskNotFound(taskID, userID)
		return errors.New("task not found")
	}
	if task.ProjectID == nil {
		logger.LabelTaskNotInProject(taskID, userID)
		return errors.New("labels can only be attached to tasks within a project")
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		logger.LabelProjectAccessDenied(userID, *task.ProjectID)
		return errors.New("you are not a member of this project")
	}

	label, err := s.LabelRepo.FindByID(labelID)
	if err != nil {
		logger.LabelNotFound(labelID, userID)
		return errors.New("label not found")
	}
	if label.ProjectID != *task.ProjectID {
		logger.LabelProjectMismatch(taskID, labelID, *task.ProjectID)
		return errors.New("label does not belong to the same project as the task")
	}

	if err := s.LabelRepo.AttachToTask(taskID, labelID); err != nil {
		logger.LabelAttachFailed(labelID, taskID, userID, err)
		return err
	}
	logger.LabelAttached(taskID, labelID, userID)
	return nil
}

func (s *LabelService) DetachLabel(userID, taskID, labelID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.LabelTaskNotFound(taskID, userID)
		return errors.New("task not found")
	}
	if task.ProjectID == nil {
		logger.LabelTaskNoLabels(taskID, userID)
		return errors.New("this task has no labels")
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		logger.LabelProjectAccessDenied(userID, *task.ProjectID)
		return errors.New("you are not a member of this project")
	}

	if err := s.LabelRepo.DetachFromTask(taskID, labelID); err != nil {
		logger.LabelDetachFailed(labelID, taskID, userID, err)
		return err
	}
	logger.LabelDetached(taskID, labelID, userID)
	return nil
}

func (s *LabelService) GetTaskLabels(userID, taskID uint) ([]model.Label, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		logger.LabelTaskNotFound(taskID, userID)
		return nil, errors.New("task not found")
	}
	if task.ProjectID == nil {
		return []model.Label{}, nil
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		logger.LabelProjectAccessDenied(userID, *task.ProjectID)
		return nil, errors.New("you are not a member of this project")
	}
	labels, err := s.LabelRepo.FindLabelsByTask(taskID)
	if err != nil {
		logger.TaskLabelsFetchFailed(taskID, err)
		return nil, err
	}
	return labels, nil
}