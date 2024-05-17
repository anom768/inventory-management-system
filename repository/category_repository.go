package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type CategoryRepository interface {
	Add(category domain.Categories) (domain.Categories, error)
	Update(category domain.Categories) (domain.Categories, error)
	Delete(categoryID int) error
	GetAll() ([]domain.Categories, error)
	GetByID(categoryID int) (domain.Categories, error)
}

type categoryRepositoryImpl struct {
	*gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db}
}

func (c *categoryRepositoryImpl) Add(category domain.Categories) (domain.Categories, error) {
	if err := c.DB.Create(&category).Error; err != nil {
		return domain.Categories{}, err
	}

	return category, nil
}

func (c *categoryRepositoryImpl) Update(category domain.Categories) (domain.Categories, error) {
	newCategory := domain.Categories{}
	if err := c.DB.First(&newCategory, "id = ?", category.ID).Error; err != nil {
		return domain.Categories{}, err
	}

	newCategory.Name = category.Name
	newCategory.Specification = category.Specification
	if err := c.DB.Save(&category).Error; err != nil {
		return domain.Categories{}, err
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

func (c *categoryRepositoryImpl) GetAll() ([]domain.Categories, error) {
	var categories []domain.Categories
	if err := c.DB.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (c *categoryRepositoryImpl) GetByID(categoryID int) (domain.Categories, error) {
	category := domain.Categories{}
	if err := c.DB.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return domain.Categories{}, err
	}

	return category, nil
}
