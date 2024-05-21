package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	if err := c.ShouldBindJSON(&categoryAddRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	err := cc.Validate.Struct(categoryAddRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("validation error"))
		return
	}

	_, err = cc.CategoryService.Add(&categoryAddRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, web.NewCreated("add category successfully"))
}

func (cc *categoryControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	var categoryUpdateRequest web.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&categoryUpdateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	categoryUpdateRequest.ID = id
	err = cc.Validate.Struct(categoryUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("validation error"))
		return
	}

	_, err = cc.CategoryService.Update(categoryUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("update category successfully"))
}

func (cc *categoryControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	err = cc.CategoryService.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("delete category successfully"))
}

func (cc *categoryControllerImpl) GetAll(c *gin.Context) {
	categories, err := cc.CategoryService.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	if len(categories) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, web.NewNotFoundResponse("record not found"))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(categories))
}

func (cc *categoryControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	category, err := cc.CategoryService.GetByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(category))
}
