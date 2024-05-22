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
	repository.ReportRepository
}

func NewReportService(reportRepository repository.ReportRepository) ReportService {
	return &reportServiceImpl{reportRepository}
}

func (r *reportServiceImpl) GetAllActivity() ([]domain.Activities, web.ErrorResponse) {
	activities, err := r.ReportRepository.GetAllActivity()
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return activities, web.ErrorResponse{}
}

func (r *reportServiceImpl) ReportStock(stockItem int) ([]domain.Items, web.ErrorResponse) {
	item, err := r.ReportRepository.ReportStock(stockItem)
	if err != nil {
		return nil, web.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "status internal server error",
			Message: err.Error(),
		}
	}

	return item, web.ErrorResponse{}
}
