package repository

import "github.com/DanKRT-Star/task-manager/internal/model"

type UserRepositoryInterface interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

type TaskRepositoryInterface interface {
	Create(task *model.Task) error
	FindByID(taskID, userID uint) (*model.Task, error)
	FindByIDOnly(taskID uint) (*model.Task, error)
	FindAll(userID uint, status string, sort string, page, limit int) ([]model.Task, int64, error)
	FindAllByProject(projectID uint, status string, sort string, page, limit int) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(taskID uint) error
}

type ProjectRepositoryInterface interface {
	Create(project *model.Project) error
	FindByID(projectID uint) (*model.Project, error)
	FindAllByUser(userID uint) ([]model.Project, error)
	Update(project *model.Project) error
	Delete(projectID uint) error
}

type ProjectMemberRepositoryInterface interface {
	AddMember(member *model.ProjectMember) error
	FindMember(projectID, userID uint) (*model.ProjectMember, error)
	FindMembersByProject(projectID uint) ([]model.ProjectMember, error)
	RemoveMember(projectID, userID uint) error
}

type EpicRepositoryInterface interface {
	Create(epic *model.Epic) error
	FindByID(epicID uint) (*model.Epic, error)
	FindAllByProject(projectID uint) ([]model.Epic, error)
	Update(epic *model.Epic) error
	Delete(epicID uint) error
}

type MilestoneRepositoryInterface interface {
	Create(milestone *model.Milestone) error
	FindByID(milestoneID uint) (*model.Milestone, error)
	FindAllByProject(projectID uint) ([]model.Milestone, error)
	Update(milestone *model.Milestone) error
	Delete(milestoneID uint) error
}