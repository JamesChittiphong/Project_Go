package repositroies

import (
	"gorm.io/gorm"
)

type CarImageRepository struct{ DB *gorm.DB }

func (r *CarImageRepository) Create(img interface{}) error {
	return r.DB.Create(img).Error
}

func (r *CarImageRepository) FindByCarID(carID uint, imgs interface{}) error {
	return r.DB.Where("car_id = ?", carID).Find(imgs).Error
}

func (r *CarImageRepository) Delete(id uint) error {
	return r.DB.Delete(&map[string]interface{}{}, id).Error
}
