package service

import (
	"github.com/go-playground/validator/v10"
	"inventory-management-system/model/domain"
	"inventory-management-system/repository"
)

type ReportService interface {
	GetAllActivity() ([]domain.Activities, error)
	ReportStock(stockItem int) ([]domain.Items, error)
}

type reportServiceImpl struct {
	repository.ReportRepository
	*validator.Validate
}

func NewReportService(reportRepository repository.ReportRepository, validate *validator.Validate) ReportService {
	return &reportServiceImpl{reportRepository, validate}
}

func (r *reportServiceImpl) GetAllActivity() ([]domain.Activities, error) {
	return r.ReportRepository.GetAllActivity()
}

func (r *reportServiceImpl) ReportStock(stockItem int) ([]domain.Items, error) {
	return r.ReportRepository.ReportStock(stockItem)
}
