package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"net/http"
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
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	if len(activities) == 0 {
		return nil, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "report service is empty",
		}
	}

	return activities, web.ErrorResponse{}
}

func (r *reportServiceImpl) ReportStock(stockItem int) ([]domain.Items, web.ErrorResponse) {
	items, err := r.HandlerRepository.ReportStock(stockItem)
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	if len(items) == 0 {
		return nil, web.ErrorResponse{
			Code:    http.StatusNotFound,
			Status:  "status not found",
			Message: "report service is empty",
		}
	}

	return items, web.ErrorResponse{}
}
