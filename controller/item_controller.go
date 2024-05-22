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

type ItemController interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
}

type itemControllerImpl struct {
	service.ItemService
	*validator.Validate
}

func NewItemController(itemService service.ItemService, validate *validator.Validate) ItemController {
	return &itemControllerImpl{itemService, validate}
}

func (i *itemControllerImpl) Add(c *gin.Context) {
	var itemAddRequest web.ItemAddRequest
	helper.ReadFromRequestBody(c, itemAddRequest)

	if err := i.Validate.Struct(&itemAddRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "validation error: " + err.Error(),
		})
		return
	}

	errResponse := i.ItemService.Add(itemAddRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusCreated, web.SuccessResponse{
		Code:    http.StatusCreated,
		Status:  "status created",
		Message: "success create item",
	})
}

func (i *itemControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	var itemUpdateRequest web.ItemUpdateRequest
	helper.ReadFromRequestBody(c, itemUpdateRequest)

	itemUpdateRequest.ID = id
	if err := i.Validate.Struct(&itemUpdateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "validation error: " + err.Error(),
		})
		return
	}

	errResponse := i.ItemService.Update(itemUpdateRequest)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status created",
		Message: "success update item",
	})
}

func (i *itemControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	errResponse := i.ItemService.Delete(id)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success delete item",
	})
}

func (i *itemControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "invalid id",
		})
		return
	}

	item, errResponse := i.ItemService.GetByID(id)
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success get item",
		Data:    item,
	})
}

func (i *itemControllerImpl) GetAll(c *gin.Context) {
	items, errResponse := i.ItemService.GetAll()
	if errResponse.Code != 0 {
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.JSON(http.StatusOK, web.SuccessResponse{
		Code:    http.StatusOK,
		Status:  "status ok",
		Message: "success get all item",
		Data:    items,
	})
}
