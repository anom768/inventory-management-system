package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model"
)

type ActivityRepository interface {
	GetAll() ([]model.Activities, error)
}

type activityRepository struct {
	*gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db}
}

func (ar *activityRepository) GetAll() ([]model.Activities, error) {
	var activities []model.Activities
	if err := ar.Find(&activities).Error; err != nil {
		return activities, err
	}

	return activities, nil
}
