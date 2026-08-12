package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type MilestoneRepository struct {
	DB *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) *MilestoneRepository {
	return &MilestoneRepository{DB: db}
}

func (r *MilestoneRepository) Create(milestone *model.Milestone) error {
	return r.DB.Create(milestone).Error
}

func (r *MilestoneRepository) FindByID(milestoneID uint) (*model.Milestone, error) {
	var milestone model.Milestone
	err := r.DB.First(&milestone, milestoneID).Error
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (r *MilestoneRepository) FindAllByProject(projectID uint) ([]model.Milestone, error) {
	var milestones []model.Milestone
	err := r.DB.Where("project_id = ?", projectID).Find(&milestones).Error
	return milestones, err
}

func (r *MilestoneRepository) Update(milestone *model.Milestone) error {
	return r.DB.Save(milestone).Error
}

func (r *MilestoneRepository) Delete(milestoneID uint) error {
	return r.DB.Delete(&model.Milestone{}, milestoneID).Error
}