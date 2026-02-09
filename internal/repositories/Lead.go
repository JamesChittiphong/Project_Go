package repositories

import (
	"gorm.io/gorm"
)

type LeadRepository struct{ DB *gorm.DB }

func (r *LeadRepository) Create(lead interface{}) error {
	return r.DB.Create(lead).Error
}

func (r *LeadRepository) FindByDealerID(dealerID uint, leads interface{}) error {
	return r.DB.Where("dealer_id = ?", dealerID).Find(leads).Error
}
