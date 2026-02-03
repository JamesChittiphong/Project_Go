package favorite

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositroies"
	"errors"
)

type FavoriteUsecase struct {
	FavoriteRepo *repositroies.FavoriteRepository
	CarRepo      *repositroies.CarRepository
}

// เพิ่มรถที่ชอบ
func (u *FavoriteUsecase) AddFavorite(fav interface{}, carID uint) error {
	// ตรวจสอบว่ารถมีอยู่
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return errors.New("ไม่พบรถ")
	}
	return u.FavoriteRepo.Create(fav)
}

// ดูรถที่ชอบทั้งหมด
func (u *FavoriteUsecase) GetFavoritesByUser(userID uint, favs interface{}) error {
	return u.FavoriteRepo.FindByUserID(userID, favs)
}

// ลบรถที่ชอบ
func (u *FavoriteUsecase) RemoveFavorite(userID, carID uint) error {
	return u.FavoriteRepo.Delete(userID, carID)
}
