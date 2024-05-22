package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"inventory-management-system/helper"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"time"
)

type UserService interface {
	Register(userRegisterRequest *web.UserRegisterRequest) error
	Login(userLoginRequest *web.UserLoginRequest) (token *string, err error)
	Update(userUpdateRequest web.UserUpdateRequest) error
	Delete(username string) error
	GetAll() ([]domain.Users, error)
	GetByUsername(username string) (domain.Users, error)
	CheckAvailable(username string) bool
}

type userServiceImpl struct {
	repository.UserRepository
	repository.SessionRepository
}

func NewUserService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository) UserService {
	return &userServiceImpl{userRepository, sessionRepository}
}

func (u *userServiceImpl) Register(userRegisterRequest *web.UserRegisterRequest) error {
	hasPassword, err := helper.HashPassword(userRegisterRequest.Password)
	if err != nil {
		return errors.New("hashing password failed")
	}

	if u.CheckAvailable(userRegisterRequest.Username) {
		return errors.New("username is already taken")
	}

	return u.UserRepository.Add(domain.Users{
		FullName: userRegisterRequest.FullName,
		Username: userRegisterRequest.Username,
		Password: hasPassword,
		Role:     userRegisterRequest.Role,
	})
}

func (u *userServiceImpl) Login(userLoginRequest *web.UserLoginRequest) (token *string, err error) {
	user, err := u.UserRepository.GetByUsername(userLoginRequest.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := helper.CheckPasswordHash(userLoginRequest.Password, user.Password)
	if !result {
		return nil, errors.New("email or password is wrong")
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &domain.JwtCustomClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(domain.JwtKey)
	if err != nil {
		return nil, err
	}

	session := domain.Sessions{
		Username:  user.Username,
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	err = u.SessionRepository.Add(session)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (u *userServiceImpl) Update(userUpdateRequest web.UserUpdateRequest) error {
	hasPassword, err := helper.HashPassword(userUpdateRequest.Password)
	if err != nil {
		return errors.New("hashing password failed")
	}

	return u.UserRepository.Update(domain.Users{
		Username: userUpdateRequest.Username,
		FullName: userUpdateRequest.FullName,
		Password: hasPassword,
		Role:     userUpdateRequest.Role,
	})
}

func (u *userServiceImpl) Delete(username string) error {
	return u.UserRepository.Delete(username)
}

func (u *userServiceImpl) GetAll() ([]domain.Users, error) {
	return u.UserRepository.GetAll()
}

func (u *userServiceImpl) GetByUsername(username string) (domain.Users, error) {
	return u.UserRepository.GetByUsername(username)
}

func (u *userServiceImpl) CheckAvailable(username string) bool {
	_, err := u.UserRepository.GetByUsername(username)
	if err != nil {
		return false
	}
	return true
}
