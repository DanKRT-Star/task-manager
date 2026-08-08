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

func (m *MockTaskRepository) FindAll(userID uint, status, sort string, page, limit int) ([]model.Task, int64, error) {
	args := m.Called(userID, status, sort, page, limit)
	return args.Get(0).([]model.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) Update(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(taskID, userID uint) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}