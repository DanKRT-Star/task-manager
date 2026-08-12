package service

import (
	"errors"

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
		return nil, errors.New("you are not a member of this project")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if color == "" {
		color = "#6b7280"
	}

	label := &model.Label{ProjectID: projectID, Name: name, Color: color}
	if err := s.LabelRepo.Create(label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) GetProjectLabels(userID, projectID uint) ([]model.Label, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	return s.LabelRepo.FindAllByProject(projectID)
}

func (s *LabelService) DeleteLabel(userID, labelID uint) error {
	label, err := s.LabelRepo.FindByID(labelID)
	if err != nil {
		return errors.New("label not found")
	}
	member, err := s.MemberRepo.FindMember(label.ProjectID, userID)
	if err != nil {
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		return errors.New("only the project owner can delete a label")
	}
	return s.LabelRepo.Delete(labelID)
}

func (s *LabelService) AttachLabel(userID, taskID, labelID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.ProjectID == nil {
		return errors.New("labels can only be attached to tasks within a project")
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		return errors.New("you are not a member of this project")
	}

	label, err := s.LabelRepo.FindByID(labelID)
	if err != nil {
		return errors.New("label not found")
	}
	if label.ProjectID != *task.ProjectID {
		return errors.New("label does not belong to the same project as the task")
	}

	return s.LabelRepo.AttachToTask(taskID, labelID)
}

func (s *LabelService) DetachLabel(userID, taskID, labelID uint) error {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return errors.New("task not found")
	}
	if task.ProjectID == nil {
		return errors.New("this task has no labels")
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		return errors.New("you are not a member of this project")
	}

	return s.LabelRepo.DetachFromTask(taskID, labelID)
}

func (s *LabelService) GetTaskLabels(userID, taskID uint) ([]model.Label, error) {
	task, err := s.TaskRepo.FindByIDOnly(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}
	if task.ProjectID == nil {
		return []model.Label{}, nil
	}
	if _, err := s.MemberRepo.FindMember(*task.ProjectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	return s.LabelRepo.FindLabelsByTask(taskID)
}