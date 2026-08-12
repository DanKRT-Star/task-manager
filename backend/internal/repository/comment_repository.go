package repository

import (
	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *CommentRepository) FindByID(commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.DB.First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) FindAllByTask(taskID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.DB.Where("task_id = ?", taskID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) Delete(commentID uint) error {
	return r.DB.Delete(&model.Comment{}, commentID).Error
}