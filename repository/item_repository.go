package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type ItemRepository interface {
	Add(item domain.Items) (domain.Items, error)
	Update(item domain.Items) (domain.Items, error)
	Delete(itemID int) error
	GetAll() ([]domain.Items, error)
	GetByItemID(itemID int) (domain.Items, error)
	GetByCategoryID(categoryID int) ([]domain.Items, error)
}

type itemRepositoryImpl struct {
	*gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepositoryImpl{db}
}

func (i *itemRepositoryImpl) Add(item domain.Items) (domain.Items, error) {
	if err := i.DB.Create(&item).Error; err != nil {
		return domain.Items{}, err
	}

	return item, nil
}

func (i *itemRepositoryImpl) Update(item domain.Items) (domain.Items, error) {
	newItem := domain.Items{}
	if err := i.DB.First(&newItem, "id = ?", item.ID).Error; err != nil {
		return domain.Items{}, err
	}

	newItem.CategoryID = item.CategoryID
	newItem.Price = item.Price
	newItem.Quantity = item.Quantity
	newItem.Specification = item.Specification
	if err := i.DB.Save(&item).Error; err != nil {
		return domain.Items{}, err
	}

	return newItem, nil
}

func (i *itemRepositoryImpl) Delete(itemID int) error {
	item, err := i.GetByItemID(itemID)
	if err != nil {
		return err
	}

	if err := i.DB.Delete(&item).Error; err != nil {
		return err
	}

	return nil
}

func (i *itemRepositoryImpl) GetAll() ([]domain.Items, error) {
	var items []domain.Items
	if err := i.DB.Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}

func (i *itemRepositoryImpl) GetByItemID(itemID int) (domain.Items, error) {
	item := domain.Items{}
	if err := i.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		return domain.Items{}, err
	}

	return item, nil
}

func (i *itemRepositoryImpl) GetByCategoryID(categoryID int) ([]domain.Items, error) {
	var items []domain.Items
	if err := i.DB.Where("category_id = ?", categoryID).Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}
