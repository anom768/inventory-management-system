package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"time"
)

type ItemService interface {
	Add(itemAddRequest web.ItemAddRequest) error
	Update(itemUpdateRequest web.ItemUpdateRequest) error
	Delete(itemID int) error
	GetAll() ([]domain.Items, error)
	GetByID(itemID int) (domain.Items, error)
}

type itemServiceImpl struct {
	repository.ItemRepository
	repository.ReportRepository
}

func NewItemService(itemRepository repository.ItemRepository, activityRepository repository.ReportRepository) ItemService {
	return &itemServiceImpl{itemRepository, activityRepository}
}

func (i *itemServiceImpl) Add(itemAddRequest web.ItemAddRequest) error {
	item := domain.Items{
		Name:          itemAddRequest.Name,
		CategoryID:    itemAddRequest.CategoryID,
		Price:         itemAddRequest.Price,
		Quantity:      itemAddRequest.Quantity,
		Specification: itemAddRequest.Specification,
	}
	err := i.ItemRepository.Add(item)
	if err != nil {
		return err
	}

	err = i.ReportRepository.AddActivity(domain.Activities{
		ItemID:         item.ID,
		Action:         "POST",
		QuantityChange: item.Quantity,
		Timestamp:      time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (i *itemServiceImpl) Update(itemUpdateRequest web.ItemUpdateRequest) error {
	itemDB, err := i.ItemRepository.GetByItemID(itemUpdateRequest.ID)
	if err != nil {
		return err
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
		return err
	}

	var quantityChange int
	if itemDB.Quantity == itemUpdateRequest.Quantity {
		quantityChange = 0
	} else {
		quantityChange = itemUpdateRequest.Quantity - itemDB.Quantity
	}

	return i.ReportRepository.AddActivity(domain.Activities{
		ItemID:         item.ID,
		Action:         "UPDATE",
		QuantityChange: quantityChange,
		Timestamp:      time.Now(),
	})
}

func (i *itemServiceImpl) Delete(itemID int) error {
	if err := i.ItemRepository.Delete(itemID); err != nil {
		return err
	}

	return i.ReportRepository.AddActivity(domain.Activities{
		ItemID:         itemID,
		Action:         "DELETE",
		QuantityChange: 0,
		Timestamp:      time.Now(),
	})
}

func (i *itemServiceImpl) GetAll() ([]domain.Items, error) {
	return i.ItemRepository.GetAll()
}

func (i *itemServiceImpl) GetByID(itemID int) (domain.Items, error) {
	return i.ItemRepository.GetByItemID(itemID)
}
