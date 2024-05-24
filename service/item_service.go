package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"net/http"
	"time"
)

type ItemService interface {
	Add(itemAddRequest web.ItemAddRequest) web.ErrorResponse
	Update(itemUpdateRequest web.ItemUpdateRequest) web.ErrorResponse
	Delete(itemID int) web.ErrorResponse
	GetAll() ([]domain.Items, web.ErrorResponse)
	GetByID(itemID int) (domain.Items, web.ErrorResponse)
	CheckAvailable(name string) bool
}

type itemServiceImpl struct {
	repository.HandlerRepository
}

func NewItemService(handlerRepository repository.HandlerRepository) ItemService {
	return &itemServiceImpl{handlerRepository}
}

func (i *itemServiceImpl) Add(itemAddRequest web.ItemAddRequest) web.ErrorResponse {
	category := domain.Categories{}
	err := i.HandlerRepository.GetByID(itemAddRequest.CategoryID, &category)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "category id not found",
		}
	}

	if ok := i.CheckAvailable(itemAddRequest.Name); ok {
		return web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "item name is already in use",
		}
	}

	item := domain.Items{
		Name:          itemAddRequest.Name,
		CategoryID:    itemAddRequest.CategoryID,
		Price:         itemAddRequest.Price,
		Quantity:      itemAddRequest.Quantity,
		Specification: itemAddRequest.Specification,
	}
	err = i.HandlerRepository.Add(&item)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	err = i.HandlerRepository.Add(&domain.Activities{
		ItemID:         item.ID,
		Action:         "POST",
		QuantityChange: item.Quantity,
		Timestamp:      time.Now(),
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (i *itemServiceImpl) Update(itemUpdateRequest web.ItemUpdateRequest) web.ErrorResponse {
	category := domain.Categories{}
	err := i.HandlerRepository.GetByID(itemUpdateRequest.CategoryID, &category)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "category id not found",
		}
	}

	itemDB := domain.Items{}
	err = i.HandlerRepository.GetByID(itemUpdateRequest.ID, &itemDB)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "item id not found",
		}
	}

	if itemDB.ID == itemUpdateRequest.ID && itemDB.Name == itemUpdateRequest.Name {

	} else {
		if ok := i.CheckAvailable(itemUpdateRequest.Name); ok {
			return web.ErrorResponse{
				Code:    http.StatusBadRequest,
				Status:  "status bad request",
				Message: "item name is already in use",
			}
		}
	}

	item := domain.Items{
		ID:            itemUpdateRequest.ID,
		Name:          itemUpdateRequest.Name,
		CategoryID:    itemUpdateRequest.CategoryID,
		Price:         itemUpdateRequest.Price,
		Quantity:      itemUpdateRequest.Quantity,
		Specification: itemUpdateRequest.Specification,
	}
	err = i.HandlerRepository.UpdateByID(itemUpdateRequest.ID, &item)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	var quantityChange int
	if itemDB.Quantity == itemUpdateRequest.Quantity {
		quantityChange = 0
	} else {
		quantityChange = itemUpdateRequest.Quantity - itemDB.Quantity
	}

	err = i.HandlerRepository.Add(&domain.Activities{
		ItemID:         item.ID,
		Action:         "UPDATE",
		QuantityChange: quantityChange,
		Timestamp:      time.Now(),
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (i *itemServiceImpl) Delete(itemID int) web.ErrorResponse {
	item := domain.Items{}
	err := i.HandlerRepository.GetByID(itemID, &item)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "item id not found",
		}
	}

	if err := i.HandlerRepository.DeleteByID(itemID, domain.Items{}); err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	err = i.HandlerRepository.Add(&domain.Activities{
		ItemID:         itemID,
		Action:         "DELETE",
		QuantityChange: 0,
		Timestamp:      time.Now(),
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (i *itemServiceImpl) GetAll() ([]domain.Items, web.ErrorResponse) {
	items := []domain.Items{}
	err := i.HandlerRepository.GetAll(&items)
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	if len(items) == 0 {
		return items, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "item not found",
		}
	}

	return items, web.ErrorResponse{}
}

func (i *itemServiceImpl) GetByID(itemID int) (domain.Items, web.ErrorResponse) {
	item := domain.Items{}
	err := i.HandlerRepository.GetByID(itemID, &item)
	if err != nil {
		return domain.Items{}, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return item, web.ErrorResponse{}
}

func (i *itemServiceImpl) CheckAvailable(name string) bool {
	err := i.HandlerRepository.GetByName(name, &domain.Items{})
	if err != nil {
		return false
	}
	return true
}
