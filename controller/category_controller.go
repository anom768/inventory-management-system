package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"inventory-management-system/helper"
	"inventory-management-system/model/web"
	"inventory-management-system/service"
	"net/http"
	"strconv"
)

type CategoryController interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
}

type categoryControllerImpl struct {
	service.CategoryService
	*validator.Validate
}

func NewCategoryController(categoryService service.CategoryService, validate *validator.Validate) CategoryController {
	return &categoryControllerImpl{categoryService, validate}
}

func (cc *categoryControllerImpl) Add(c *gin.Context) {
	var categoryAddRequest web.CategoryAddRequest
	if err := helper.ReadFromRequestBody(c, &categoryAddRequest); err != nil {
		return
	}

	err := cc.Validate.Struct(categoryAddRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("validation error: "+err.Error()))
		return
	}

	errResponse := cc.CategoryService.Add(&categoryAddRequest)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusCreated, web.NewStatusCreated("success add category"))
}

func (cc *categoryControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	var categoryUpdateRequest web.CategoryUpdateRequest
	categoryUpdateRequest.ID = id
	if err := helper.ReadFromRequestBody(c, &categoryUpdateRequest); err != nil {
		return
	}

	err = cc.Validate.Struct(categoryUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("validation error: "+err.Error()))
		return
	}

	errResponse := cc.CategoryService.Update(categoryUpdateRequest)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKMessage("success update category"))
}

func (cc *categoryControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	errResponse := cc.CategoryService.Delete(id)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKMessage("success delete category"))
}

func (cc *categoryControllerImpl) GetAll(c *gin.Context) {
	categories, errResponse := cc.CategoryService.GetAll()
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKData("success get all category", categories))
}

func (cc *categoryControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	category, errResponse := cc.CategoryService.GetByID(id)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKData("success get category", category))
}
