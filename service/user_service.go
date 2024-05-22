package service

import (
	"github.com/golang-jwt/jwt"
	"inventory-management-system/helper"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"net/http"
	"time"
)

type UserService interface {
	Register(userRegisterRequest *web.UserRegisterRequest) web.ErrorResponse
	Login(userLoginRequest *web.UserLoginRequest) (*string, web.ErrorResponse)
	Update(userUpdateRequest web.UserUpdateRequest) web.ErrorResponse
	Delete(username string) web.ErrorResponse
	GetAll() ([]domain.Users, web.ErrorResponse)
	GetByUsername(username string) (domain.Users, web.ErrorResponse)
	CheckAvailable(username string) bool
}

type userServiceImpl struct {
	repository.UserRepository
	repository.SessionRepository
}

func NewUserService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository) UserService {
	return &userServiceImpl{userRepository, sessionRepository}
}

func (u *userServiceImpl) Register(userRegisterRequest *web.UserRegisterRequest) web.ErrorResponse {
	hasPassword, err := helper.HashPassword(userRegisterRequest.Password)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: "failed to hash password",
		}
	}

	if u.CheckAvailable(userRegisterRequest.Username) {
		return web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "username is already taken",
		}
	}

	err = u.UserRepository.Add(domain.Users{
		FullName: userRegisterRequest.FullName,
		Username: userRegisterRequest.Username,
		Password: hasPassword,
		Role:     userRegisterRequest.Role,
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (u *userServiceImpl) Login(userLoginRequest *web.UserLoginRequest) (*string, web.ErrorResponse) {
	user, err := u.UserRepository.GetByUsername(userLoginRequest.Username)
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "user not found",
		}
	}

	result := helper.CheckPasswordHash(userLoginRequest.Password, user.Password)
	if !result {
		return nil, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "username or password is wrong",
		}
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
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	session := domain.Sessions{
		Username:  user.Username,
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	err = u.SessionRepository.Add(session)
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return &tokenString, web.ErrorResponse{}
}

func (u *userServiceImpl) Update(userUpdateRequest web.UserUpdateRequest) web.ErrorResponse {
	hasPassword, err := helper.HashPassword(userUpdateRequest.Password)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: "failed to hash password",
		}
	}

	err = u.UserRepository.Update(domain.Users{
		Username: userUpdateRequest.Username,
		FullName: userUpdateRequest.FullName,
		Password: hasPassword,
		Role:     userUpdateRequest.Role,
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (u *userServiceImpl) Delete(username string) web.ErrorResponse {
	err := u.UserRepository.Delete(username)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (u *userServiceImpl) GetAll() ([]domain.Users, web.ErrorResponse) {
	users, err := u.UserRepository.GetAll()
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	if len(users) == 0 {
		return nil, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "users not found",
		}
	}

	return users, web.ErrorResponse{}
}

func (u *userServiceImpl) GetByUsername(username string) (domain.Users, web.ErrorResponse) {
	user, err := u.UserRepository.GetByUsername(username)
	if err != nil {
		return domain.Users{}, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return user, web.ErrorResponse{}
}

func (u *userServiceImpl) CheckAvailable(username string) bool {
	_, err := u.UserRepository.GetByUsername(username)
	if err != nil {
		return false
	}
	return true
}
