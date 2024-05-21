package controller

import (
	"github.com/gin-gonic/gin"
	"inventory-management-system/model/web"
	"inventory-management-system/service"
	"net/http"
)

type ActivityController interface {
	GetAll(c *gin.Context)
}

type activityControllerImpl struct {
	service.ActivityService
}

func NewActivityController(activityService service.ActivityService) ActivityController {
	return &activityControllerImpl{activityService}
}

func (a *activityControllerImpl) GetAll(c *gin.Context) {
	activities, err := a.ActivityService.GetAll()
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
