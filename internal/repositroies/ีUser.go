package repositroies

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uint) (*entities.User, error)
	FindAll(users interface{}) error
	Update(user *entities.User) error
	Delete(id uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepo) FindAll(users interface{}) error {
	return r.db.Find(users).Error
}

func (r *userRepo) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&entities.User{}, id).Error
}
