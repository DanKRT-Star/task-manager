package dto

import "time"

type CreateProjectRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateProjectRequest struct {
	Name        string     `json:"name" validate:"omitempty,min=1,max=200"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	Deadline    *time.Time `json:"deadline"`
}

type AddMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ProjectMemberResponse struct {
	ProjectMemberID uint   `json:"projectMemberId"`
	UserID          uint   `json:"userId"`
	UserName        string `json:"userName"`
	Email           string `json:"email"`
	Role            string `json:"role"`
}