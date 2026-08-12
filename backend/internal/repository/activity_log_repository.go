package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	DB *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

func (r *ActivityLogRepository) Create(log *model.ActivityLog) error {
	return r.DB.Create(log).Error
}

func (r *ActivityLogRepository) FindAllByTask(taskID uint) ([]model.ActivityLog, error) {
	var logs []model.ActivityLog
	err := r.DB.Where("task_id = ?", taskID).Order("created_at desc").Find(&logs).Error
	return logs, err
}