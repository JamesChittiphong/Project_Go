package review

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
	"errors"
)

type ReviewUsecase struct {
	ReviewRepo *repositories.ReviewRepository
	DealerRepo *repositories.DealerRepository
}

// ลูกค้ารีวิวร้าน
func (u *ReviewUsecase) CreateReview(review interface{}, dealerID uint) error {
	// ตรวจสอบร้าน
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err != nil {
		return errors.New("ไม่พบร้านค้า")
	}
	return u.ReviewRepo.Create(review)
}

// ดูรีวิวร้าน
func (u *ReviewUsecase) GetReviewsByDealer(dealerID uint, reviews interface{}) error {
	return u.ReviewRepo.FindByDealerID(dealerID, reviews)
}
