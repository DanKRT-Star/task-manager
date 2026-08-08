package dto

import "github.com/DanKRT-Star/task-manager/internal/model"

type CreateTaskRequest struct {
	Title       string           `json:"title" validate:"required,min=1,max=200"`
	Description string           `json:"description" validate:"max=1000"`
	Status      model.TaskStatus `json:"status"`
	Deadline    string           `json:"deadline"`
}

type UpdateTaskRequest struct {
	Title       string           `json:"title" validate:"omitempty,min=1,max=200"`
	Description string           `json:"description" validate:"omitempty,max=1000"`
	Status      model.TaskStatus `json:"status"`
	Deadline    string           `json:"deadline"`
}

type TaskListResponse struct {
	Data  []model.Task `json:"data"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}