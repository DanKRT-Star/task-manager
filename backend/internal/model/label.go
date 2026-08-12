package model

import "time"

type Label struct {
	LabelID   uint      `gorm:"primaryKey" json:"labelId"`
	ProjectID uint      `gorm:"index;not null" json:"projectId"`
	Name      string    `gorm:"not null" json:"name"`
	Color     string    `gorm:"not null;default:'#6b7280'" json:"color"`
	CreatedAt time.Time `json:"createdAt"`
}

// TaskLabel là bảng trung gian many-to-many giữa Task và Label
type TaskLabel struct {
	TaskID  uint `gorm:"primaryKey" json:"taskId"`
	LabelID uint `gorm:"primaryKey" json:"labelId"`
}
