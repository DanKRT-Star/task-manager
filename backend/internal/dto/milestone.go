package dto

import "time"

type CreateMilestoneRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	DueDate     *time.Time `json:"dueDate"`
}

type UpdateMilestoneRequest struct {
	Title       string     `json:"title" validate:"omitempty,min=1,max=200"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	DueDate     *time.Time `json:"dueDate"`
}