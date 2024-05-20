package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"time"
)

type UserService interface {
	Register(userRegisterRequest *web.UserRegisterRequest) (domain.Users, error)
	Login(userLoginRequest *web.UserLoginRequest) (token *string, err error)
	Update(userUpdateRequest web.UserUpdateRequest) (domain.Users, error)
	Delete(username string) error
	GetAll() ([]domain.Users, error)
	GetByUsername(username string) (domain.Users, error)
	CheckAvailable(username string) bool
}

type userServiceImpl struct {
	repository.UserRepository
	repository.SessionRepository
	*validator.Validate
}

func NewUserService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository, validate *validator.Validate) UserService {
	return &userServiceImpl{userRepository, sessionRepository, validate}
}

func (u *userServiceImpl) Register(userRegisterRequest *web.UserRegisterRequest) (domain.Users, error) {
	err := u.Validate.Struct(userRegisterRequest)
	if err != nil {
		return domain.Users{}, err
	}

	hasPassword, err := hashPassword(userRegisterRequest.Password)
	if err != nil {
		return domain.Users{}, errors.New("hashing password failed")
	}

	if u.CheckAvailable(userRegisterRequest.Username) {
		return domain.Users{}, errors.New("username is already taken")
	}

	user, err := u.UserRepository.Add(domain.Users{
		FullName: userRegisterRequest.FullName,
		Username: userRegisterRequest.Username,
		Password: hasPassword,
		Role:     userRegisterRequest.Role,
	})
	if err != nil {
		return domain.Users{}, errors.New("user creation failed")
	}

	return user, nil
}

func (u *userServiceImpl) Login(userLoginRequest *web.UserLoginRequest) (token *string, err error) {
	err = u.Validate.Struct(userLoginRequest)
	if err != nil {
		return nil, err
	}

	user, err := u.UserRepository.GetByUsername(userLoginRequest.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := checkPasswordHash(userLoginRequest.Password, user.Password)
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

	_, err = u.SessionRepository.Add(session)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (u *userServiceImpl) Update(userUpdateRequest web.UserUpdateRequest) (domain.Users, error) {
	err := u.Validate.Struct(userUpdateRequest)
	if err != nil {
		return domain.Users{}, err
	}

	hasPassword, err := hashPassword(userUpdateRequest.Password)
	if err != nil {
		return domain.Users{}, errors.New("hashing password failed")
	}

	user, err := u.UserRepository.Update(domain.Users{
		Username: userUpdateRequest.Username,
		FullName: userUpdateRequest.FullName,
		Password: hasPassword,
		Role:     userUpdateRequest.Role,
	})
	if err != nil {
		return domain.Users{}, err
	}

	return user, nil
}

func (u *userServiceImpl) Delete(username string) error {
	return u.UserRepository.Delete(username)
}

func (u *userServiceImpl) GetAll() ([]domain.Users, error) {
	users, err := u.UserRepository.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].Password = "-"
	}
	return users, nil
}

func (u *userServiceImpl) GetByUsername(username string) (domain.Users, error) {
	user, err := u.UserRepository.GetByUsername(username)
	if err != nil {
		return domain.Users{}, err
	}
	user.Password = "-"
	return user, nil
}

func (u *userServiceImpl) CheckAvailable(username string) bool {
	_, err := u.UserRepository.GetByUsername(username)
	if err != nil {
		return false
	}
	return true
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
