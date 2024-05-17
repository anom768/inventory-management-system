package repository

import (
	"gorm.io/gorm"
	"inventory-management-system/model/domain"
)

type ActivityRepository interface {
	Add(activity domain.Activities) (domain.Activities, error)
	GetAll() ([]domain.Activities, error)
}

type activityRepository struct {
	*gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db}
}

func (ar *activityRepository) Add(activity domain.Activities) (domain.Activities, error) {
	if err := ar.DB.Create(&activity).Error; err != nil {
		return domain.Activities{}, err
	}

	return activity, nil
}

func (ar *activityRepository) GetAll() ([]domain.Activities, error) {
	var activities []domain.Activities
	if err := ar.Find(&activities).Error; err != nil {
		return activities, err
	}

	return activities, nil
}
