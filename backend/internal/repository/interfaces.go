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

type SprintRepositoryInterface interface {
	Create(sprint *model.Sprint) error
	FindByID(sprintID uint) (*model.Sprint, error)
	FindAllByProject(projectID uint) ([]model.Sprint, error)
	Update(sprint *model.Sprint) error
	Delete(sprintID uint) error
}

type CommentRepositoryInterface interface {
	Create(comment *model.Comment) error
	FindByID(commentID uint) (*model.Comment, error)
	FindAllByTask(taskID uint) ([]model.Comment, error)
	Delete(commentID uint) error
}

type LabelRepositoryInterface interface {
	Create(label *model.Label) error
	FindByID(labelID uint) (*model.Label, error)
	FindAllByProject(projectID uint) ([]model.Label, error)
	Delete(labelID uint) error
	AttachToTask(taskID, labelID uint) error
	DetachFromTask(taskID, labelID uint) error
	FindLabelsByTask(taskID uint) ([]model.Label, error)
}

type ActivityLogRepositoryInterface interface {
	Create(log *model.ActivityLog) error
	FindAllByTask(taskID uint) ([]model.ActivityLog, error)
}

type RefreshTokenRepositoryInterface interface {
	Create(token *model.RefreshToken) error
	FindByHash(tokenHash string) (*model.RefreshToken, error)
	Revoke(tokenID uint) error
}