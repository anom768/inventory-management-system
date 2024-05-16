package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model"
)

type UserRepository interface {
	Add(user model.Users) (model.Users, error)
	Update(user model.Users) (model.Users, error)
	Delete(username string) error
	GetAll() ([]model.Users, error)
	GetByUsername(username string) (model.Users, error)
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (u *userRepositoryImpl) Add(user model.Users) (model.Users, error) {
	if err := u.DB.Create(&user).Error; err != nil {
		return model.Users{}, err
	}

	return user, nil
}

func (u *userRepositoryImpl) Update(user model.Users) (model.Users, error) {
	newUser := model.Users{}
	if err := u.DB.First(&newUser, "username = ?", user.Username).Error; err != nil {
		return model.Users{}, err
	}

	newUser.Password = user.Password
	newUser.Role = user.Role
	if err := u.DB.Save(&user).Error; err != nil {
		return model.Users{}, err
	}

	return newUser, nil
}

func (u *userRepositoryImpl) Delete(username string) error {
	user, err := u.GetByUsername(username)
	if err != nil {
		return err
	}

	if err := u.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepositoryImpl) GetAll() ([]model.Users, error) {
	var users []model.Users
	if err := u.DB.Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (u *userRepositoryImpl) GetByUsername(username string) (model.Users, error) {
	user := model.Users{}
	if err := u.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return model.Users{}, err
	}

	return user, nil
}
