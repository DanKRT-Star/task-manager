package service

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository giả lập UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// MockTaskRepository giả lập TaskRepositoryInterface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(taskID, userID uint) (*model.Task, error) {
	args := m.Called(taskID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByIDOnly(taskID uint) (*model.Task, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll(userID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	args := m.Called(userID, status, sort, page, limit)
	return args.Get(0).([]model.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) FindAllByProject(projectID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	args := m.Called(projectID, status, sort, page, limit)
	return args.Get(0).([]model.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) Update(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(taskID uint) error {
	args := m.Called(taskID)
	return args.Error(0)
}

// MockActivityLogRepository giả lập ActivityLogRepositoryInterface
type MockActivityLogRepository struct {
	mock.Mock
}

func (m *MockActivityLogRepository) Create(log *model.ActivityLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockActivityLogRepository) FindAllByTask(taskID uint) ([]model.ActivityLog, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ActivityLog), args.Error(1)
}

// MockProjectMemberRepository giả lập ProjectMemberRepositoryInterface
type MockProjectMemberRepository struct {
	mock.Mock
}

func (m *MockProjectMemberRepository) AddMember(member *model.ProjectMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) FindMember(projectID, userID uint) (*model.ProjectMember, error) {
	args := m.Called(projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) FindMembersByProject(projectID uint) ([]model.ProjectMember, error) {
	args := m.Called(projectID)
	return args.Get(0).([]model.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) RemoveMember(projectID, userID uint) error {
	args := m.Called(projectID, userID)
	return args.Error(0)
}

// MockProjectRepository giả lập ProjectRepositoryInterface
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) FindByID(projectID uint) (*model.Project, error) {
	args := m.Called(projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepository) FindAllByUser(userID uint) ([]model.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(projectID uint) error {
	args := m.Called(projectID)
	return args.Error(0)
}

// MockEpicRepository giả lập EpicRepositoryInterface
type MockEpicRepository struct {
	mock.Mock
}

func (m *MockEpicRepository) Create(epic *model.Epic) error {
	args := m.Called(epic)
	return args.Error(0)
}

func (m *MockEpicRepository) FindByID(epicID uint) (*model.Epic, error) {
	args := m.Called(epicID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Epic), args.Error(1)
}

func (m *MockEpicRepository) FindAllByProject(projectID uint) ([]model.Epic, error) {
	args := m.Called(projectID)
	return args.Get(0).([]model.Epic), args.Error(1)
}

func (m *MockEpicRepository) Update(epic *model.Epic) error {
	args := m.Called(epic)
	return args.Error(0)
}

func (m *MockEpicRepository) Delete(epicID uint) error {
	args := m.Called(epicID)
	return args.Error(0)
}

// MockMilestoneRepository giả lập MilestoneRepositoryInterface
type MockMilestoneRepository struct {
	mock.Mock
}

func (m *MockMilestoneRepository) Create(milestone *model.Milestone) error {
	args := m.Called(milestone)
	return args.Error(0)
}

func (m *MockMilestoneRepository) FindByID(milestoneID uint) (*model.Milestone, error) {
	args := m.Called(milestoneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Milestone), args.Error(1)
}

func (m *MockMilestoneRepository) FindAllByProject(projectID uint) ([]model.Milestone, error) {
	args := m.Called(projectID)
	return args.Get(0).([]model.Milestone), args.Error(1)
}

func (m *MockMilestoneRepository) Update(milestone *model.Milestone) error {
	args := m.Called(milestone)
	return args.Error(0)
}

func (m *MockMilestoneRepository) Delete(milestoneID uint) error {
	args := m.Called(milestoneID)
	return args.Error(0)
}

// MockSprintRepository giả lập SprintRepositoryInterface
type MockSprintRepository struct {
	mock.Mock
}

func (m *MockSprintRepository) Create(sprint *model.Sprint) error {
	args := m.Called(sprint)
	return args.Error(0)
}

func (m *MockSprintRepository) FindByID(sprintID uint) (*model.Sprint, error) {
	args := m.Called(sprintID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Sprint), args.Error(1)
}

func (m *MockSprintRepository) FindAllByProject(projectID uint) ([]model.Sprint, error) {
	args := m.Called(projectID)
	return args.Get(0).([]model.Sprint), args.Error(1)
}

func (m *MockSprintRepository) Update(sprint *model.Sprint) error {
	args := m.Called(sprint)
	return args.Error(0)
}

func (m *MockSprintRepository) Delete(sprintID uint) error {
	args := m.Called(sprintID)
	return args.Error(0)
}
