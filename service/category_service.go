package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
)

type CategoryService interface {
	Add(categoryAddRequest *web.CategoryAddRequest) web.ErrorResponse
	Update(categoryUpdateRequest web.CategoryUpdateRequest) web.ErrorResponse
	Delete(categoryID int) web.ErrorResponse
	GetAll() ([]domain.Categories, web.ErrorResponse)
	GetByID(categoryID int) (domain.Categories, web.ErrorResponse)
	CheckAvailable(name string) bool
}

type categoryServiceImpl struct {
	repository.HandlerRepository
}

func NewCategoryService(handleRepository repository.HandlerRepository) CategoryService {
	return &categoryServiceImpl{handleRepository}
}

func (c *categoryServiceImpl) Add(categoryAddRequest *web.CategoryAddRequest) web.ErrorResponse {
	result := c.CheckAvailable(categoryAddRequest.Name)
	if result {
		return web.NewBadRequestError("category already exists")
	}

	err := c.HandlerRepository.Add(&domain.Categories{
		Name: categoryAddRequest.Name,
	})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (c *categoryServiceImpl) Update(categoryUpdateRequest web.CategoryUpdateRequest) web.ErrorResponse {
	ok := c.CheckAvailable(categoryUpdateRequest.Name)
	if ok {
		return web.NewBadRequestError("category already exists")
	}

	category := domain.Categories{}
	err := c.HandlerRepository.GetByID(categoryUpdateRequest.ID, &category)
	if err != nil {
		return web.NewBadRequestError("category id not exists")
	}

	result := c.CheckAvailable(category.Name)
	if !result {
		return web.NewBadRequestError("category name exists")
	}

	err = c.HandlerRepository.UpdateByID(categoryUpdateRequest.ID, &domain.Categories{
		ID:   categoryUpdateRequest.ID,
		Name: categoryUpdateRequest.Name,
	})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (c *categoryServiceImpl) Delete(categoryID int) web.ErrorResponse {
	category := domain.Categories{}
	err := c.HandlerRepository.GetByID(categoryID, &category)
	if err != nil {
		return web.NewBadRequestError("category id not exists")
	}

	err = c.HandlerRepository.DeleteByID(categoryID, &domain.Categories{})
	if err != nil {
		return web.NewInternalServerErrorError(err.Error())
	}

	return nil
}

func (c *categoryServiceImpl) GetAll() ([]domain.Categories, web.ErrorResponse) {
	categories := []domain.Categories{}
	err := c.HandlerRepository.GetAll(&categories)
	if err != nil {
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	if len(categories) == 0 {
		return nil, web.NewNotFoundError("category not found")
	}

	return categories, nil
}

func (c *categoryServiceImpl) GetByID(categoryID int) (domain.Categories, web.ErrorResponse) {
	category := domain.Categories{}
	err := c.HandlerRepository.GetByID(categoryID, &category)
	if err != nil {
		return domain.Categories{}, web.NewInternalServerErrorError(err.Error())
	}

	return category, nil
}

func (c *categoryServiceImpl) CheckAvailable(name string) bool {
	err := c.HandlerRepository.GetByName(name, &domain.Categories{})
	if err != nil {
		return false
	}
	return true
}
