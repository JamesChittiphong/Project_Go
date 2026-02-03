package lend

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositroies"
	"errors"
)

type LeadUsecase struct {
	LeadRepo   *repositroies.LeadRepository
	CarRepo    *repositroies.CarRepository
	DealerRepo *repositroies.DealerRepository
}

// ลูกค้าส่งข้อมูลติดต่อร้าน
func (u *LeadUsecase) CreateLead(lead interface{}, carID uint, dealerID uint) error {
	// ตรวจสอบว่ารถมีอยู่จริง
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return errors.New("ไม่พบรถที่ต้องการติดต่อ")
	}

	// ตรวจสอบร้านค้า
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err != nil {
		return errors.New("ไม่พบร้านค้า")
	}

	// บันทึก lead
	return u.LeadRepo.Create(lead)
}

// ร้านค้าดูรายชื่อ Lead ของตัวเอง
func (u *LeadUsecase) GetLeadsByDealer(dealerID uint, leads interface{}) error {
	return u.LeadRepo.FindByDealerID(dealerID, leads)
}
