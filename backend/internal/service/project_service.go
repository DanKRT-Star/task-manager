package service

import (
	"errors"

	"time"
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/DanKRT-Star/task-manager/internal/repository"
)

type ProjectService struct {
	ProjectRepo repository.ProjectRepositoryInterface
	MemberRepo  repository.ProjectMemberRepositoryInterface
	UserRepo    repository.UserRepositoryInterface
}

func NewProjectService(
	projectRepo repository.ProjectRepositoryInterface,
	memberRepo repository.ProjectMemberRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *ProjectService {
	return &ProjectService{ProjectRepo: projectRepo, MemberRepo: memberRepo, UserRepo: userRepo}
}

// CreateProject creates a project and adds the creator as owner.
// Note: not fully atomic across two tables; acceptable for current scope.
func (s *ProjectService) CreateProject(userID uint, name, description string, deadline *time.Time) (*model.Project, error) {
	if name == "" {
		return nil, errors.New("project name is required")
	}

	project := &model.Project{
		Name:        name,
		Description: description,
		Deadline:    deadline,
		OwnerID:     userID,
	}

	if err := s.ProjectRepo.Create(project); err != nil {
		return nil, err
	}

	member := &model.ProjectMember{
		ProjectID: project.ProjectID,
		UserID:    userID,
		Role:      model.RoleOwner,
	}
	if err := s.MemberRepo.AddMember(member); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetProject(userID, projectID uint) (*model.Project, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("project not found or access denied")
	}

	project, err := s.ProjectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return project, nil
}

func (s *ProjectService) GetUserProjects(userID uint) ([]model.Project, error) {
	return s.ProjectRepo.FindAllByUser(userID)
}

func (s *ProjectService) UpdateProject(userID, projectID uint, name, description string, deadline *time.Time) (*model.Project, error) {
	member, err := s.MemberRepo.FindMember(projectID, userID)
	if err != nil || member.Role != model.RoleOwner {
		return nil, errors.New("only the project owner can update the project")
	}

	project, err := s.ProjectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if deadline != nil {
		project.Deadline = deadline
	}

	if err := s.ProjectRepo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) DeleteProject(userID, projectID uint) error {
	member, err := s.MemberRepo.FindMember(projectID, userID)
	if err != nil || member.Role != model.RoleOwner {
		return errors.New("only the project owner can delete the project")
	}
	return s.ProjectRepo.Delete(projectID)
}

func (s *ProjectService) AddMember(userID, projectID uint, email string) (*model.ProjectMember, error) {
	requester, err := s.MemberRepo.FindMember(projectID, userID)
	if err != nil || requester.Role != model.RoleOwner {
		return nil, errors.New("only the project owner can add members")
	}

	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user with this email does not exist")
	}

	existing, _ := s.MemberRepo.FindMember(projectID, user.UserID)
	if existing != nil {
		return nil, errors.New("user is already a member of this project")
	}

	member := &model.ProjectMember{
		ProjectID: projectID,
		UserID:    user.UserID,
		Role:      model.RoleMember,
	}
	if err := s.MemberRepo.AddMember(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *ProjectService) RemoveMember(userID, projectID, targetUserID uint) error {
	requester, err := s.MemberRepo.FindMember(projectID, userID)
	if err != nil || requester.Role != model.RoleOwner {
		return errors.New("only the project owner can remove members")
	}

	target, err := s.MemberRepo.FindMember(projectID, targetUserID)
	if err != nil {
		return errors.New("member not found")
	}
	if target.Role == model.RoleOwner {
		return errors.New("cannot remove the project owner")
	}

	return s.MemberRepo.RemoveMember(projectID, targetUserID)
}

func (s *ProjectService) GetMembers(userID, projectID uint) ([]model.ProjectMember, error) {
	if _, err := s.MemberRepo.FindMember(projectID, userID); err != nil {
		return nil, errors.New("project not found or access denied")
	}
	return s.MemberRepo.FindMembersByProject(projectID)
}