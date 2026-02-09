package admin

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
)

type AdminUsecase struct {
	UserRepo   repositories.UserRepository
	DealerRepo *repositories.DealerRepository
	ReportRepo *repositories.ReportRepository
	CarRepo    *repositories.CarRepository
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

// Admin Get All Cars (including pending, deleted_requested)
func (u *AdminUsecase) GetAllCars(cars *[]*entities.Car) error {
	return u.CarRepo.FindAll(cars)
}

// Approve or reject a dealer
func (u *AdminUsecase) SetDealerApproval(dealerID uint, approve bool) error {
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err != nil {
		return err
	}
	if approve {
		dealer.Status = "approved"
		dealer.IsApproved = true
	} else {
		dealer.Status = "rejected"
		dealer.IsApproved = false
	}
	return u.DealerRepo.Update(&dealer)
}

// RejectDealer explicitly
func (u *AdminUsecase) RejectDealer(dealerID uint) error {
	return u.SetDealerApproval(dealerID, false)
}

// ApproveCar
func (u *AdminUsecase) ApproveCar(carID uint) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	car.Status = "approved"
	car.IsHidden = false
	return u.CarRepo.Update(&car)
}

// RejectCar
func (u *AdminUsecase) RejectCar(carID uint, reason string) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	car.Status = "rejected"
	car.ViolationReason = reason
	car.IsHidden = true
	return u.CarRepo.Update(&car)
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

// BanUser disables a user account
func (u *AdminUsecase) BanUser(userID uint) error {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.IsActive = false
	return u.UserRepo.Update(user)
}

// UnbanUser enables a user account
func (u *AdminUsecase) UnbanUser(userID uint) error {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.IsActive = true
	return u.UserRepo.Update(user)
}

// SuspendDealer unapproves a dealer and bans their user account
func (u *AdminUsecase) SuspendDealer(dealerID uint) error {
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err != nil {
		return err
	}

	// 1. Unapprove dealer
	dealer.Status = "suspended"
	dealer.IsApproved = false
	if err := u.DealerRepo.Update(&dealer); err != nil {
		return err
	}

	// 2. Ban the underlying user
	return u.BanUser(dealer.UserID)
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
