package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"inventory-management-system/helper"
	"inventory-management-system/model/web"
	"inventory-management-system/service"
	"net/http"
	"time"
)

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	GetByUsername(c *gin.Context)
}

type userControllerImpl struct {
	service.UserService
	*validator.Validate
}

func NewUserController(userService service.UserService, validate *validator.Validate) UserController {
	return &userControllerImpl{userService, validate}
}

func (u *userControllerImpl) Register(c *gin.Context) {
	var userRegisterRequest web.UserRegisterRequest
	if err := helper.ReadFromRequestBody(c, &userRegisterRequest); err != nil {
		return
	}

	err := u.Validate.Struct(userRegisterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "validation error: " + err.Error(),
		})
		return
	}

	errResponse := u.UserService.Register(&userRegisterRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusCreated, web.SuccessResponseMessage{
		Code:    http.StatusCreated,
		Status:  "status created",
		Message: "register user success",
	})
}

func (u *userControllerImpl) Login(c *gin.Context) {
	var userLoginRequest web.UserLoginRequest
	if err := helper.ReadFromRequestBody(c, &userLoginRequest); err != nil {
		return
	}

	err := u.Validate.Struct(userLoginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "validation error: " + err.Error(),
		})
		return
	}

	tokenString, errResponse := u.UserService.Login(&userLoginRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   *tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})

	c.JSON(http.StatusOK, web.SuccessResponseMessage{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "login user success",
	})
}

func (u *userControllerImpl) Update(c *gin.Context) {
	var userUpdateRequest web.UserUpdateRequest
	userUpdateRequest.Username = c.Param("username")
	if err := helper.ReadFromRequestBody(c, &userUpdateRequest); err != nil {
		return
	}

	err := u.Validate.Struct(userUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "validation error: " + err.Error(),
		})
		return
	}

	errResponse := u.UserService.Update(userUpdateRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponseMessage{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "update user success",
	})
}

func (u *userControllerImpl) Delete(c *gin.Context) {
	username := c.Param("username")
	errResponse := u.UserService.Delete(username)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponseMessage{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "delete user success",
	})
}

func (u *userControllerImpl) GetAll(c *gin.Context) {
	users, errResponse := u.UserService.GetAll()
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponseData{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "get all user success",
		Data:    users,
	})
}

func (u *userControllerImpl) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	user, errResponse := u.UserService.GetByUsername(username)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponseData{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "get user success",
		Data:    user,
	})
}
