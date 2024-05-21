package controller

import (
	"github.com/gin-gonic/gin"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/service"
	"net/http"
	"strconv"
)

type ReportController interface {
	GetAllActivity(c *gin.Context)
	ReportStock(c *gin.Context)
}

type reportControllerImpl struct {
	service.ReportService
}

func NewReportController(reportService service.ReportService) ReportController {
	return &reportControllerImpl{reportService}
}

func (r *reportControllerImpl) GetAllActivity(c *gin.Context) {
	activities, err := r.ReportService.GetAllActivity()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	if len(activities) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, web.NewNotFoundResponse("activity not found"))
		return
	}

	c.JSON(http.StatusOK, web.NewResponseModel(activities))
}

func (r *reportControllerImpl) ReportStock(c *gin.Context) {
	totalStock, err := strconv.Atoi(c.Param("itemStock"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("invalid item stock"))
		return
	}

	items, err := r.ReportService.ReportStock(totalStock)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, web.NewInternalServerError(err.Error()))
		return
	}

	if len(items) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, web.NewNotFoundResponse("item not found"))
		return
	}

	reportStock := domain.ReportStock{}
	for _, item := range items {
		reportStock.TotalItems += 1
		reportStock.TotalQuantity += item.Quantity
		reportStock.TotalInventoryValue += item.Price
	}

	reportStock.Items = items
	c.JSON(http.StatusOK, web.NewResponseModel(reportStock))
}
