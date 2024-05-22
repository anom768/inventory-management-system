package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type ReportRepository interface {
	AddActivity(activity domain.Activities) error
	GetAllActivity() ([]domain.Activities, error)
	ReportStock(itemStock int) ([]domain.Items, error)
}

type reportRepositoryImpl struct {
	*gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepositoryImpl{db}
}

func (r *reportRepositoryImpl) AddActivity(activity domain.Activities) error {
	return r.DB.Create(&activity).Error
}

func (r *reportRepositoryImpl) GetAllActivity() ([]domain.Activities, error) {
	var activities []domain.Activities
	if err := r.Find(&activities).Error; err != nil {
		return activities, err
	}

	return activities, nil
}

func (r *reportRepositoryImpl) ReportStock(itemStock int) ([]domain.Items, error) {
	var items []domain.Items
	if err := r.Where("quantity <= ?", itemStock).Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}
