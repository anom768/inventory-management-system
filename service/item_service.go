package service

import (
	"github.com/go-playground/validator/v10"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
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
	*validator.Validate
}

func NewItemService(itemRepository repository.ItemRepository, validate *validator.Validate) ItemService {
	return &itemServiceImpl{itemRepository, validate}
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

	return item, nil
}

func (i *itemServiceImpl) Update(itemUpdateRequest web.ItemUpdateRequest) (domain.Items, error) {
	err := i.Validate.Struct(itemUpdateRequest)
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

	return item, nil
}

func (i *itemServiceImpl) Delete(itemID int) error {
	return i.ItemRepository.Delete(itemID)
}

func (i *itemServiceImpl) GetAll() ([]domain.Items, error) {
	return i.ItemRepository.GetAll()
}

func (i *itemServiceImpl) GetByID(itemID int) (domain.Items, error) {
	return i.ItemRepository.GetByItemID(itemID)
}
