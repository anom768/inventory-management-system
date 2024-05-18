package service

import (
	"github.com/go-playground/validator/v10"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
)

type ActivityService interface {
	Add(activityAddRequest web.ActivityAddRequest) (domain.Activities, error)
	GetAll() ([]domain.Activities, error)
}

type activityServiceImpl struct {
	repository.ActivityRepository
	*validator.Validate
}

func NewActivityService(activityRepository repository.ActivityRepository, validate *validator.Validate) ActivityService {
	return &activityServiceImpl{activityRepository, validate}
}

func (s *activityServiceImpl) Add(activityAddRequest web.ActivityAddRequest) (domain.Activities, error) {
	err := s.Validate.Struct(activityAddRequest)
	if err != nil {
		return domain.Activities{}, err
	}

	activity, err := s.ActivityRepository.Add(domain.Activities{
		ItemID:         activityAddRequest.ItemID,
		Action:         activityAddRequest.Action,
		QuantityChange: activityAddRequest.QuantityChane,
		Timestamp:      activityAddRequest.Timestamp,
		PerformedBy:    activityAddRequest.PerformedBy,
	})
	if err != nil {
		return domain.Activities{}, err
	}

	return activity, nil
}

func (s *activityServiceImpl) GetAll() ([]domain.Activities, error) {
	return s.ActivityRepository.GetAll()
}
