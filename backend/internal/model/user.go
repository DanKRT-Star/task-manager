package model

import "time"

type User struct {
	UserID         uint      `gorm:"primaryKey" json:"userId"`
	UserName       string    `gorm:"not null" json:"userName"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	HashedPassword string    `gorm:"not null" json:"-"`
	Tasks          []Task    `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}