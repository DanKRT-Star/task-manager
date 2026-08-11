package model

import "time"

type ProjectRole string

const (
	RoleOwner  ProjectRole = "owner"
	RoleMember ProjectRole = "member"
)

func (r ProjectRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleMember:
		return true
	default:
		return false
	}
}

type Project struct {
	ProjectID   uint            `gorm:"primaryKey" json:"projectId"`
	Name        string          `gorm:"not null" json:"name"`
	Description string          `json:"description"`
	Deadline    *time.Time      `json:"deadline,omitempty"`
	OwnerID     uint            `gorm:"index;not null" json:"ownerId"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE;" json:"members,omitempty"`
	Tasks       []Task          `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL;" json:"tasks,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type ProjectMember struct {
	ProjectMemberID uint        `gorm:"primaryKey" json:"projectMemberId"`
	ProjectID       uint        `gorm:"index;not null" json:"projectId"`
	UserID          uint        `gorm:"index;not null" json:"userId"`
	Role            ProjectRole `gorm:"type:varchar(20);not null;check:role IN ('owner','member')" json:"role"`
	User            User        `gorm:"-" json:"user,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
}