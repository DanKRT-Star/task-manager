package model

import "time"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// IsValid kiểm tra giá trị status có hợp lệ không, dùng ở tầng service trước khi lưu DB
func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	TaskID      uint       `gorm:"primaryKey" json:"taskId"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','in_progress','done')" json:"status"`
	ProjectID   *uint      `gorm:"index" json:"projectId,omitempty"`
	EpicID      *uint      `gorm:"index" json:"epicId,omitempty"`
	MilestoneID *uint      `gorm:"index" json:"milestoneId,omitempty"`
	SprintID    *uint      `gorm:"index" json:"sprintId,omitempty"`
	UserID      uint       `gorm:"index;not null" json:"userId"`
	AssigneeID  *uint      `gorm:"index" json:"assigneeId,omitempty"`
	Deadline    time.Time  `json:"deadline"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
