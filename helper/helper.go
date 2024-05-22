package helper

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"net/http"
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

func ReadFromRequestBody(c *gin.Context, v any) {
	err := c.ShouldBindJSON(v)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid body request",
		})
		return
	}
}

//func Validate(c *gin.Context, v any) {
//	err := validator.Validate.Struct(validator.Validate{}, v)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
//			Code:    http.StatusBadRequest,
//			Status:  "status bad request",
//			Message: "validation error: " + err.Error(),
//		})
//		return
//	}
//}

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
