package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return r.DB.Create(task).Error
}

// FindByID chỉ trả về task nếu đúng userID sở hữu (chặn truy cập chéo user ngay ở tầng query)
func (r *TaskRepository) FindByID(taskID, userID uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// FindAll hỗ trợ filter theo status, sort theo deadline, và pagination
func (r *TaskRepository) FindAll(userID uint, status string, sort string, page, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.DB.Model(&model.Task{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Đếm tổng số bản ghi trước khi phân trang (dùng cho frontend hiển thị số trang)
	query.Count(&total)

	if sort == "deadline_desc" {
		query = query.Order("deadline desc")
	} else {
		query = query.Order("deadline asc")
	}

	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Find(&tasks).Error

	return tasks, total, err
}

func (r *TaskRepository) Update(task *model.Task) error {
	return r.DB.Save(task).Error
}

func (r *TaskRepository) Delete(taskID, userID uint) error {
	return r.DB.Where("task_id = ? AND user_id = ?", taskID, userID).Delete(&model.Task{}).Error
}