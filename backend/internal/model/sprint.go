package model

import "time"

type SprintStatus string

const (
	SprintPlanned   SprintStatus = "planned"
	SprintActive    SprintStatus = "active"
	SprintCompleted SprintStatus = "completed"
)

func (s SprintStatus) IsValid() bool {
	switch s {
	case SprintPlanned, SprintActive, SprintCompleted:
		return true
	default:
		return false
	}
}

type Sprint struct {
	SprintID  uint         `gorm:"primaryKey" json:"sprintId"`
	ProjectID uint         `gorm:"index;not null" json:"projectId"`
	Name      string       `gorm:"not null" json:"name"`
	Status    SprintStatus `gorm:"type:varchar(20);default:'planned';check:status IN ('planned','active','completed')" json:"status"`
	StartDate *time.Time   `json:"startDate,omitempty"`
	EndDate   *time.Time   `json:"endDate,omitempty"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}