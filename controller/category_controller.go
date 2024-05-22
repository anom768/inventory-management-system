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
	helper.ReadFromRequestBody(c, &categoryAddRequest)

	err := cc.Validate.Struct(categoryAddRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "validation error")
		return
	}

	errResponse := cc.CategoryService.Add(&categoryAddRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusCreated, web.SuccessResponse{
		Code:    http.StatusCreated,
		Status:  "status created",
		Message: "success add category",
	})
}

func (cc *categoryControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	var categoryUpdateRequest web.CategoryUpdateRequest
	helper.ReadFromRequestBody(c, categoryUpdateRequest)

	categoryUpdateRequest.ID = id
	err = cc.Validate.Struct(categoryUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	errResponse := cc.CategoryService.Update(categoryUpdateRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success update category",
	})
}

func (cc *categoryControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	errResponse := cc.CategoryService.Delete(id)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success delete category",
	})
}

func (cc *categoryControllerImpl) GetAll(c *gin.Context) {
	categories, errResponse := cc.CategoryService.GetAll()
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success get all category",
		Data:    categories,
	})
}

func (cc *categoryControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	category, errResponse := cc.CategoryService.GetByID(id)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success get category",
		Data:    category,
	})
}
