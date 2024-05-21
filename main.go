package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/app"
	"inventory-management-system/controller"
	"inventory-management-system/model/domain"
	"inventory-management-system/repository"
	"inventory-management-system/service"
)

func main() {
	var postgres app.Postgres
	var apiServer *gin.Engine
	var validate validator.Validate
	var userRepository repository.UserRepository
	var activityRepository repository.ActivityRepository
	var sessionRepository repository.SessionRepository
	var categoryRepository repository.CategoryRepository
	var itemRepository repository.ItemRepository
	var userService service.UserService
	var activityService service.ActivityService
	var categoryService service.CategoryService
	var itemService service.ItemService
	var userController controller.UserController
	var activityController controller.ActivityController
	var categoryController controller.CategoryController
	var itemController controller.ItemController

	dbCredential := domain.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "inventory_test",
		Port:         5432,
		Schema:       "public",
	}

	postgres = *app.NewDB()
	connection, err := postgres.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	//err = connection.AutoMigrate(domain.Users{}, domain.Sessions{}, domain.Items{}, domain.Categories{}, domain.Activities{})
	//if err != nil {
	//	panic(err)
	//}

	validate = *validator.New()
	userRepository = repository.NewUserRepository(connection)
	activityRepository = repository.NewActivityRepository(connection)
	sessionRepository = repository.NewSessionRepository(connection)
	categoryRepository = repository.NewCategoryRepository(connection)
	itemRepository = repository.NewItemRepository(connection)
	userService = service.NewUserService(userRepository, sessionRepository, &validate)
	activityService = service.NewActivityService(activityRepository, &validate)
	categoryService = service.NewCategoryService(categoryRepository, &validate)
	itemService = service.NewItemService(itemRepository, activityRepository, &validate)
	userController = controller.NewUserController(userService, &validate)
	activityController = controller.NewActivityController(activityService)
	categoryController = controller.NewCategoryController(categoryService, &validate)
	itemController = controller.NewItemController(itemService, &validate)

	postgres.Reset(connection, "users")
	//postgres.Reset(connection, "sessions")
	registerAdmin(userRepository)

	apiServer = gin.New()
	app.UserRouter(apiServer, userController)
	app.CategoryRouter(apiServer, categoryController)
	app.ItemRouter(apiServer, itemController)
	app.ActivityRouter(apiServer, activityController)
	err = apiServer.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func registerAdmin(userRepository repository.UserRepository) {
	password, err := hashPassword("admin123")
	if err != nil {
		panic(err)
	}
	admin := domain.Users{
		Username: "admin",
		FullName: "Administrator",
		Password: password,
		Role:     "admin",
	}
	_, err = userRepository.Add(admin)
	if err != nil {
		panic(err)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
