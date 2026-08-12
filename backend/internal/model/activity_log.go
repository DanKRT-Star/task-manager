package model

import "time"

type ActivityAction string

const (
	ActionCreated            ActivityAction = "created"
	ActionStatusChanged      ActivityAction = "status_changed"
	ActionAssigneeChanged    ActivityAction = "assignee_changed"
	ActiconDeadlineChanged   ActivityAction = "deadline_changed"
	ActionDescriptionUpdated ActivityAction = "description_updated"
	ActionUpdated            ActivityAction = "updated"
)

type ActivityLog struct {
	ActivityLogID uint           `gorm:"primaryKey" json:"activityLogId"`
	TaskID        uint           `gorm:"index;not null" json:"taskId"`
	UserID        uint           `gorm:"index;not null" json:"userId"`
	Action        ActivityAction `gorm:"type:varchar(30);not null" json:"action"`
	Detail        string         `json:"detail"`
	CreatedAt     time.Time      `json:"createdAt"`
}
