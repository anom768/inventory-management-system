package app

import (
	"github.com/gin-gonic/gin"
	"inventory-management-system/controller"
	"inventory-management-system/middleware"
)

func UserRouter(apiServer *gin.Engine, userController controller.UserController) *gin.Engine {
	user := apiServer.Group("/api/v1")
	user.POST("/login", userController.Login)

	user.Use(middleware.Auth())
	user.POST("/users", userController.Register)
	user.GET("/users", userController.GetAll)
	user.GET("/users/:username", userController.GetByUsername)
	user.PUT("/users/:username", userController.Update)
	user.DELETE("/users/:username", userController.Delete)

	return apiServer
}

func CategoryRouter(apiServer *gin.Engine, categoryController controller.CategoryController) *gin.Engine {
	category := apiServer.Group("/api/v1")
	category.Use(middleware.Auth())
	category.GET("/category", categoryController.GetAll)
	category.PUT("/category/:categoryID", categoryController.Update)
	category.DELETE("/category/:categoryID", categoryController.Delete)
	category.POST("/category", categoryController.Add)
	category.GET("/category/:categoryID", categoryController.GetByID)

	return apiServer
}
