package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"inventory-management-system/app"
	"inventory-management-system/controller"
	"inventory-management-system/helper"
	"inventory-management-system/model/domain"
	"inventory-management-system/repository"
	"inventory-management-system/service"
)

func main() {
	dbCredential := domain.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "inventory_test",
		Port:         5432,
		Schema:       "public",
	}

	postgres := *app.NewDB()
	connection, err := postgres.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	//err = connection.AutoMigrate(domain.Users{}, domain.Sessions{}, domain.Items{}, domain.Categories{}, domain.Activities{})
	//if err != nil {
	//	panic(err)
	//}

	validate := *validator.New()
	userRepository := repository.NewUserRepository(connection)
	reportRepository := repository.NewReportRepository(connection)
	sessionRepository := repository.NewSessionRepository(connection)
	categoryRepository := repository.NewCategoryRepository(connection)
	itemRepository := repository.NewItemRepository(connection)
	userService := service.NewUserService(userRepository, sessionRepository)
	reportService := service.NewReportService(reportRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	itemService := service.NewItemService(itemRepository, reportRepository)
	userController := controller.NewUserController(userService, &validate)
	reportController := controller.NewReportController(reportService)
	categoryController := controller.NewCategoryController(categoryService, &validate)
	itemController := controller.NewItemController(itemService, &validate)

	postgres.Reset(connection, "users")
	helper.RegisterAdmin(userRepository)

	apiServer := gin.New()
	app.UserRouter(apiServer, userController)
	app.CategoryRouter(apiServer, categoryController)
	app.ItemRouter(apiServer, itemController)
	app.ReportRouter(apiServer, reportController)
	err = apiServer.Run(":8080")
	if err != nil {
		panic(err)
	}
}
