package model

import "time"

type RefreshToken struct {
	RefreshTokenID uint       `gorm:"primaryKey" json:"-"`
	UserID         uint       `gorm:"index;not null" json:"-"`
	TokenHash      string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt      time.Time  `json:"-"`
	RevokedAt      *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"-"`
}