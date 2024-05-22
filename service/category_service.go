package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
)

type CategoryService interface {
	Add(categoryAddRequest *web.CategoryAddRequest) (domain.Categories, error)
	Update(categoryUpdateRequest web.CategoryUpdateRequest) (domain.Categories, error)
	Delete(categoryID int) error
	GetAll() ([]domain.Categories, error)
	GetByID(categoryID int) (domain.Categories, error)
	CheckAvailable(name string) bool
}

type categoryServiceImpl struct {
	repository.CategoryRepository
	*validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, validate *validator.Validate) CategoryService {
	return &categoryServiceImpl{categoryRepository, validate}
}

func (c *categoryServiceImpl) Add(categoryAddRequest *web.CategoryAddRequest) (domain.Categories, error) {
	err := c.Validate.Struct(categoryAddRequest)
	if err != nil {
		return domain.Categories{}, err
	}

	result := c.CheckAvailable(categoryAddRequest.Name)
	if result {
		return domain.Categories{}, errors.New("category already exist")
	}

	category, err := c.CategoryRepository.Add(domain.Categories{
		Name: categoryAddRequest.Name,
	})
	if err != nil {
		return domain.Categories{}, err
	}

	return category, nil
}

func (c *categoryServiceImpl) Update(categoryUpdateRequest web.CategoryUpdateRequest) (domain.Categories, error) {
	err := c.Validate.Struct(categoryUpdateRequest)
	if err != nil {
		return domain.Categories{}, err
	}

	category, err := c.CategoryRepository.Update(domain.Categories{
		ID:   categoryUpdateRequest.ID,
		Name: categoryUpdateRequest.Name,
	})
	if err != nil {
		return domain.Categories{}, err
	}

	return category, nil
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
