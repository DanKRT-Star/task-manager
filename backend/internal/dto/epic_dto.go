package dto

type CreateEpicRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type UpdateEpicRequest struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=200"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}