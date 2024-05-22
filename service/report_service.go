package service

import (
	"inventory-management-system/model/domain"
	"inventory-management-system/repository"
)

type ReportService interface {
	GetAllActivity() ([]domain.Activities, error)
	ReportStock(stockItem int) ([]domain.Items, error)
}

type reportServiceImpl struct {
	repository.ReportRepository
}

func NewReportService(reportRepository repository.ReportRepository) ReportService {
	return &reportServiceImpl{reportRepository}
}

func (r *reportServiceImpl) GetAllActivity() ([]domain.Activities, error) {
	return r.ReportRepository.GetAllActivity()
}

func (r *reportServiceImpl) ReportStock(stockItem int) ([]domain.Items, error) {
	return r.ReportRepository.ReportStock(stockItem)
}
