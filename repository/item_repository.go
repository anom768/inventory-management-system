package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model"
)

type ItemRepository interface {
	Add(item model.Items) (model.Items, error)
	Update(item model.Items) (model.Items, error)
	Delete(itemID int) error
	GetAll() ([]model.Items, error)
	GetByItemID(itemID int) (model.Items, error)
	GetByCategoryID(categoryID int) ([]model.Items, error)
}

type itemRepositoryImpl struct {
	*gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepositoryImpl{db}
}

func (i *itemRepositoryImpl) Add(item model.Items) (model.Items, error) {
	if err := i.DB.Create(&item).Error; err != nil {
		return model.Items{}, err
	}

	return item, nil
}

func (i *itemRepositoryImpl) Update(item model.Items) (model.Items, error) {
	newItem := model.Items{}
	if err := i.DB.First(&newItem, "id = ?", item.ID).Error; err != nil {
		return model.Items{}, err
	}

	newItem.CategoryID = item.CategoryID
	newItem.Price = item.Price
	newItem.Quantity = item.Quantity
	newItem.Description = item.Description
	if err := i.DB.Save(&item).Error; err != nil {
		return model.Items{}, err
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

func (i *itemRepositoryImpl) GetAll() ([]model.Items, error) {
	var items []model.Items
	if err := i.DB.Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}

func (i *itemRepositoryImpl) GetByItemID(itemID int) (model.Items, error) {
	item := model.Items{}
	if err := i.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		return model.Items{}, err
	}

	return item, nil
}

func (i *itemRepositoryImpl) GetByCategoryID(categoryID int) ([]model.Items, error) {
	var items []model.Items
	if err := i.DB.Where("category_id = ?", categoryID).Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}
