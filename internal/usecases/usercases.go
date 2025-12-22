package usecases

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/middleware"
	"Backend_Go/internal/repositroies"
	"Backend_Go/utils"
	"errors"
	"time"
)

// ==================== CAR USECASE ====================
// ใช้กับหน้าเว็บลูกค้า + ร้านค้า

type CarUsecase struct {
	CarRepo      *repositroies.CarRepository
	ImageRepo    *repositroies.CarImageRepository
	DealerRepo   *repositroies.DealerRepository
	LeadRepo     *repositroies.LeadRepository
	FavoriteRepo *repositroies.FavoriteRepository
}

// Update an existing car
func (u *CarUsecase) UpdateCar(car *entities.Car) error {
	return u.CarRepo.Update(car)
}

// Add additional images to a car
func (u *CarUsecase) AddImages(images []interface{}) error {
	for _, img := range images {
		if err := u.ImageRepo.Create(img); err != nil {
			return err
		}
	}
	return nil
}

// Set car status (available, contacted, sold)
func (u *CarUsecase) SetStatus(carID uint, status string) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	car.Status = status
	return u.CarRepo.Update(&car)
}

// Record a contact (call/line) and create a Lead entry
func (u *CarUsecase) RecordContact(carID uint, dealerID uint, via string) error {
	// increment counters on car
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	switch via {
	case "call":
		car.CallCount++
	case "line":
		car.LineCount++
	default:
		// unknown contact method still counts as a lead
	}
	car.LeadCount++
	if err := u.CarRepo.Update(&car); err != nil {
		return err
	}

	// save lead record
	lead := &entities.Lead{CarID: carID, DealerID: dealerID, ContactVia: via}
	return u.LeadRepo.Create(lead)
}

// Get contact statistics for a car
func (u *CarUsecase) GetStats(carID uint) (*entities.Car, error) {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return nil, err
	}
	return &car, nil
}

// Promote a car for a duration (days)
func (u *CarUsecase) PromoteCar(carID uint, days int) error {
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return err
	}
	now := time.Now()
	until := now.Add(time.Duration(days) * 24 * time.Hour)
	car.IsPromoted = true
	car.PromotedUntil = &until
	return u.CarRepo.Update(&car)
}

type AuthUsecase interface {
	RegisterUser(name, email, password, role string) error
	RegisterDealer(
		name, email, password,
		shopName, phone, lineID string,
	) error

	Login(email, password string) (accessToken, refreshToken string, err error)
	Refresh(refreshToken string) (string, error)
	Logout(refreshToken string) error
	GetUser(id uint) (*entities.User, error)
}

type authUsecase struct {
	userRepo    repositroies.UserRepository
	dealerRepo  *repositroies.DealerRepository
	refreshRepo repositroies.RefreshTokenRepository
}

func NewAuthUsecase(
	userRepo repositroies.UserRepository,
	dealerRepo *repositroies.DealerRepository,
	refreshRepo repositroies.RefreshTokenRepository,
) AuthUsecase {
	return &authUsecase{userRepo, dealerRepo, refreshRepo}
}

func (u *authUsecase) RegisterUser(name, email, password, role string) error {
	hash, _ := utils.HashPassword(password)

	user := &entities.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
	}

	return u.userRepo.Create(user)
}

func (u *authUsecase) Login(email, password string) (string, string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || !utils.CheckPassword(user.Password, password) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, _ := middleware.GenerateToken(user.ID, user.Role)

	refreshToken := middleware.GenerateRefreshToken()

	rt := &entities.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	u.refreshRepo.Create(rt)

	return accessToken, refreshToken, nil
}

func (u *authUsecase) Refresh(token string) (string, error) {
	rt, err := u.refreshRepo.Find(token)
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return "", errors.New("invalid refresh token")
	}

	user, _ := u.userRepo.FindByID(rt.UserID)
	return middleware.GenerateToken(user.ID, user.Role)
}

func (u *authUsecase) Logout(token string) error {
	return u.refreshRepo.Revoke(token)
}

func (u *authUsecase) GetUser(id uint) (*entities.User, error) {
	return u.userRepo.FindByID(id)
}

// สร้างรถใหม่ (ร้านค้า)
func (u *CarUsecase) CreateCar(car *entities.Car, images []interface{}) error {
	// 1. บันทึกรถ
	if err := u.CarRepo.Create(car); err != nil {
		return err
	}

	// 2. บันทึกรูปรถ
	for _, img := range images {
		if err := u.ImageRepo.Create(img); err != nil {
			return err
		}
	}
	return nil
}

// ดึงรถทั้งหมด (ลูกค้า)
func (u *CarUsecase) GetAllCars(cars *[]*entities.Car) error {
	return u.CarRepo.FindAll(cars)
}

// ดูรายละเอียดรถ (ลูกค้า)
func (u *CarUsecase) GetCarDetail(carID uint, car *entities.Car, images interface{}) error {
	if err := u.CarRepo.FindByID(carID, car); err != nil {
		return err
	}
	return u.ImageRepo.FindByCarID(carID, images)
}

// ลบรถ (ร้านค้า)
func (u *CarUsecase) DeleteCar(carID uint) error {
	return u.CarRepo.Delete(carID)
}

// ==================== LEAD USECASE ====================
// ใช้เมื่อ ลูกค้ากดติดต่อร้าน

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

// ==================== FAVORITE USECASE ====================
// รถที่ลูกค้าชอบ

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

// ==================== REVIEW USECASE ====================
// รีวิวร้านค้า

type ReviewUsecase struct {
	ReviewRepo *repositroies.ReviewRepository
	DealerRepo *repositroies.DealerRepository
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

// ==================== ADMIN USECASE ====================
// สำหรับแอดมินระบบ

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

// USER USECASE

// จัดการผู้ใช้งาน (สมัคร / โปรไฟล์ / แอดมิน)
type UserUsecase struct {
	UserRepo repositroies.UserRepository
}

// สร้างผู้ใช้ใหม่
func (u *UserUsecase) CreateUser(user *entities.User) error {
	return u.UserRepo.Create(user)
}

// ดูข้อมูลผู้ใช้
func (u *UserUsecase) GetUserByID(id uint) (*entities.User, error) {
	return u.UserRepo.FindByID(id)
}

// แก้ไขข้อมูลผู้ใช้
func (u *UserUsecase) UpdateUser(user *entities.User) error {
	return u.UserRepo.Update(user)
}

// ลบผู้ใช้ (แอดมิน)
func (u *UserUsecase) DeleteUser(id uint) error {
	return u.UserRepo.Delete(id)
}

// DEALER USECASE

// จัดการร้านค้า
type DealerUsecase struct {
	DealerRepo *repositroies.DealerRepository
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

// แก้ไขข้อมูลร้าน
func (u *DealerUsecase) UpdateDealer(dealer *entities.Dealer) error {
	return u.DealerRepo.Update(dealer)
}

// CAR IMAGE USECASE

// จัดการรูปรถ (แยกจาก Car)
type CarImageUsecase struct {
	ImageRepo *repositroies.CarImageRepository
	CarRepo   *repositroies.CarRepository
}

// เพิ่มรูปรถ
func (u *CarImageUsecase) AddImage(image interface{}, carID uint) error {
	// business rule: ต้องมีรถก่อน
	var car entities.Car
	if err := u.CarRepo.FindByID(carID, &car); err != nil {
		return errors.New("ไม่พบรถ")
	}
	return u.ImageRepo.Create(image)
}

// ดูรูปรถทั้งหมด
func (u *CarImageUsecase) GetImagesByCar(carID uint, images interface{}) error {
	return u.ImageRepo.FindByCarID(carID, images)
}

// ลบรูปรถ
func (u *CarImageUsecase) DeleteImage(imageID uint) error {
	return u.ImageRepo.Delete(imageID)
}

func (u *authUsecase) RegisterDealer(
	name, email, password,
	shopName, phone, lineID string,
) error {

	// 1. hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// 2. create user (role = dealer)
	user := &entities.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     "dealer",
	}

	if err := u.userRepo.Create(user); err != nil {
		return err
	}

	// 3. create dealer profile
	dealer := &entities.Dealer{
		UserID:     user.ID,
		ShopName:   shopName,
		Phone:      phone,
		LineID:     lineID,
		IsApproved: false,
	}

	return u.dealerRepo.Create(dealer)
}
