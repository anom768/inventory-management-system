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
}

type itemServiceImpl struct {
	repository.ItemRepository
	repository.ReportRepository
}

func NewItemService(itemRepository repository.ItemRepository, activityRepository repository.ReportRepository) ItemService {
	return &itemServiceImpl{itemRepository, activityRepository}
}

func (i *itemServiceImpl) Add(itemAddRequest web.ItemAddRequest) web.ErrorResponse {
	item := domain.Items{
		Name:          itemAddRequest.Name,
		CategoryID:    itemAddRequest.CategoryID,
		Price:         itemAddRequest.Price,
		Quantity:      itemAddRequest.Quantity,
		Specification: itemAddRequest.Specification,
	}
	err := i.ItemRepository.Add(item)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	err = i.ReportRepository.AddActivity(domain.Activities{
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
	itemDB, err := i.ItemRepository.GetByItemID(itemUpdateRequest.ID)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
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
	err = i.ItemRepository.Update(item)
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

	err = i.ReportRepository.AddActivity(domain.Activities{
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
	if err := i.ItemRepository.Delete(itemID); err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	err := i.ReportRepository.AddActivity(domain.Activities{
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
	items, err := i.ItemRepository.GetAll()
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return items, web.ErrorResponse{}
}

func (i *itemServiceImpl) GetByID(itemID int) (domain.Items, web.ErrorResponse) {
	item, err := i.ItemRepository.GetByItemID(itemID)
	if err != nil {
		return domain.Items{}, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return item, web.ErrorResponse{}
}
