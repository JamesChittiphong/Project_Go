package user

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
)

// จัดการผู้ใช้งาน (สมัคร / โปรไฟล์ / แอดมิน)
type UserUsecase struct {
	UserRepo repositories.UserRepository
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
