package repositories

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type FavoriteRepository struct{ DB *gorm.DB }

func (r *FavoriteRepository) Create(fav *entities.Favorite) error {
	return r.DB.Create(fav).Error
}

func (r *FavoriteRepository) Exists(userID, carID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Favorite{}).
		Where("user_id = ? AND car_id = ?", userID, carID).
		Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepository) FindByUserID(userID uint, favs interface{}) error {
	return r.DB.
		Preload("Car").
		Preload("Car.CarImages").
		Preload("Car.Dealer"). // Load Dealer info for display
		Where("user_id = ?", userID).
		Find(favs).Error
}

func (r *FavoriteRepository) Delete(userID, carID uint) error {
	return r.DB.Where("user_id = ? AND car_id = ?", userID, carID).Delete(&entities.Favorite{}).Error
}

func (r *FavoriteRepository) Toggle(userID, carID uint) (string, error) {
	var status string
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var existing entities.Favorite
		// Use Unscoped to find soft-deleted records that might trigger unique constraint
		err := tx.Unscoped().Where("user_id = ? AND car_id = ?", userID, carID).First(&existing).Error

		if err == nil {
			// Record exists
			if existing.DeletedAt.Valid {
				// If soft-deleted, restore it
				if err := tx.Unscoped().Model(&existing).Update("deleted_at", nil).Error; err != nil {
					return err
				}
				status = "added"
			} else {
				// If active, soft delete it
				if err := tx.Delete(&existing).Error; err != nil {
					return err
				}
				status = "removed"
			}
			return nil
		}

		// Handle not found
		if err != nil && err.Error() != "record not found" {
			return err
		}

		// Not found at all, create new
		newFav := entities.Favorite{
			UserID: userID,
			CarID:  carID,
		}
		if err := tx.Create(&newFav).Error; err != nil {
			return err
		}
		status = "added"
		return nil
	})

	return status, err
}
