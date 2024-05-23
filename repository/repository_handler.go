package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type HandlerRepository interface {
	Add(v any) error
	UpdateByID(id int, new any) error
	UpdateByUsername(username string, new any) error
	DeleteByID(id int, v any) error
	DeleteByUsername(username string, v any) error
	GetAll(v any) error
	GetByID(id int, v any) error
	GetByUsername(username string, v any) error
	ReportStock(itemStock int) ([]domain.Items, error)
}

type handlerRepositoryImpl struct {
	*gorm.DB
}

func NewHandlerRepository(db *gorm.DB) HandlerRepository {
	return &handlerRepositoryImpl{db}
}

func (h *handlerRepositoryImpl) Add(v any) error {
	return h.DB.Create(v).Error
}

func (h *handlerRepositoryImpl) UpdateByID(id int, new any) error {
	return h.DB.Where("id = ?", id).Updates(new).Error
}

func (h *handlerRepositoryImpl) UpdateByUsername(username string, new any) error {
	return h.DB.Where("username = ?", username).Updates(new).Error
}

func (h *handlerRepositoryImpl) DeleteByID(id int, v any) error {
	return h.DB.Where("id = ?", id).Delete(v).Error
}

func (h *handlerRepositoryImpl) DeleteByUsername(username string, v any) error {
	return h.DB.Where("username = ?", username).Delete(v).Error
}

func (h *handlerRepositoryImpl) GetAll(v any) error {
	return h.DB.Find(v).Error
}

func (h *handlerRepositoryImpl) GetByID(id int, v any) error {
	return h.DB.Where("id = ?", id).First(&v).Error
}

func (h *handlerRepositoryImpl) GetByUsername(username string, v any) error {
	return h.DB.Where("username = ?", username).First(v).Error
}

func (h *handlerRepositoryImpl) ReportStock(itemStock int) ([]domain.Items, error) {
	var items []domain.Items
	if err := h.Where("quantity <= ?", itemStock).Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}
