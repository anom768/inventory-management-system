package app

import (
	"github.com/gin-gonic/gin"
	"inventory-management-system/controller"
	"inventory-management-system/middleware"
)

func UserRouter(gin *gin.Engine, userController controller.UserController) *gin.Engine {
	user := gin.Group("/api/v1")
	user.POST("/login", userController.Login)

	user.Use(middleware.Auth())
	user.POST("/users", userController.Register)
	user.GET("/users", userController.GetAll)
	user.GET("/users/:username", userController.GetByUsername)
	user.PUT("/users/:username", userController.Update)
	user.DELETE("/users/:username", userController.Delete)

	return gin
}
