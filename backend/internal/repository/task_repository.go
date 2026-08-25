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
	if err := r.DB.Create(task).Error; err != nil {
		return err
	}
	r.populateProject(task)
	return nil
}

// FindByID giữ nguyên logic cũ - dùng cho task cá nhân (không thuộc project)
func (r *TaskRepository) FindByID(taskID, userID uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.Preload("Labels").Where("task_id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// FindByIDOnly lấy task theo ID, không check quyền - việc check quyền chuyển cho service xử lý
func (r *TaskRepository) FindByIDOnly(taskID uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.Preload("Labels").First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	r.populateProject(&task)
	return &task, nil
}

func (r *TaskRepository) FindMyTasks(userID uint, status string, sort string, page, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.DB.Model(&model.Task{}).
		Where("(project_id IS NULL AND user_id = ?) OR assignee_id = ?", userID, userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if sort == "deadline_desc" {
		query = query.Order("deadline desc")
	} else {
		query = query.Order("deadline asc")
	}

	offset := (page - 1) * limit
	if err := query.Preload("Labels").Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	// Project dùng gorm:"-" (tránh xung đột association 2 chiều), nên tự nạp thủ công
	// thay vì Preload — chỉ query những project thực sự cần (task có ProjectID khác nil).
	projectCache := make(map[uint]*model.Project)
	for i := range tasks {
		if tasks[i].ProjectID == nil {
			continue
		}
		pid := *tasks[i].ProjectID
		if cached, ok := projectCache[pid]; ok {
			tasks[i].Project = cached
			continue
		}
		var project model.Project
		if err := r.DB.First(&project, pid).Error; err == nil {
			projectCache[pid] = &project
			tasks[i].Project = &project
		}
	}

	return tasks, total, nil
}

// FindAllByProject lấy toàn bộ task thuộc 1 project (quyền xem đã check ở service trước khi gọi)
func (r *TaskRepository) FindAllByProject(projectID uint, status string, sort string, page, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.DB.Model(&model.Task{}).Where("project_id = ?", projectID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if sort == "deadline_desc" {
		query = query.Order("deadline desc")
	} else {
		query = query.Order("deadline asc")
	}

	offset := (page - 1) * limit
	err := query.Preload("Labels").Limit(limit).Offset(offset).Find(&tasks).Error

	return tasks, total, err
}

func (r *TaskRepository) Update(task *model.Task) error {
	if err := r.DB.Save(task).Error; err != nil {
		return err
	}
	r.populateProject(task)
	return nil
}

func (r *TaskRepository) Delete(taskID uint) error {
	return r.DB.Delete(&model.Task{}, taskID).Error
}

func (r *TaskRepository) populateProject(task *model.Task) {
	if task.ProjectID == nil {
		return
	}
	var project model.Project
	if err := r.DB.First(&project, *task.ProjectID).Error; err == nil {
		task.Project = &project
	}
}