package auth

import "Backend_Go/internal/entities"

type AuthUsecase interface {
	RegisterUser(name, email, phone, password, role string) error
	RegisterDealer(
		name, email, phone, password,
		shopName, lineID,
		address, province, latitude, longitude string,
	) error

	Login(email, password string) (accessToken, refreshToken string, user *entities.User, err error)
	Refresh(refreshToken string) (string, error)
	Logout(refreshToken string) error
	GetUser(id uint) (*entities.User, error)
	GetDealerByUserID(userID uint, dealer *entities.Dealer) error
}
