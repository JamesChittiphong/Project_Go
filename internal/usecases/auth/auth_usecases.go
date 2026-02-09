package auth

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/middleware"
	"Backend_Go/internal/repositories"
	"Backend_Go/utils"
	"errors"
	"log"
	"time"
)

type authUsecase struct {
	userRepo    repositories.UserRepository
	dealerRepo  *repositories.DealerRepository
	refreshRepo repositories.RefreshTokenRepository
}

func NewAuthUsecase(
	userRepo repositories.UserRepository,
	dealerRepo *repositories.DealerRepository,
	refreshRepo repositories.RefreshTokenRepository,
) AuthUsecase {
	return &authUsecase{userRepo, dealerRepo, refreshRepo}
}

func (u *authUsecase) RegisterUser(
	name, email, phone, password, role string,
) error {

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		return err
	}

	user := &entities.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: hash,
		Role:     role,
		IsActive: true,
	}

	log.Printf("Registering user - Email: %s, Role: %s, Password Length: %d", email, role, len(password))
	log.Printf("Password Hash: %s", hash)
	err = u.userRepo.Create(user)
	if err != nil {
		log.Printf("User creation error: %v", err)
		return err
	}

	log.Printf("User registered successfully - Email: %s", email)
	return nil
}

func (u *authUsecase) Login(email, password string) (string, string, *entities.User, error) {
	log.Printf("Login attempt - Email: %s, Password Length: %d", email, len(password))

	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("User not found - Email: %s, Error: %v", email, err)
		return "", "", nil, errors.New("user not found")
	}

	log.Printf("User found - Email: %s, Stored Hash: %s", user.Email, user.Password)
	log.Printf("Checking password - Input: %s (len:%d)", password, len(password))

	isPasswordValid := utils.CheckPassword(user.Password, password)
	log.Printf("Password check result: %v", isPasswordValid)

	if !isPasswordValid {
		log.Printf("Invalid password - Email: %s", email)
		return "", "", nil, errors.New("invalid password")
	}

	log.Printf("Login successful - Email: %s, UserID: %d", email, user.ID)

	accessToken, _ := middleware.GenerateToken(user.ID, user.Role)

	refreshToken := middleware.GenerateRefreshToken()

	rt := &entities.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	u.refreshRepo.Create(rt)

	return accessToken, refreshToken, user, nil
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

func (u *authUsecase) RegisterDealer(
	name, email, phone, password,
	shopName, lineID,
	address, province, latitude, longitude string,
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
		Phone:    phone,
		Password: hash,
		Role:     "dealer",
		IsActive: true,
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
		Address:    address,
		Province:   province,
		Latitude:   latitude,
		Longitude:  longitude,
		Status:     "pending",
		IsApproved: false,
	}

	return u.dealerRepo.Create(dealer)
}

func (u *authUsecase) GetDealerByUserID(userID uint, dealer *entities.Dealer) error {
	return u.dealerRepo.FindByUserID(userID, dealer)
}
