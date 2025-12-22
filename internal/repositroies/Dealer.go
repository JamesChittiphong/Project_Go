package repositroies

import (
	"gorm.io/gorm"
)

type DealerRepository struct{ DB *gorm.DB }

func (r *DealerRepository) Create(dealer interface{}) error {
	return r.DB.Create(dealer).Error
}
func (r *DealerRepository) FindAll(dealers interface{}) error {
	return r.DB.Find(dealers).Error
}
func (r *DealerRepository) FindByID(id uint, dealer interface{}) error {
	return r.DB.First(dealer, id).Error
}
func (r *DealerRepository) Update(dealer interface{}) error {
	return r.DB.Save(dealer).Error
}
func (r *DealerRepository) Delete(id uint) error {
	return r.DB.Delete(&map[string]interface{}{}, id).Error
}

func (r *DealerRepository) FindByUserID(userID uint, dealer interface{}) error {
	return r.DB.Where("user_id = ?", userID).First(dealer).Error
}
