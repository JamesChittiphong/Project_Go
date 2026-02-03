package admin

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositroies"
)

type AdminUsecase struct {
	UserRepo   repositroies.UserRepository
	DealerRepo *repositroies.DealerRepository
	ReportRepo *repositroies.ReportRepository
	CarRepo    *repositroies.CarRepository
}

// ดูผู้ใช้ทั้งหมด
func (u *AdminUsecase) GetAllUsers(users interface{}) error {
	return u.UserRepo.FindAll(users)
}

// ดูร้านค้าทั้งหมด
func (u *AdminUsecase) GetAllDealers(dealers *[]*entities.Dealer) error {
	return u.DealerRepo.FindAll(dealers)
}

// ดูรายงานปัญหา
func (u *AdminUsecase) GetAllReports(reports interface{}) error {
	return u.ReportRepo.FindAll(reports)
}

// Approve or reject a dealer
func (u *AdminUsecase) SetDealerApproval(dealerID uint, approve bool) error {
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err != nil {
		return err
	}
	dealer.IsApproved = approve
	return u.DealerRepo.Update(&dealer)
}

// Hide or unhide a car
func (u *AdminUsecase) SetCarHidden(carID uint, hide bool) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	car.IsHidden = hide
	return u.CarRepo.Update(&car)
}

// Flag a car as violating rules with a reason
func (u *AdminUsecase) FlagCar(carID uint, reason string) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	car.Flagged = true
	car.ViolationReason = reason
	// optionally hide when flagged
	car.IsHidden = true
	return u.CarRepo.Update(&car)
}

// Admin delete car
func (u *AdminUsecase) DeleteCar(carID uint) error {
	return u.CarRepo.Delete(carID)
}

// IsAdmin checks if the user has admin role and is active
func (u *AdminUsecase) IsAdmin(userID uint) (bool, error) {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	if user.Role == "admin" && user.IsActive {
		return true, nil
	}
	return false, nil
}
