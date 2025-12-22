package repositroies

import (
	"gorm.io/gorm"
)

type FavoriteRepository struct{ DB *gorm.DB }

func (r *FavoriteRepository) Create(fav interface{}) error {
	return r.DB.Create(fav).Error
}

func (r *FavoriteRepository) FindByUserID(userID uint, favs interface{}) error {
	return r.DB.Where("user_id = ?", userID).Find(favs).Error
}

func (r *FavoriteRepository) Delete(userID, carID uint) error {
	return r.DB.Where("user_id = ? AND car_id = ?", userID, carID).Delete(&map[string]interface{}{}).Error
}
