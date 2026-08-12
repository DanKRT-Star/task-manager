package service

import (
	"errors"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type EpicService struct {
	EpicRepo   repository.EpicRepositoryInterface
	MemberRepo repository.ProjectMemberRepositoryInterface
}

func NewEpicService(epicRepo repository.EpicRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *EpicService {
	return &EpicService{EpicRepo: epicRepo, MemberRepo: memberRepo}
}

func (s *EpicService) CreateEpic(userID, projectID uint, title, description string) (*model.Epic, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}

	epic := &model.Epic{ProjectID: projectID, Title: title, Description: description}
	if err := s.EpicRepo.Create(epic); err != nil {
		return nil, err
	}
	return epic, nil
}

func (s *EpicService) GetProjectEpics(userID, projectID uint) ([]model.Epic, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	return s.EpicRepo.FindAllByProject(projectID)
}

func (s *EpicService) UpdateEpic(userID, epicID uint, title, description string) (*model.Epic, error) {
	epic, err := s.EpicRepo.FindByID(epicID)
	if err != nil {
		return nil, errors.New("epic not found")
	}
	if _, err := s.MemberRepo.FindMember(epic.ProjectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}

	if title != "" {
		epic.Title = title
	}
	if description != "" {
		epic.Description = description
	}

	if err := s.EpicRepo.Update(epic); err != nil {
		return nil, err
	}
	return epic, nil
}

func (s *EpicService) DeleteEpic(userID, epicID uint) error {
	epic, err := s.EpicRepo.FindByID(epicID)
	if err != nil {
		return errors.New("epic not found")
	}
	member, err := s.MemberRepo.FindMember(epic.ProjectID, userID)
	if err != nil {
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		return errors.New("only the project owner can delete an epic")
	}
	return s.EpicRepo.Delete(epicID)
}