package repository

import (
	"time"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (r *RefreshTokenRepository) Create(token *model.RefreshToken) error {
	return r.DB.Create(token).Error
}

func (r *RefreshTokenRepository) FindByHash(tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.DB.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(tokenID uint) error {
	now := time.Now()
	return r.DB.Model(&model.RefreshToken{}).Where("refresh_token_id = ?", tokenID).Update("revoked_at", now).Error
}