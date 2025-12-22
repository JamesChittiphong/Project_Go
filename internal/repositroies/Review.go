package repositroies

import (
	"gorm.io/gorm"
)

type ReviewRepository struct{ DB *gorm.DB }

func (r *ReviewRepository) Create(review interface{}) error {
	return r.DB.Create(review).Error
}

func (r *ReviewRepository) FindByDealerID(dealerID uint, reviews interface{}) error {
	return r.DB.Where("dealer_id = ?", dealerID).Find(reviews).Error
}
