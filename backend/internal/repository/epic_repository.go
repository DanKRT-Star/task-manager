package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type EpicRepository struct {
	DB *gorm.DB
}

func NewEpicRepository(db *gorm.DB) *EpicRepository {
	return &EpicRepository{DB: db}
}

func (r *EpicRepository) Create(epic *model.Epic) error {
	return r.DB.Create(epic).Error
}

func (r *EpicRepository) FindByID(epicID uint) (*model.Epic, error) {
	var epic model.Epic
	err := r.DB.First(&epic, epicID).Error
	if err != nil {
		return nil, err
	}
	return &epic, nil
}

func (r *EpicRepository) FindAllByProject(projectID uint) ([]model.Epic, error) {
	var epics []model.Epic
	err := r.DB.Where("project_id = ?", projectID).Find(&epics).Error
	return epics, err
}

func (r *EpicRepository) Update(epic *model.Epic) error {
	return r.DB.Save(epic).Error
}

func (r *EpicRepository) Delete(epicID uint) error {
	return r.DB.Delete(&model.Epic{}, epicID).Error
}