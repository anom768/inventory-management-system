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

func ReadFromRequestBody(c *gin.Context, v any) error {
	err := c.ShouldBindJSON(v)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid body request"))
		return err
	}
	return nil
}

func RegisterAdmin(handleRepository repository.HandlerRepository) {
	pwd, err := HashPassword("admin123")
	if err != nil {
		panic(err)
	}

	err = handleRepository.Add(&domain.Users{
		FullName: "Administrator",
		Username: "administrator",
		Password: pwd,
		Role:     "admin",
	})
	if err != nil {
		panic(err)
	}
}
