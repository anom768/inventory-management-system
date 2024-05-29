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
	activities, errResponse := r.ReportService.GetAllActivity()
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	c.JSON(http.StatusOK, web.NewStatusOKData("success get all activities", activities))
}

func (r *reportControllerImpl) ReportStock(c *gin.Context) {
	totalStock, err := strconv.Atoi(c.Param("itemStock"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("invalid item stock"))
		return
	}

	items, errResponse := r.ReportService.ReportStock(totalStock)
	if errResponse != nil {
		c.AbortWithStatusJSON(errResponse.Code(), errResponse)
		return
	}

	reportStock := domain.ReportStock{}
	for _, item := range items {
		reportStock.TotalItems += 1
		reportStock.TotalQuantity += item.Quantity
		reportStock.TotalInventoryValue += item.Price
	}

	reportStock.Items = items
	c.JSON(http.StatusOK, web.NewStatusOKData("success get all report stock", reportStock))
}
