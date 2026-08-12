package service

import (
	"errors"
	"time"

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
		return nil, errors.New("you are not a member of this project")
	}
	if name == "" {
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
		return nil, err
	}
	return sprint, nil
}

func (s *SprintService) GetProjectSprints(userID, projectID uint) ([]model.Sprint, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	return s.SprintRepo.FindAllByProject(projectID)
}

func (s *SprintService) UpdateSprint(userID, sprintID uint, name string, status model.SprintStatus, startDate, endDate *time.Time) (*model.Sprint, error) {
	sprint, err := s.SprintRepo.FindByID(sprintID)
	if err != nil {
		return nil, errors.New("sprint not found")
	}
	if _, err := s.MemberRepo.FindMember(sprint.ProjectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}

	if name != "" {
		sprint.Name = name
	}
	if status != "" {
		if !status.IsValid() {
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
		return nil, err
	}
	return sprint, nil
}

func (s *SprintService) DeleteSprint(userID, sprintID uint) error {
	sprint, err := s.SprintRepo.FindByID(sprintID)
	if err != nil {
		return errors.New("sprint not found")
	}
	member, err := s.MemberRepo.FindMember(sprint.ProjectID, userID)
	if err != nil {
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		return errors.New("only the project owner can delete a sprint")
	}
	return s.SprintRepo.Delete(sprintID)
}