package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	if err := c.ShouldBindJSON(&itemAddRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	if err := i.Validate.Struct(&itemAddRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("validation error"))
		return
	}

	err := i.ItemService.Add(itemAddRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, web.NewCreated("create item successful"))
}

func (i *itemControllerImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	var itemUpdateRequest web.ItemUpdateRequest
	if err := c.ShouldBindJSON(&itemUpdateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid body request"))
		return
	}

	itemUpdateRequest.ID = id
	if err := i.Validate.Struct(&itemUpdateRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("validation error"))
		return
	}

	err = i.ItemService.Update(itemUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("update item successful"))
}

func (i *itemControllerImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	err = i.ItemService.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewOk("delete item successful"))
}

func (i *itemControllerImpl) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid id"))
		return
	}

	item, err := i.ItemService.GetByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(item))
}

func (i *itemControllerImpl) GetAll(c *gin.Context) {
	items, err := i.ItemService.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	if len(items) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, web.NewNotFoundResponse("item not found"))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(items))
}
