package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	err := c.ShouldBind(&userRegisterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	err = u.Validate.Struct(userRegisterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	_, err = u.UserService.Register(&userRegisterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, web.NewCreated("registration successful"))
}

func (u *userControllerImpl) Login(c *gin.Context) {
	var userLoginRequest web.UserLoginRequest
	err := c.ShouldBindJSON(&userLoginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	err = u.Validate.Struct(userLoginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	tokenString, err := u.UserService.Login(&userLoginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   *tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})

	c.JSON(http.StatusOK, web.NewOk("login successful"))
}

func (u *userControllerImpl) Update(c *gin.Context) {
	var userUpdateRequest web.UserUpdateRequest
	err := c.ShouldBind(&userUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	err = u.Validate.Struct(userUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	_, err = u.UserService.Update(userUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("update successful"))
}

func (u *userControllerImpl) Delete(c *gin.Context) {
	username := c.Param("username")
	err := u.UserService.Delete(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("delete successful"))
}

func (u *userControllerImpl) GetAll(c *gin.Context) {
	users, err := u.UserService.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(users))
}

func (u *userControllerImpl) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := u.UserService.GetByUsername(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(user))
}
