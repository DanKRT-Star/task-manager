package service

import (
	"errors"

	"github.com/DanKRT-Star/task-manager/internal/logger"
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
		logger.EpicCreateAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	if title == "" {
		logger.EpicCreateMissingTitle(userID, projectID)
		return nil, errors.New("title is required")
	}

	epic := &model.Epic{ProjectID: projectID, Title: title, Description: description}
	if err := s.EpicRepo.Create(epic); err != nil {
		logger.EpicCreateFailed(projectID, userID, title, err)
		return nil, err
	}
	logger.EpicCreated(epic.EpicID, projectID, userID)
	return epic, nil
}

func (s *EpicService) GetProjectEpics(userID, projectID uint) ([]model.Epic, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.EpicFetchAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	epics, err := s.EpicRepo.FindAllByProject(projectID)
	if err != nil {
		logger.EpicFetchFailed(projectID, err)
		return nil, err
	}
	return epics, nil
}

func (s *EpicService) UpdateEpic(userID, epicID uint, title, description string) (*model.Epic, error) {
	epic, err := s.EpicRepo.FindByID(epicID)
	if err != nil {
		logger.EpicNotFound(epicID, userID)
		return nil, errors.New("epic not found")
	}
	if _, err := s.MemberRepo.FindMember(epic.ProjectID, userID); err != nil {
		logger.EpicUpdateAccessDenied(epicID, userID, epic.ProjectID)
		return nil, errors.New("you are not a member of this project")
	}

	if title != "" {
		epic.Title = title
	}
	if description != "" {
		epic.Description = description
	}

	if err := s.EpicRepo.Update(epic); err != nil {
		logger.EpicUpdateFailed(epicID, userID, err)
		return nil, err
	}
	logger.EpicUpdated(epicID, userID)
	return epic, nil
}

func (s *EpicService) DeleteEpic(userID, epicID uint) error {
	epic, err := s.EpicRepo.FindByID(epicID)
	if err != nil {
		logger.EpicNotFound(epicID, userID)
		return errors.New("epic not found")
	}
	member, err := s.MemberRepo.FindMember(epic.ProjectID, userID)
	if err != nil {
		logger.EpicDeleteAccessDenied(epicID, userID, epic.ProjectID)
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		logger.EpicOwnerRequired(epicID, epic.ProjectID, userID)
		return errors.New("only the project owner can delete an epic")
	}
	if err := s.EpicRepo.Delete(epicID); err != nil {
		logger.EpicDeleteFailed(epicID, userID, err)
		return err
	}
	logger.EpicDeleted(epicID, userID)
	return nil
}