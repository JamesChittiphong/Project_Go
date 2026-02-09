package carimage

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
	"errors"
)

type CarImageUsecase struct {
	CarImageRepo *repositories.CarImageRepository
	CarRepo      *repositories.CarRepository
}

// CreateCarImage creates a new car image
func (u *CarImageUsecase) CreateCarImage(image *entities.CarImage) error {
	// Validate required fields
	if image.CarID == 0 {
		return errors.New("car_id is required")
	}

	if image.ImageURL == "" {
		return errors.New("image_url is required")
	}

	// Verify car exists
	var car entities.Car
	if err := u.CarRepo.FindByID(image.CarID, &car); err != nil {
		return errors.New("car not found")
	}

	return u.CarImageRepo.Create(image)
}

// GetCarImages retrieves all images for a car
func (u *CarImageUsecase) GetCarImages(carID uint) ([]*entities.CarImage, error) {
	if carID == 0 {
		return nil, errors.New("car_id is required")
	}

	// Verify car exists
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return nil, errors.New("car not found")
	}

	return u.CarImageRepo.FindByCarID(carID)
}

// GetCarImage retrieves a specific image
func (u *CarImageUsecase) GetCarImage(imageID uint) (*entities.CarImage, error) {
	if imageID == 0 {
		return nil, errors.New("image_id is required")
	}

	return u.CarImageRepo.FindByID(imageID)
}

// UpdateCarImage updates a car image
func (u *CarImageUsecase) UpdateCarImage(image *entities.CarImage) error {
	if image.ID == 0 {
		return errors.New("image_id is required")
	}

	if image.ImageURL == "" {
		return errors.New("image_url is required")
	}

	// Verify image exists
	_, err := u.CarImageRepo.FindByID(image.ID)
	if err != nil {
		return errors.New("image not found")
	}

	return u.CarImageRepo.Update(image)
}

// DeleteCarImage deletes a car image
func (u *CarImageUsecase) DeleteCarImage(imageID uint) error {
	if imageID == 0 {
		return errors.New("image_id is required")
	}

	// Verify image exists
	_, err := u.CarImageRepo.FindByID(imageID)
	if err != nil {
		return errors.New("image not found")
	}

	return u.CarImageRepo.Delete(imageID)
}

// DeleteCarImages deletes all images for a car
func (u *CarImageUsecase) DeleteCarImages(carID uint) error {
	if carID == 0 {
		return errors.New("car_id is required")
	}

	// Verify car exists
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return errors.New("car not found")
	}

	return u.CarImageRepo.DeleteByCarID(carID)
}
