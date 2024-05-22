package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type UserRepository interface {
	Add(user domain.Users) error
	Update(user domain.Users) error
	Delete(username string) error
	GetAll() ([]domain.Users, error)
	GetByUsername(username string) (domain.Users, error)
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (u *userRepositoryImpl) Add(user domain.Users) error {
	return u.DB.Create(&user).Error
}

func (u *userRepositoryImpl) Update(user domain.Users) error {
	newUser := domain.Users{}
	if err := u.DB.First(&newUser, "username = ?", user.Username).Error; err != nil {
		return err
	}

	newUser.FullName = user.FullName
	newUser.Password = user.Password
	newUser.Role = user.Role
	return u.DB.Save(&newUser).Error
}

func (u *userRepositoryImpl) Delete(username string) error {
	user, err := u.GetByUsername(username)
	if err != nil {
		return err
	}

	return u.DB.Delete(&user).Error
}

func (u *userRepositoryImpl) GetAll() ([]domain.Users, error) {
	var users []domain.Users
	if err := u.DB.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (u *userRepositoryImpl) GetByUsername(username string) (domain.Users, error) {
	user := domain.Users{}
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return domain.Users{}, err
	}

	return user, nil
}
