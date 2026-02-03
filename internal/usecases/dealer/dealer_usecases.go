package dealer

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositroies"
)

// จัดการร้านค้า
type DealerUsecase struct {
	DealerRepo *repositroies.DealerRepository
	CarRepo    *repositroies.CarRepository
	ReviewRepo *repositroies.ReviewRepository
}

// สมัคร / สร้างร้านค้า
func (u *DealerUsecase) CreateDealer(dealer *entities.Dealer) error {
	return u.DealerRepo.Create(dealer)
}

// ดูร้านค้าทั้งหมด (หน้าเว็บ / แอดมิน)
func (u *DealerUsecase) GetAllDealers(dealers *[]*entities.Dealer) error {
	return u.DealerRepo.FindAll(dealers)
}

// ดูร้านค้ารายเดียว
func (u *DealerUsecase) GetDealerByID(id uint, dealer *entities.Dealer) error {
	return u.DealerRepo.FindByID(id, dealer)
}

// ดูร้านค้าจาก User ID
func (u *DealerUsecase) GetDealerByUserID(userID uint, dealer *entities.Dealer) error {
	return u.DealerRepo.FindByUserID(userID, dealer)
}

// แก้ไขข้อมูลร้าน
func (u *DealerUsecase) UpdateDealer(dealer *entities.Dealer) error {
	return u.DealerRepo.Update(dealer)
}

// GetDealerStats retrieves dealer rating and review statistics
func (u *DealerUsecase) GetDealerStats(dealerID uint) (map[string]interface{}, error) {
	// Get all reviews for this dealer
	var reviews []*entities.Review
	if err := u.ReviewRepo.FindByDealerID(dealerID, &reviews); err != nil {
		return map[string]interface{}{
			"total_rating":        0,
			"review_count":        0,
			"rating_distribution": make(map[int]int),
		}, nil
	}

	// Calculate statistics
	totalRating := 0.0
	ratingDistribution := make(map[int]int)

	for _, review := range reviews {
		totalRating += float64(review.Rating)
		ratingDistribution[review.Rating]++
	}

	averageRating := 0.0
	if len(reviews) > 0 {
		averageRating = totalRating / float64(len(reviews))
	}

	return map[string]interface{}{
		"total_rating":        averageRating,
		"review_count":        len(reviews),
		"rating_distribution": ratingDistribution,
	}, nil
}

// GetDealerCars retrieves all cars from a dealer
func (u *DealerUsecase) GetDealerCars(dealerID uint) ([]*entities.Car, error) {
	var cars []*entities.Car
	if err := u.CarRepo.FindByDealerID(dealerID, &cars); err != nil {
		return nil, err
	}
	return cars, nil
}

// GetDealerReviews retrieves all reviews for a dealer
func (u *DealerUsecase) GetDealerReviews(dealerID uint) ([]*entities.Review, error) {
	var reviews []*entities.Review
	if err := u.ReviewRepo.FindByDealerID(dealerID, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}
