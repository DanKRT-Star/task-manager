package service

import (
	"errors"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type SprintService struct {
	SprintRepo repository.SprintRepositoryInterface
	MemberRepo repository.ProjectMemberRepositoryInterface
}

func NewSprintService(sprintRepo repository.SprintRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *SprintService {
	return &SprintService{SprintRepo: sprintRepo, MemberRepo: memberRepo}
}

func (s *SprintService) CreateSprint(userID, projectID uint, name string, startDate, endDate *time.Time) (*model.Sprint, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.SprintProjectAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	if name == "" {
		logger.SprintCreateMissingName(userID, projectID)
		return nil, errors.New("name is required")
	}

	sprint := &model.Sprint{
		ProjectID: projectID,
		Name:      name,
		Status:    model.SprintPlanned,
		StartDate: startDate,
		EndDate:   endDate,
	}
	if err := s.SprintRepo.Create(sprint); err != nil {
		logger.SprintCreateFailed(projectID, userID, name, err)
		return nil, err
	}
	logger.SprintCreated(sprint.SprintID, projectID, userID)
	return sprint, nil
}

func (s *SprintService) GetProjectSprints(userID, projectID uint) ([]model.Sprint, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.SprintProjectAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	sprints, err := s.SprintRepo.FindAllByProject(projectID)
	if err != nil {
		logger.SprintFetchFailed(projectID, err)
		return nil, err
	}
	return sprints, nil
}

func (s *SprintService) UpdateSprint(userID, sprintID uint, name string, status model.SprintStatus, startDate, endDate *time.Time) (*model.Sprint, error) {
	sprint, err := s.SprintRepo.FindByID(sprintID)
	if err != nil {
		logger.SprintNotFound(sprintID, userID)
		return nil, errors.New("sprint not found")
	}
	if _, err := s.MemberRepo.FindMember(sprint.ProjectID, userID); err != nil {
		logger.SprintProjectAccessDenied(userID, sprint.ProjectID)
		return nil, errors.New("you are not a member of this project")
	}

	if name != "" {
		sprint.Name = name
	}
	if status != "" {
		if !status.IsValid() {
			logger.SprintInvalidStatus(sprintID, string(status))
			return nil, errors.New("invalid status value")
		}
		sprint.Status = status
	}
	if startDate != nil {
		sprint.StartDate = startDate
	}
	if endDate != nil {
		sprint.EndDate = endDate
	}

	if err := s.SprintRepo.Update(sprint); err != nil {
		logger.SprintUpdateFailed(sprintID, userID, err)
		return nil, err
	}
	logger.SprintUpdated(sprintID, userID)
	return sprint, nil
}

func (s *SprintService) DeleteSprint(userID, sprintID uint) error {
	sprint, err := s.SprintRepo.FindByID(sprintID)
	if err != nil {
		logger.SprintNotFound(sprintID, userID)
		return errors.New("sprint not found")
	}
	member, err := s.MemberRepo.FindMember(sprint.ProjectID, userID)
	if err != nil {
		logger.SprintProjectAccessDenied(userID, sprint.ProjectID)
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		logger.SprintOwnerRequired(sprintID, sprint.ProjectID, userID)
		return errors.New("only the project owner can delete a sprint")
	}
	if err := s.SprintRepo.Delete(sprintID); err != nil {
		logger.SprintDeleteFailed(sprintID, userID, err)
		return err
	}
	logger.SprintDeleted(sprintID, userID)
	return nil
}