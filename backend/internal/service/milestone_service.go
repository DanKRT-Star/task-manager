package service

import (
	"errors"
	"time"

	"github.com/DanKRT-Star/task-manager/internal/logger"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type MilestoneService struct {
	MilestoneRepo repository.MilestoneRepositoryInterface
	MemberRepo    repository.ProjectMemberRepositoryInterface
}

func NewMilestoneService(milestoneRepo repository.MilestoneRepositoryInterface, memberRepo repository.ProjectMemberRepositoryInterface) *MilestoneService {
	return &MilestoneService{MilestoneRepo: milestoneRepo, MemberRepo: memberRepo}
}

func (s *MilestoneService) CreateMilestone(userID, projectID uint, title, description string, dueDate *time.Time) (*model.Milestone, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.MilestoneCreateAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	if title == "" {
		logger.MilestoneCreateMissingTitle(userID, projectID)
		return nil, errors.New("title is required")
	}

	milestone := &model.Milestone{ProjectID: projectID, Title: title, Description: description, DueDate: dueDate}
	if err := s.MilestoneRepo.Create(milestone); err != nil {
		logger.MilestoneCreateFailed(projectID, userID, title, err)
		return nil, err
	}
	logger.MilestoneCreated(milestone.MilestoneID, projectID, userID)
	return milestone, nil
}

func (s *MilestoneService) GetProjectMilestones(userID, projectID uint) ([]model.Milestone, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		logger.MilestoneFetchAccessDenied(userID, projectID)
		return nil, errors.New("you are not a member of this project")
	}
	milestones, err := s.MilestoneRepo.FindAllByProject(projectID)
	if err != nil {
		logger.MilestoneFetchFailed(projectID, err)
		return nil, err
	}
	return milestones, nil
}

func (s *MilestoneService) UpdateMilestone(userID, milestoneID uint, title, description string, dueDate *time.Time) (*model.Milestone, error) {
	milestone, err := s.MilestoneRepo.FindByID(milestoneID)
	if err != nil {
		logger.MilestoneNotFound(milestoneID, userID)
		return nil, errors.New("milestone not found")
	}
	if _, err := s.MemberRepo.FindMember(milestone.ProjectID, userID); err != nil {
		logger.MilestoneUpdateAccessDenied(milestoneID, userID, milestone.ProjectID)
		return nil, errors.New("you are not a member of this project")
	}

	if title != "" {
		milestone.Title = title
	}
	if description != "" {
		milestone.Description = description
	}
	if dueDate != nil {
		milestone.DueDate = dueDate
	}

	if err := s.MilestoneRepo.Update(milestone); err != nil {
		logger.MilestoneUpdateFailed(milestoneID, userID, err)
		return nil, err
	}
	logger.MilestoneUpdated(milestoneID, userID)
	return milestone, nil
}

func (s *MilestoneService) DeleteMilestone(userID, milestoneID uint) error {
	milestone, err := s.MilestoneRepo.FindByID(milestoneID)
	if err != nil {
		logger.MilestoneNotFound(milestoneID, userID)
		return errors.New("milestone not found")
	}
	member, err := s.MemberRepo.FindMember(milestone.ProjectID, userID)
	if err != nil {
		logger.MilestoneDeleteAccessDenied(milestoneID, userID, milestone.ProjectID)
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		logger.MilestoneOwnerRequired(milestoneID, milestone.ProjectID, userID)
		return errors.New("only the project owner can delete a milestone")
	}
	if err := s.MilestoneRepo.Delete(milestoneID); err != nil {
		logger.MilestoneDeleteFailed(milestoneID, userID, err)
		return err
	}
	logger.MilestoneDeleted(milestoneID, userID)
	return nil
}