package dto

import (
	"time"

	"github.com/DanKRT-Star/task-manager/internal/model"
)

type CreateSprintRequest struct {
	Name      string     `json:"name" validate:"required,min=1,max=200"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

type UpdateSprintRequest struct {
	Name      string           `json:"name" validate:"omitempty,min=1,max=200"`
	Status    model.SprintStatus `json:"status"`
	StartDate *time.Time       `json:"startDate"`
	EndDate   *time.Time       `json:"endDate"`
}