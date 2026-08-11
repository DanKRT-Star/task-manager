package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type ProjectMemberRepository struct {
	DB *gorm.DB
}

func NewProjectMemberRepository(db *gorm.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{DB: db}
}

func (r *ProjectMemberRepository) AddMember(member *model.ProjectMember) error {
	return r.DB.Create(member).Error
}

func (r *ProjectMemberRepository) FindMember(projectID, userID uint) (*model.ProjectMember, error) {
	var member model.ProjectMember
	err := r.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// FindMembersByProject lấy danh sách member, tự query User thủ công (không dùng Preload vì User dùng gorm:"-")
func (r *ProjectMemberRepository) FindMembersByProject(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	if err := r.DB.Where("project_id = ?", projectID).Find(&members).Error; err != nil {
		return nil, err
	}

	for i := range members {
		var user model.User
		if err := r.DB.First(&user, members[i].UserID).Error; err == nil {
			members[i].User = user
		}
	}

	return members, nil
}

func (r *ProjectMemberRepository) RemoveMember(projectID, userID uint) error {
	return r.DB.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.ProjectMember{}).Error
}