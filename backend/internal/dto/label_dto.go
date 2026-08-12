package dto

type CreateLabelRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Color string `json:"color" validate:"omitempty,hexcolor"`
}