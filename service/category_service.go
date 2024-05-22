package service

import (
	"errors"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
)

type CategoryService interface {
	Add(categoryAddRequest *web.CategoryAddRequest) error
	Update(categoryUpdateRequest web.CategoryUpdateRequest) error
	Delete(categoryID int) error
	GetAll() ([]domain.Categories, error)
	GetByID(categoryID int) (domain.Categories, error)
	CheckAvailable(name string) bool
}

type categoryServiceImpl struct {
	repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{categoryRepository}
}

func (c *categoryServiceImpl) Add(categoryAddRequest *web.CategoryAddRequest) error {
	result := c.CheckAvailable(categoryAddRequest.Name)
	if result {
		return errors.New("category already exist")
	}

	return c.CategoryRepository.Add(domain.Categories{
		Name: categoryAddRequest.Name,
	})
}

func (c *categoryServiceImpl) Update(categoryUpdateRequest web.CategoryUpdateRequest) error {
	return c.CategoryRepository.Update(domain.Categories{
		ID:   categoryUpdateRequest.ID,
		Name: categoryUpdateRequest.Name,
	})
}

func (c *categoryServiceImpl) Delete(categoryID int) error {
	return c.CategoryRepository.Delete(categoryID)
}

func (c *categoryServiceImpl) GetAll() ([]domain.Categories, error) {
	return c.CategoryRepository.GetAll()
}

func (c *categoryServiceImpl) GetByID(categoryID int) (domain.Categories, error) {
	return c.CategoryRepository.GetByID(categoryID)
}

func (c *categoryServiceImpl) CheckAvailable(name string) bool {
	_, err := c.CategoryRepository.GetByName(name)
	if err != nil {
		return false
	}
	return true
}
