package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"net/http"
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
	repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{categoryRepository}
}

func (c *categoryServiceImpl) Add(categoryAddRequest *web.CategoryAddRequest) web.ErrorResponse {
	result := c.CheckAvailable(categoryAddRequest.Name)
	if result {
		return web.ErrorResponse{
			Code:    http.StatusBadRequest,
			Status:  "status bad request",
			Message: "category already exists",
		}
	}

	err := c.CategoryRepository.Add(domain.Categories{
		Name: categoryAddRequest.Name,
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (c *categoryServiceImpl) Update(categoryUpdateRequest web.CategoryUpdateRequest) web.ErrorResponse {
	err := c.CategoryRepository.Update(domain.Categories{
		ID:   categoryUpdateRequest.ID,
		Name: categoryUpdateRequest.Name,
	})
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (c *categoryServiceImpl) Delete(categoryID int) web.ErrorResponse {
	err := c.CategoryRepository.Delete(categoryID)
	if err != nil {
		return web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return web.ErrorResponse{}
}

func (c *categoryServiceImpl) GetAll() ([]domain.Categories, web.ErrorResponse) {
	categories, err := c.CategoryRepository.GetAll()
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return categories, web.ErrorResponse{}
}

func (c *categoryServiceImpl) GetByID(categoryID int) (domain.Categories, web.ErrorResponse) {
	category, err := c.CategoryRepository.GetByID(categoryID)
	if err != nil {
		return domain.Categories{}, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "internal server error",
			Message: err.Error(),
		}
	}

	return category, web.ErrorResponse{}
}

func (c *categoryServiceImpl) CheckAvailable(name string) bool {
	_, err := c.CategoryRepository.GetByName(name)
	if err != nil {
		return false
	}
	return true
}
