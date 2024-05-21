package service

import (
	"github.com/go-playground/validator/v10"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"time"
)

type ItemService interface {
	Add(itemAddRequest web.ItemAddRequest) (domain.Items, error)
	Update(itemUpdateRequest web.ItemUpdateRequest) (domain.Items, error)
	Delete(itemID int) error
	GetAll() ([]domain.Items, error)
	GetByID(itemID int) (domain.Items, error)
}

type itemServiceImpl struct {
	repository.ItemRepository
	repository.ActivityRepository
	*validator.Validate
}

func NewItemService(itemRepository repository.ItemRepository, activityRepository repository.ActivityRepository, validate *validator.Validate) ItemService {
	return &itemServiceImpl{itemRepository, activityRepository, validate}
}

func (i *itemServiceImpl) Add(itemAddRequest web.ItemAddRequest) (domain.Items, error) {
	err := i.Validate.Struct(itemAddRequest)
	if err != nil {
		return domain.Items{}, err
	}

	item, err := i.ItemRepository.Add(domain.Items{
		Name:        itemAddRequest.Name,
		CategoryID:  itemAddRequest.CategoryID,
		Price:       itemAddRequest.Price,
		Quantity:    itemAddRequest.Quantity,
		Description: itemAddRequest.Description,
	})
	if err != nil {
		return domain.Items{}, err
	}

	i.ActivityRepository.Add(domain.Activities{
		ItemID:         item.ID,
		Action:         "POST",
		QuantityChange: item.Quantity,
		Timestamp:      time.Now(),
	})

	return item, nil
}

func (i *itemServiceImpl) Update(itemUpdateRequest web.ItemUpdateRequest) (domain.Items, error) {
	err := i.Validate.Struct(itemUpdateRequest)
	if err != nil {
		return domain.Items{}, err
	}

	itemDB, err := i.ItemRepository.GetByItemID(itemUpdateRequest.ID)
	if err != nil {
		return domain.Items{}, err
	}

	item, err := i.ItemRepository.Update(domain.Items{
		ID:          itemUpdateRequest.ID,
		Name:        itemUpdateRequest.Name,
		CategoryID:  itemUpdateRequest.CategoryID,
		Price:       itemUpdateRequest.Price,
		Quantity:    itemUpdateRequest.Quantity,
		Description: itemUpdateRequest.Description,
	})
	if err != nil {
		return domain.Items{}, err
	}

	var quantityChange int
	if itemDB.Quantity == itemUpdateRequest.Quantity {
		quantityChange = 0
	} else {
		quantityChange = itemUpdateRequest.Quantity - itemDB.Quantity
	}
	//if itemDB.Quantity > itemUpdateRequest.Quantity {
	//
	//} else if itemDB.Quantity < itemUpdateRequest.Quantity {
	//	quantityChange = itemUpdateRequest.Quantity - itemDB.Quantity
	//} else {
	//	quantityChange = 0
	//}
	if _, err := i.ActivityRepository.Add(domain.Activities{
		ItemID:         item.ID,
		Action:         "UPDATE",
		QuantityChange: quantityChange,
		Timestamp:      time.Now(),
	}); err != nil {
		return domain.Items{}, err
	}

	return item, nil
}

func (i *itemServiceImpl) Delete(itemID int) error {
	if err := i.ItemRepository.Delete(itemID); err != nil {
		return err
	}

	if _, err := i.ActivityRepository.Add(domain.Activities{
		ItemID:         itemID,
		Action:         "DELETE",
		QuantityChange: 0,
		Timestamp:      time.Now(),
	}); err != nil {
		return nil
	}

	return nil
}

func (i *itemServiceImpl) GetAll() ([]domain.Items, error) {
	return i.ItemRepository.GetAll()
}

func (i *itemServiceImpl) GetByID(itemID int) (domain.Items, error) {
	return i.ItemRepository.GetByItemID(itemID)
}
