package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
)

type ReportService interface {
	GetAllActivity() ([]domain.Activities, web.ErrorResponse)
	ReportStock(stockItem int) ([]domain.Items, web.ErrorResponse)
}

type reportServiceImpl struct {
	repository.HandlerRepository
}

func NewReportService(handleRepository repository.HandlerRepository) ReportService {
	return &reportServiceImpl{handleRepository}
}

func (r *reportServiceImpl) GetAllActivity() ([]domain.Activities, web.ErrorResponse) {
	activities := []domain.Activities{}
	err := r.HandlerRepository.GetAll(&activities)
	if err != nil {
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	if len(activities) == 0 {
		return nil, web.NewNotFoundError("report not found")
	}

	return activities, nil
}

func (r *reportServiceImpl) ReportStock(stockItem int) ([]domain.Items, web.ErrorResponse) {
	items, err := r.HandlerRepository.ReportStock(stockItem)
	if err != nil {
		return nil, web.NewInternalServerErrorError(err.Error())
	}

	if len(items) == 0 {
		return nil, web.NewNotFoundError("report not found")
	}

	return items, nil
}
