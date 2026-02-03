package car

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositroies"
	"errors"
	"time"
)

type CarUsecase struct {
	CarRepo      *repositroies.CarRepository
	DealerRepo   *repositroies.DealerRepository
	LeadRepo     *repositroies.LeadRepository
	FavoriteRepo *repositroies.FavoriteRepository
}

// ---------- Core ----------

func (u *CarUsecase) CreateCar(car *entities.Car) error {
	if car.DealerID == 0 {
		return errors.New("dealer_id is required")
	}

	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(car.DealerID, &dealer); err != nil {
		return errors.New("dealer not found")
	}

	return u.CarRepo.Create(car)
}

func (u *CarUsecase) GetAllCars(cars *[]*entities.Car) error {
	return u.CarRepo.FindAll(cars)
}

func (u *CarUsecase) GetCarDetail(id uint) (*entities.Car, error) {
	var car entities.Car
	if err := u.CarRepo.FindByID(id, &car); err != nil {
		return nil, err
	}
	return &car, nil
}

func (u *CarUsecase) UpdateCar(car *entities.Car) error {
	if car.ID == 0 {
		return errors.New("car_id is required")
	}
	return u.CarRepo.Update(car)
}

func (u *CarUsecase) DeleteCarByUser(carID uint, userID uint) error {
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByUserID(userID, &dealer); err != nil {
		return errors.New("forbidden")
	}

	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}

	if car.DealerID != dealer.ID {
		return errors.New("forbidden")
	}

	return u.CarRepo.Delete(carID)
}

// ---------- Business ----------

func (u *CarUsecase) SetStatus(carID uint, status string) error {
	if status == "" {
		return errors.New("status is required")
	}

	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}

	car.Status = status
	return u.CarRepo.Update(&car)
}

func (u *CarUsecase) RecordContact(carID uint, dealerID uint, via string) error {
	if dealerID == 0 {
		return errors.New("dealer_id is required")
	}

	if via != "call" && via != "line" {
		return errors.New("invalid contact method")
	}

	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}

	if via == "call" {
		car.CallCount++
	} else {
		car.LineCount++
	}
	car.LeadCount++

	if err := u.CarRepo.Update(&car); err != nil {
		return err
	}

	return u.LeadRepo.Create(&entities.Lead{
		CarID:      carID,
		DealerID:   dealerID,
		ContactVia: via,
	})
}

func (u *CarUsecase) GetStats(carID uint) (*entities.Car, error) {
	return u.GetCarDetail(carID)
}

func (u *CarUsecase) PromoteCar(carID uint, days int) error {
	if days <= 0 {
		days = 7
	}

	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}

	until := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	car.IsPromoted = true
	car.PromotedUntil = &until

	return u.CarRepo.Update(&car)
}

func (u *CarUsecase) GetCarsByDealer(dealerID uint, cars *[]*entities.Car) error {
	if dealerID == 0 {
		return errors.New("dealer_id is required")
	}
	return u.CarRepo.FindByDealerID(dealerID, cars)
}
