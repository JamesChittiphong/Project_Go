package repositroies

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *entities.RefreshToken) error
	Find(token string) (*entities.RefreshToken, error)
	Revoke(token string) error
}

type refreshRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshRepo{db}
}

func (r *refreshRepo) Create(token *entities.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshRepo) Find(token string) (*entities.RefreshToken, error) {
	var rt entities.RefreshToken
	err := r.db.Where("token = ? AND revoked = false", token).First(&rt).Error
	return &rt, err
}

func (r *refreshRepo) Revoke(token string) error {
	return r.db.Model(&entities.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}
