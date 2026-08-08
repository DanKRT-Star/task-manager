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
	FindAll(userID uint, status string, sort string, page, limit int) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(taskID, userID uint) error
}