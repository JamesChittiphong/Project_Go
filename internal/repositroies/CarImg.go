package repositroies

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type CarImageRepository struct{ DB *gorm.DB }

func (r *CarImageRepository) Create(img *entities.CarImage) error {
	return r.DB.Create(img).Error
}

func (r *CarImageRepository) FindByCarID(carID uint) ([]*entities.CarImage, error) {
	var images []*entities.CarImage
	err := r.DB.Where("car_id = ?", carID).Find(&images).Error
	return images, err
}

func (r *CarImageRepository) FindByID(id uint) (*entities.CarImage, error) {
	var image *entities.CarImage
	err := r.DB.First(&image, id).Error
	return image, err
}

func (r *CarImageRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.CarImage{}, id).Error
}

func (r *CarImageRepository) DeleteByCarID(carID uint) error {
	return r.DB.Where("car_id = ?", carID).Delete(&entities.CarImage{}).Error
}

func (r *CarImageRepository) Update(img *entities.CarImage) error {
	return r.DB.Save(img).Error
}
