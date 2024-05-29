package service

import (
	"github.com/golang-jwt/jwt"
	"inventory-management-system/helper"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"time"
)

type UserService interface {
	Register(userRegisterRequest *web.UserRegisterRequest) web.ErrorResponse
	Login(userLoginRequest *web.UserLoginRequest) (*string, web.ErrorResponse)
	Logout() web.ErrorResponse
	Update(userUpdateRequest web.UserUpdateRequest) web.ErrorResponse
	Delete(username string) web.ErrorResponse
	GetAll() ([]domain.Users, web.ErrorResponse)
	GetByUsername(username string) (domain.Users, web.ErrorResponse)
	CheckAvailable(username string) bool
}

type userServiceImpl struct {
	repository.HandlerRepository
}

func NewUserService(handlerRepository repository.HandlerRepository) UserService {
	return &userServiceImpl{handlerRepository}
}

func (u *userServiceImpl) Register(userRegisterRequest *web.UserRegisterRequest) web.ErrorResponse {
	hasPassword, err := helper.HashPassword(userRegisterRequest.Password)
	if err != nil {
		return web.NewInternalServerErrorError("failed to hash password")
	}

	if u.CheckAvailable(userRegisterRequest.Username) {
		return web.NewBadRequestError("username is already taken")
	}

	err = u.HandlerRepository.Add(&domain.Users{
		FullName: userRegisterRequest.FullName,
		Username: userRegisterRequest.Username,
		Password: hasPassword,
		Role:     userRegisterRequest.Role,
	})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (u *userServiceImpl) Login(userLoginRequest *web.UserLoginRequest) (*string, web.ErrorResponse) {
	user := domain.Users{}
	err := u.HandlerRepository.GetByUsername(userLoginRequest.Username, &user)
	if err != nil {
		return nil, web.NewNotFoundError("user not found")
	}

	result := helper.CheckPasswordHash(userLoginRequest.Password, user.Password)
	if !result {
		return nil, web.NewBadRequestError("invalid username or password")
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
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	session := domain.Sessions{
		Username:  user.Username,
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	err = u.HandlerRepository.Add(&session)
	if err != nil {
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	return &tokenString, nil
}

func (u *userServiceImpl) Logout() web.ErrorResponse {
	//TODO implement me
	panic("implement me")
}

func (u *userServiceImpl) Update(userUpdateRequest web.UserUpdateRequest) web.ErrorResponse {
	if !u.CheckAvailable(userUpdateRequest.Username) {
		return web.NewNotFoundError("user not found")
	}

	hasPassword, err := helper.HashPassword(userUpdateRequest.Password)
	if err != nil {
		return web.NewInternalServerErrorError("failed to hash password")
	}

	err = u.HandlerRepository.UpdateByUsername(userUpdateRequest.Username, &domain.Users{
		Username: userUpdateRequest.Username,
		FullName: userUpdateRequest.FullName,
		Password: hasPassword,
		Role:     userUpdateRequest.Role,
	})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (u *userServiceImpl) Delete(username string) web.ErrorResponse {
	if !u.CheckAvailable(username) {
		return web.NewNotFoundError("user not found")
	}

	err := u.HandlerRepository.DeleteByUsername(username, &domain.Users{})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (u *userServiceImpl) GetAll() ([]domain.Users, web.ErrorResponse) {
	users := []domain.Users{}
	err := u.HandlerRepository.GetAll(&users)
	if err != nil {
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	if len(users) == 0 {
		return nil, web.NewNotFoundError("user not found")
	}

	for i, _ := range users {
		users[i].Password = "-"
	}

	return users, nil
}

func (u *userServiceImpl) GetByUsername(username string) (domain.Users, web.ErrorResponse) {
	user := domain.Users{}
	err := u.HandlerRepository.GetByUsername(username, &user)
	if err != nil {
		return user, web.NewNotFoundError("user not found")
	}

	user.Password = "-"
	return user, nil
}

func (u *userServiceImpl) CheckAvailable(username string) bool {
	err := u.HandlerRepository.GetByUsername(username, &domain.Users{})
	if err != nil {
		return false
	}
	return true
}
