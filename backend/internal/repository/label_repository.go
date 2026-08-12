package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type LabelRepository struct {
	DB *gorm.DB
}

func NewLabelRepository(db *gorm.DB) *LabelRepository {
	return &LabelRepository{DB: db}
}

func (r *LabelRepository) Create(label *model.Label) error {
	return r.DB.Create(label).Error
}

func (r *LabelRepository) FindByID(labelID uint) (*model.Label, error) {
	var label model.Label
	err := r.DB.First(&label, labelID).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *LabelRepository) FindAllByProject(projectID uint) ([]model.Label, error) {
	var labels []model.Label
	err := r.DB.Where("project_id = ?", projectID).Find(&labels).Error
	return labels, err
}

func (r *LabelRepository) Delete(labelID uint) error {
	return r.DB.Delete(&model.Label{}, labelID).Error
}

func (r *LabelRepository) AttachToTask(taskID, labelID uint) error {
	return r.DB.Create(&model.TaskLabel{TaskID: taskID, LabelID: labelID}).Error
}

func (r *LabelRepository) DetachFromTask(taskID, labelID uint) error {
	return r.DB.Where("task_id = ? AND label_id = ?", taskID, labelID).Delete(&model.TaskLabel{}).Error
}

func (r *LabelRepository) FindLabelsByTask(taskID uint) ([]model.Label, error) {
	var labels []model.Label
	err := r.DB.
		Joins("JOIN task_labels ON task_labels.label_id = labels.label_id").
		Where("task_labels.task_id = ?", taskID).
		Find(&labels).Error
	return labels, err
}