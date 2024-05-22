package helper

import (
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/model/domain"
	"inventory-management-system/repository"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func RegisterAdmin(userRepository repository.UserRepository) {
	err := userRepository.Add(domain.Users{
		FullName: "Administrator",
		Username: "admin",
		Password: "admin123",
		Role:     "admin",
	})
	if err != nil {
		panic(err)
	}
}
