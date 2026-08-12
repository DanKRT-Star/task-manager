package model

import "time"

type Comment struct {
	CommentID uint      `gorm:"primaryKey" json:"commentId"`
	TaskID    uint      `gorm:"index;not null" json:"taskId"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}