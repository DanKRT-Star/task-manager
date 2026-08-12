package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type SprintRepository struct {
	DB *gorm.DB
}

func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{DB: db}
}

func (r *SprintRepository) Create(sprint *model.Sprint) error {
	return r.DB.Create(sprint).Error
}

func (r *SprintRepository) FindByID(sprintID uint) (*model.Sprint, error) {
	var sprint model.Sprint
	err := r.DB.First(&sprint, sprintID).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (r *SprintRepository) FindAllByProject(projectID uint) ([]model.Sprint, error) {
	var sprints []model.Sprint
	err := r.DB.Where("project_id = ?", projectID).Find(&sprints).Error
	return sprints, err
}

func (r *SprintRepository) Update(sprint *model.Sprint) error {
	return r.DB.Save(sprint).Error
}

func (r *SprintRepository) Delete(sprintID uint) error {
	return r.DB.Delete(&model.Sprint{}, sprintID).Error
}