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
	if err := helper.ReadFromRequestBody(c, &itemAddRequest); err != nil {
		return
	}

	if err := i.Validate.Struct(&itemAddRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("validation error: "+err.Error()))
		return
	}

	username, _ := c.Get("username")
	errResponse := i.ItemService.Add(itemAddRequest, username.(string))
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusCreated, web.NewStatusCreated("success create item"))
}

func (i *itemControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	var itemUpdateRequest web.ItemUpdateRequest
	if err := helper.ReadFromRequestBody(c, &itemUpdateRequest); err != nil {
		return
	}

	itemUpdateRequest.ID = id
	if err := i.Validate.Struct(&itemUpdateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("validation error: "+err.Error()))
		return
	}

	username, _ := c.Get("username")
	errResponse := i.ItemService.Update(itemUpdateRequest, username.(string))
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKMessage("success update item"))
}

func (i *itemControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	username, _ := c.Get("username")
	errResponse := i.ItemService.Delete(id, username.(string))
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKMessage("success delete item"))
}

func (i *itemControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid id"))
		return
	}

	item, errResponse := i.ItemService.GetByID(id)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKData("success get item", item))
}

func (i *itemControllerImpl) GetAll(c *gin.Context) {
	items, errResponse := i.ItemService.GetAll()
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKData("success get all item", items))
}
