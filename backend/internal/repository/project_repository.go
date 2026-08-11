package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.DB.Create(project).Error
}

func (r *ProjectRepository) FindByID(projectID uint) (*model.Project, error) {
	var project model.Project
	err := r.DB.First(&project, projectID).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindAllByUser trả về project mà user là owner HOẶC member
func (r *ProjectRepository) FindAllByUser(userID uint) ([]model.Project, error) {
	var projects []model.Project
	err := r.DB.
		Joins("JOIN project_members ON project_members.project_id = projects.project_id").
		Where("project_members.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Update(project *model.Project) error {
	return r.DB.Save(project).Error
}

func (r *ProjectRepository) Delete(projectID uint) error {
	return r.DB.Delete(&model.Project{}, projectID).Error
}