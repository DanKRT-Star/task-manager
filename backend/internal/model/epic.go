package model

import "time"

type Epic struct {
	EpicID      uint      `gorm:"primaryKey" json:"epicId"`
	ProjectID   uint       `gorm:"index;not null" json:"projectId"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}