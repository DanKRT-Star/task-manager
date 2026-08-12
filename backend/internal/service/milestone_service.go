package service

import (
	"errors"
	"time"

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
		return nil, errors.New("you are not a member of this project")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}

	milestone := &model.Milestone{ProjectID: projectID, Title: title, Description: description, DueDate: dueDate}
	if err := s.MilestoneRepo.Create(milestone); err != nil {
		return nil, err
	}
	return milestone, nil
}

func (s *MilestoneService) GetProjectMilestones(userID, projectID uint) ([]model.Milestone, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("you are not a member of this project")
	}
	return s.MilestoneRepo.FindAllByProject(projectID)
}

func (s *MilestoneService) UpdateMilestone(userID, milestoneID uint, title, description string, dueDate *time.Time) (*model.Milestone, error) {
	milestone, err := s.MilestoneRepo.FindByID(milestoneID)
	if err != nil {
		return nil, errors.New("milestone not found")
	}
	if _, err := s.MemberRepo.FindMember(milestone.ProjectID, userID); err != nil {
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
		return nil, err
	}
	return milestone, nil
}

func (s *MilestoneService) DeleteMilestone(userID, milestoneID uint) error {
	milestone, err := s.MilestoneRepo.FindByID(milestoneID)
	if err != nil {
		return errors.New("milestone not found")
	}
	member, err := s.MemberRepo.FindMember(milestone.ProjectID, userID)
	if err != nil {
		return errors.New("you are not a member of this project")
	}
	if member.Role != model.RoleOwner {
		return errors.New("only the project owner can delete a milestone")
	}
	return s.MilestoneRepo.Delete(milestoneID)
}