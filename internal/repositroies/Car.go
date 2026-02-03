package repositroies

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type CarRepository struct {
	DB *gorm.DB
}

func (r *CarRepository) Create(car *entities.Car) error {
	return r.DB.Create(car).Error
}

func (r *CarRepository) FindAll(cars *[]*entities.Car) error {
	return r.DB.
		Preload("CarImages").
		Preload("Dealer").
		Preload("Dealer.User").
		Find(cars).Error
}

func (r *CarRepository) FindByID(id uint, car *entities.Car) error {
	return r.DB.
		Preload("CarImages").
		Preload("Dealer").
		Preload("Dealer.User").
		First(car, id).Error
}

func (r *CarRepository) Update(car *entities.Car) error {
	return r.DB.Save(car).Error
}

func (r *CarRepository) Delete(id uint) error {
	return r.DB.Delete(&entities.Car{}, id).Error
}

func (r *CarRepository) FindByDealerID(dealerID uint, cars *[]*entities.Car) error {
	return r.DB.
		Preload("CarImages").
		Preload("Dealer").
		Preload("Dealer.User").
		Where("dealer_id = ?", dealerID).
		Find(cars).Error
}
