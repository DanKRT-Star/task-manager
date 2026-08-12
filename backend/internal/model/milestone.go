package model

import "time"

type Milestone struct {
	MilestoneID uint       `gorm:"primaryKey" json:"milestoneId"`
	ProjectID   uint       `gorm:"index;not null" json:"projectId"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}