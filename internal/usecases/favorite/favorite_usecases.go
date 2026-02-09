package favorite

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
	"errors"
)

type FavoriteUsecase struct {
	FavoriteRepo *repositories.FavoriteRepository
	CarRepo      *repositories.CarRepository
}

// ToggleFavourite เพิ่มหรือลบรถที่ชอบ
func (u *FavoriteUsecase) ToggleFavorite(userID, carID uint) (string, error) {
	// ตรวจสอบว่ารถมีอยู่
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return "", errors.New("ไม่พบรถ")
	}

	return u.FavoriteRepo.Toggle(userID, carID)
}

// ดูรถที่ชอบทั้งหมด
func (u *FavoriteUsecase) GetFavoritesByUser(userID uint, favs interface{}) error {
	return u.FavoriteRepo.FindByUserID(userID, favs)
}

// ลบรถที่ชอบ
func (u *FavoriteUsecase) RemoveFavorite(userID, carID uint) error {
	return u.FavoriteRepo.Delete(userID, carID)
}
