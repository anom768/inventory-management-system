package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model"
)

type CategoryRepository interface {
	Add(category model.Categories) (model.Categories, error)
	Update(category model.Categories) (model.Categories, error)
	Delete(categoryID int) error
	GetAll() ([]model.Categories, error)
	GetByID(categoryID int) (model.Categories, error)
}

type categoryRepositoryImpl struct {
	*gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db}
}

func (c *categoryRepositoryImpl) Add(category model.Categories) (model.Categories, error) {
	if err := c.DB.Create(&category).Error; err != nil {
		return model.Categories{}, err
	}

	return category, nil
}

func (c *categoryRepositoryImpl) Update(category model.Categories) (model.Categories, error) {
	newCategory := model.Categories{}
	if err := c.DB.First(&newCategory, "username = ?", category.ID).Error; err != nil {
		return model.Categories{}, err
	}

	newCategory.Name = category.Name
	newCategory.Specification = category.Specification
	if err := c.DB.Save(&category).Error; err != nil {
		return model.Categories{}, err
	}

	return newCategory, nil
}

func (c *categoryRepositoryImpl) Delete(categoryID int) error {
	category, err := c.GetByID(categoryID)
	if err != nil {
		return err
	}

	if err := c.DB.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}

func (c *categoryRepositoryImpl) GetAll() ([]model.Categories, error) {
	var categories []model.Categories
	if err := c.DB.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (c *categoryRepositoryImpl) GetByID(categoryID int) (model.Categories, error) {
	category := model.Categories{}
	if err := c.DB.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return model.Categories{}, err
	}

	return category, nil
}
