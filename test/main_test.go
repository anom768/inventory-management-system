package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	//"github.com/go-playground/validator/v10"
	//"inventory-management-system/controller"
	"inventory-management-system/model/web"
	"net/http"
	"net/http/httptest"

	//"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/app"
	model "inventory-management-system/model/domain"
	//"inventory-management-system/model/web"
	//"inventory-management-system/service"
	//"inventory-management-system/service"
)

var _ = Describe("Digital Inventory Management API", func() {
	var apiServer *gin.Engine
	//var userRepository repository.UserRepository
	//var categoryRepository repository.CategoryRepository
	//var itemRepository repository.ItemRepository
	//var activityRepository repository.ActivityRepository
	//var sessionRepository repository.SessionRepository
	//var userService service.UserService
	//var categoryService service.CategoryService
	//var itemService service.ItemService
	//var activityService service.ActivityService
	//var userController controller.UserController

	db := app.NewDB()
	credential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "inventory_test",
		Port:         5432,
		Schema:       "public",
	}
	connection, err := db.Connect(&credential)
	Expect(err).ShouldNot(HaveOccurred())

	//validate := validator.New()
	//userRepository = repository.NewUserRepository(connection)
	//categoryRepository = repository.NewCategoryRepository(connection)
	//itemRepository = repository.NewItemRepository(connection)
	//activityRepository = repository.NewActivityRepository(connection)
	//sessionRepository = repository.NewSessionRepository(connection)
	//userService = service.NewUserService(userRepository, sessionRepository, validate)
	//categoryService = service.NewCategoryService(categoryRepository, validate)
	//itemService = service.NewItemService(itemRepository, validate)
	//activityService = service.NewActivityService(activityRepository, validate)
	//userController = controller.NewUserController(userService, validate)

	BeforeEach(func() {
		err = connection.Migrator().DropTable("users", "categories", "items", "activities", "sessions")
		Expect(err).ShouldNot(HaveOccurred())

		err := connection.AutoMigrate(&model.Users{}, &model.Categories{}, &model.Items{}, model.Activities{}, &model.Sessions{})
		Expect(err).ShouldNot(HaveOccurred())

		err = db.Reset(connection, "users")
		err = db.Reset(connection, "categories")
		err = db.Reset(connection, "items")
		err = db.Reset(connection, "activities")
		err = db.Reset(connection, "sessions")
		Expect(err).ShouldNot(HaveOccurred())
	})

	//Describe("Repository", func() {
	//
	//	/*
	//		########################################################################################
	//										USER REPOSITORY
	//		########################################################################################
	//	*/
	//	Describe("User Repository", func() {
	//		When("add new user to users table in database postgres", func() {
	//			It("should save data user to users table in database postgres", func() {
	//				_, err := userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "bangkit",
	//					Password: "rahasia",
	//					Role:     "admin",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = db.Reset(connection, "users")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get user by username is success", func() {
	//			It("should return data user", func() {
	//				user, err := userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "bangkit",
	//					Password: "rahasia",
	//					Role:     "admin",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := userRepository.GetByUsername("bangkit")
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.ID).To(Equal(uint(1)))
	//				Expect(result.FullName).To(Equal(user.FullName))
	//				Expect(result.Username).To(Equal(user.Username))
	//				Expect(result.Password).To(Equal(user.Password))
	//				Expect(result.Role).To(Equal(user.Role))
	//
	//				err = db.Reset(connection, "users")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get all data users table in database postgres", func() {
	//			It("should return all data users", func() {
	//				_, err := userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "bangkit",
	//					Password: "rahasia",
	//					Role:     "admin",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				_, err = userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "anom",
	//					Password: "rahasia",
	//					Role:     "user",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				results, err := userRepository.GetAll()
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(results).Should(HaveLen(2))
	//
	//				err = db.Reset(connection, "users")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("update user data to users table in database postgres", func() {
	//			It("should save new data user", func() {
	//				user, err := userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "bangkit",
	//					Password: "rahasia",
	//					Role:     "admin",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newUser, err := userRepository.GetByUsername(user.Username)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newUser.FullName = "Anom Sedhayu"
	//				newUser.Password = "newpassword"
	//				newUser.Role = "user"
	//				_, err = userRepository.Update(newUser)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := userRepository.GetByUsername(user.Username)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.FullName).To(Equal(newUser.FullName))
	//				Expect(result.Username).To(Equal(newUser.Username))
	//				Expect(result.Password).To(Equal(newUser.Password))
	//				Expect(result.Role).To(Equal(newUser.Role))
	//
	//				err = db.Reset(connection, "users")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("delete user data by username to users table in database postgres", func() {
	//			It("should soft delete data user by username", func() {
	//				user, err := userRepository.Add(model.Users{
	//					FullName: "Bangkit Anom",
	//					Username: "bangkit",
	//					Password: "rahasia",
	//					Role:     "admin",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = userRepository.Delete(user.Username)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := userRepository.GetByUsername(user.Username)
	//				Expect(err).Should(HaveOccurred())
	//				Expect(result).To(Equal(model.Users{}))
	//
	//				err = db.Reset(connection, "users")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//	})
	//
	//	/*
	//		########################################################################################
	//										CATEGORY REPOSITORY
	//		########################################################################################
	//	*/
	//	Describe("Category Repository", func() {
	//		When("add new category to categories table in database postgres", func() {
	//			It("should save data category to categories table in database postgres", func() {
	//				_, err := categoryRepository.Add(model.Categories{
	//					Name:          "VGA",
	//					Specification: "RTX-3060, RAM 4 GB",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = db.Reset(connection, "categories")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get category by id is success", func() {
	//			It("should return data category", func() {
	//				category, err := categoryRepository.Add(model.Categories{
	//					Name:          "VGA",
	//					Specification: "RTX-3060, RAM 4 GB",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := categoryRepository.GetByID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.ID).To(Equal(uint(1)))
	//				Expect(result.Name).To(Equal(category.Name))
	//				Expect(result.Specification).To(Equal(category.Specification))
	//
	//				err = db.Reset(connection, "categories")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get all data categories table in database postgres", func() {
	//			It("should return all data categories", func() {
	//				_, err := categoryRepository.Add(model.Categories{
	//					Name:          "VGA",
	//					Specification: "RTX-3060, RAM 4 GB",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				_, err = categoryRepository.Add(model.Categories{
	//					Name:          "Monitor",
	//					Specification: "14 in Ajhua",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				results, err := categoryRepository.GetAll()
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(results).Should(HaveLen(2))
	//
	//				err = db.Reset(connection, "categories")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("update category data to categories table in database postgres", func() {
	//			It("should save new data category", func() {
	//				_, err := categoryRepository.Add(model.Categories{
	//					Name:          "VGA",
	//					Specification: "RTX-3060, RAM 4 GB",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newCategory, err := categoryRepository.GetByID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newCategory.Name = "CPU"
	//				newCategory.Specification = "4 core"
	//				_, err = categoryRepository.Update(newCategory)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := categoryRepository.GetByID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.Name).To(Equal(newCategory.Name))
	//				Expect(result.Specification).To(Equal(newCategory.Specification))
	//
	//				err = db.Reset(connection, "categories")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("delete category data by id to categories table in database postgres", func() {
	//			It("should soft delete data category by id", func() {
	//				_, err := categoryRepository.Add(model.Categories{
	//					Name:          "VGA",
	//					Specification: "RTX-3060, RAM 4 GB",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = categoryRepository.Delete(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := categoryRepository.GetByID(1)
	//				Expect(err).Should(HaveOccurred())
	//				Expect(result).To(Equal(model.Categories{}))
	//
	//				err = db.Reset(connection, "categories")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//	})
	//
	//	/*
	//		########################################################################################
	//										ITEM REPOSITORY
	//		########################################################################################
	//	*/
	//	Describe("Item Repository", func() {
	//		When("add new item to items table in database postgres", func() {
	//			It("should save data item to items table in database postgres", func() {
	//				_, err := itemRepository.Add(model.Items{
	//					Name:        "VGA",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get item by item id is success", func() {
	//			It("should return data item", func() {
	//				item, err := itemRepository.Add(model.Items{
	//					Name:        "VGA",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := itemRepository.GetByItemID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.ID).To(Equal(uint(1)))
	//				Expect(result.Name).To(Equal(item.Name))
	//				Expect(result.CategoryID).To(Equal(item.CategoryID))
	//				Expect(result.Quantity).To(Equal(item.Quantity))
	//				Expect(result.Price).To(Equal(item.Price))
	//				Expect(result.Description).To(Equal(item.Description))
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get all data items table in database postgres", func() {
	//			It("should return all data items", func() {
	//				_, err := itemRepository.Add(model.Items{
	//					Name:        "VGA",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				_, err = itemRepository.Add(model.Items{
	//					Name:        "VGA2",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				results, err := itemRepository.GetAll()
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(results).Should(HaveLen(2))
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("update item data to items table in database postgres", func() {
	//			It("should save new data item", func() {
	//				_, err := itemRepository.Add(model.Items{
	//					Name:        "VGA2",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newItem, err := itemRepository.GetByItemID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				newItem.Name = "VGA"
	//				newItem.CategoryID = 2
	//				newItem.Quantity = 5
	//				newItem.Price = 4500000.00
	//				newItem.Description = "desc"
	//				_, err = itemRepository.Update(newItem)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := itemRepository.GetByItemID(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.ID).To(Equal(uint(1)))
	//				Expect(result.Name).To(Equal(newItem.Name))
	//				Expect(result.CategoryID).To(Equal(newItem.CategoryID))
	//				Expect(result.Quantity).To(Equal(newItem.Quantity))
	//				Expect(result.Price).To(Equal(newItem.Price))
	//				Expect(result.Description).To(Equal(newItem.Description))
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("delete item data by item id to items table in database postgres", func() {
	//			It("should soft delete data item by id", func() {
	//				_, err := itemRepository.Add(model.Items{
	//					Name:        "VGA2",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = itemRepository.Delete(1)
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := itemRepository.GetByItemID(1)
	//				Expect(err).Should(HaveOccurred())
	//				Expect(result).To(Equal(model.Items{}))
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get all data items table by category id in database postgres", func() {
	//			It("should return all data items by category id", func() {
	//				_, err := itemRepository.Add(model.Items{
	//					Name:        "VGA",
	//					CategoryID:  1,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				_, err = itemRepository.Add(model.Items{
	//					Name:        "VGA2",
	//					CategoryID:  2,
	//					Quantity:    10,
	//					Price:       5000000.00,
	//					Description: "",
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				results, err := itemRepository.GetByCategoryID(2)
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(results).Should(HaveLen(1))
	//
	//				err = db.Reset(connection, "items")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//	})
	//
	//	/*
	//		########################################################################################
	//										ACTIVITY REPOSITORY
	//		########################################################################################
	//	*/
	//	Describe("Item Repository", func() {
	//		When("add new activity to activities table in database postgres", func() {
	//			It("should save data activity to activities table in database postgres", func() {
	//				_, err := activityRepository.Add(model.Activities{
	//					ItemID:         1,
	//					Action:         "POST",
	//					QuantityChange: 5,
	//					PerformedBy:    1,
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = db.Reset(connection, "activities")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get all data activities table in database postgres", func() {
	//			It("should return all data activities", func() {
	//				_, err := activityRepository.Add(model.Activities{
	//					ItemID:         1,
	//					Action:         "POST",
	//					QuantityChange: 5,
	//					PerformedBy:    1,
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				_, err = activityRepository.Add(model.Activities{
	//					ItemID:         2,
	//					Action:         "POST",
	//					QuantityChange: -2,
	//					PerformedBy:    1,
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				results, err := activityRepository.GetAll()
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(results).Should(HaveLen(2))
	//
	//				err = db.Reset(connection, "activities")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//	})
	//
	//	/*
	//		########################################################################################
	//										SESSION REPOSITORY
	//		########################################################################################
	//	*/
	//	Describe("Session Repository", func() {
	//		When("add new session to sessions table in database postgres", func() {
	//			It("should save data session to sessions table in database postgres", func() {
	//				_, err := sessionRepository.Add(model.Sessions{
	//					Username:  "bangkit",
	//					Token:     "token",
	//					ExpiresAt: time.Now().Add(5 * time.Minute),
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = db.Reset(connection, "sessions")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("get session by username is success", func() {
	//			It("should return data session", func() {
	//				session, err := sessionRepository.Add(model.Sessions{
	//					Username:  "bangkit",
	//					Token:     "token",
	//					ExpiresAt: time.Now().Add(5 * time.Minute),
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := sessionRepository.GetByUsername("bangkit")
	//				Expect(err).ShouldNot(HaveOccurred())
	//				Expect(result.ID).To(Equal(uint(1)))
	//				Expect(result.Username).To(Equal(session.Username))
	//				Expect(result.Token).To(Equal(session.Token))
	//
	//				err = db.Reset(connection, "sessions")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//
	//		When("delete session data by username to sessions table in database postgres", func() {
	//			It("should soft delete data session by username", func() {
	//				_, err := sessionRepository.Add(model.Sessions{
	//					Username:  "bangkit",
	//					Token:     "token",
	//					ExpiresAt: time.Now().Add(5 * time.Minute),
	//				})
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				err = sessionRepository.Delete("bangkit")
	//				Expect(err).ShouldNot(HaveOccurred())
	//
	//				result, err := sessionRepository.GetByUsername("bangkit")
	//				Expect(err).Should(HaveOccurred())
	//				Expect(result).To(Equal(model.Sessions{}))
	//
	//				err = db.Reset(connection, "sessions")
	//				Expect(err).ShouldNot(HaveOccurred())
	//			})
	//		})
	//	})
	//})

	//====================================================================================================
	//====================================================================================================
	//====================================================================================================

	Describe("Service", func() {

		/*
			########################################################################################
											USER SERVICE
			########################################################################################
		*/
		//Describe("User Service", func() {
		//	Describe("Register", func() {
		//		When("register is successful", func() {
		//			It("should register user", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				user, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(user.FullName).To(Equal(request.FullName))
		//				Expect(user.Username).To(Equal(request.Username))
		//				Expect(true).To(Equal(checkPasswordHash(request.Password, user.Password)))
		//				Expect(user.Role).To(Equal(request.Role))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("register with blank field", func() {
		//			It("should return error", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "",
		//					Username: "",
		//					Password: "",
		//					Role:     "",
		//				}
		//				user, err := userService.Register(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(user).To(Equal(model.Users{}))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("register with duplicate username", func() {
		//			It("should return error", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				user, err := userService.Register(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(user).To(Equal(model.Users{}))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Login", func() {
		//		When("login with blank field", func() {
		//			It("should return error", func() {
		//				request := web.UserLoginRequest{
		//					Username: "",
		//					Password: "",
		//				}
		//				_, err := userService.Login(&request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("login with wrong username", func() {
		//			It("should return error", func() {
		//				request := web.UserLoginRequest{
		//					Username: "wrongusername",
		//					Password: "wrongpassword",
		//				}
		//				_, err := userService.Login(&request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("login with wrong password", func() {
		//			It("should return error", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := web.UserLoginRequest{
		//					Username: "bangkit",
		//					Password: "wrongpassword",
		//				}
		//				_, err = userService.Login(&request2)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("login is successful", func() {
		//			It("should create session and create token", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := web.UserLoginRequest{
		//					Username: "bangkit",
		//					Password: "12345678",
		//				}
		//				tokenString, err := userService.Login(&request2)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				session, err := sessionRepository.GetByUsername(request.Username)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(session.Token).ToNot(Equal(tokenString))
		//				Expect(session.Username).To(Equal(request.Username))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get By Username", func() {
		//		When("get by username is success", func() {
		//			It("should return user data", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				user, err := userService.GetByUsername(request.Username)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(user.Username).To(Equal(request.Username))
		//				Expect(user.FullName).To(Equal(request.FullName))
		//				Expect(true).To(Equal(checkPasswordHash(request.Password, user.Password)))
		//				Expect(user.Role).To(Equal(request.Role))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get by username is failed", func() {
		//			It("should return empty user data", func() {
		//				user, err := userService.GetByUsername("wrong")
		//				Expect(err).Should(HaveOccurred())
		//				Expect(user).To(Equal(model.Users{}))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get All", func() {
		//		When("get all data users success", func() {
		//			It("should return all user data", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit2",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err = userService.Register(request2)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				users, err := userService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(users).To(HaveLen(2))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get all data users empty", func() {
		//			It("should return empty user data", func() {
		//				users, err := userService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(users).To(HaveLen(0))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Delete", func() {
		//		When("delete existing user data", func() {
		//			It("should delete data user", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				err = userService.Delete(request.Username)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				user, err := userService.GetByUsername(request.Username)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(user).To(Equal(model.Users{}))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("delete not existing user data", func() {
		//			It("should return error", func() {
		//				err = userService.Delete("empty")
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Update", func() {
		//		When("update user data blank field", func() {
		//			It("should return error", func() {
		//				update := web.UserUpdateRequest{
		//					FullName: "",
		//					Password: "",
		//					Role:     "",
		//					Username: "",
		//				}
		//				_, err = userService.Update(update)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update user data validation error", func() {
		//			It("should return error", func() {
		//				update := web.UserUpdateRequest{
		//					FullName: "Anom",
		//					Password: "123",
		//					Role:     "admin",
		//					Username: "bangkit",
		//				}
		//				_, err = userService.Update(update)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update not existing user data", func() {
		//			It("should return error", func() {
		//				update := web.UserUpdateRequest{
		//					FullName: "Anoman",
		//					Password: "123123123",
		//					Role:     "admin",
		//					Username: "bangkit",
		//				}
		//				_, err = userService.Update(update)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update user data is success", func() {
		//			It("should update data user", func() {
		//				request := &web.UserRegisterRequest{
		//					FullName: "Bangkit Anom",
		//					Username: "bangkit",
		//					Password: "12345678",
		//					Role:     "admin",
		//				}
		//				_, err := userService.Register(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				update := web.UserUpdateRequest{
		//					FullName: "Anom Sedhayu",
		//					Password: "newpassword",
		//					Role:     "user",
		//					Username: "bangkit",
		//				}
		//				_, err = userService.Update(update)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				user, err := userService.GetByUsername(request.Username)
		//				Expect(user.Username).To(Equal(update.Username))
		//				Expect(user.Role).To(Equal(update.Role))
		//				Expect(user.FullName).To(Equal(update.FullName))
		//				Expect(true).To(Equal(checkPasswordHash(update.Password, user.Password)))
		//
		//				err = db.Reset(connection, "users")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//})

		/*
			########################################################################################
											CATEGORY SERVICE
			########################################################################################
		*/
		//Describe("Category Service", func() {
		//	Describe("Add Category", func() {
		//		When("add category is successful", func() {
		//			It("should add category to database", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				category, err := categoryService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(category.ID).To(Equal(request.ID))
		//				Expect(category.Name).To(Equal(request.Name))
		//				Expect(category.Specification).To(Equal(request.Specification))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("add category with blank field", func() {
		//			It("should return error", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "",
		//					Specification: "",
		//				}
		//				category, err := categoryService.Add(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(category).To(Equal(model.Categories{}))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("add category with duplicate id", func() {
		//			It("should return error", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err := categoryService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				category, err := categoryService.Add(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(category).To(Equal(model.Categories{}))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get By ID", func() {
		//		When("get by id is success", func() {
		//			It("should return category data", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err := categoryService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				category, err := categoryService.GetByID(request.ID)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(category.ID).To(Equal(request.ID))
		//				Expect(category.Name).To(Equal(request.Name))
		//				Expect(category.Specification).To(Equal(request.Specification))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get by id is failed", func() {
		//			It("should return empty category data", func() {
		//				category, err := categoryService.GetByID(1)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(category).To(Equal(model.Categories{}))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get All", func() {
		//		When("get all data category success", func() {
		//			It("should return all category data", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err := categoryService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := &web.CategoryAddRequest{
		//					ID:            2,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err = categoryService.Add(request2)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				category, err := categoryService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(category).To(HaveLen(2))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get all data category empty", func() {
		//			It("should return empty category data", func() {
		//				categories, err := categoryService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(categories).To(HaveLen(0))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Delete", func() {
		//		When("delete existing category data", func() {
		//			It("should delete data user", func() {
		//				request := &web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err := categoryService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				err = categoryService.Delete(request.ID)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				category, err := categoryService.GetByID(request.ID)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(category).To(Equal(model.Categories{}))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("delete not existing category data", func() {
		//			It("should return error", func() {
		//				err = categoryService.Delete(1)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Update", func() {
		//		When("update category data blank field", func() {
		//			It("should return error", func() {
		//				request := web.CategoryUpdateRequest{
		//					ID:            1,
		//					Name:          "",
		//					Specification: "",
		//				}
		//				_, err = categoryService.Update(request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update not existing category data", func() {
		//			It("should return error", func() {
		//				request := web.CategoryUpdateRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3090",
		//				}
		//				_, err = categoryService.Update(request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update category data is success", func() {
		//			It("should update data category", func() {
		//				request := web.CategoryAddRequest{
		//					ID:            1,
		//					Name:          "VGA",
		//					Specification: "RTX 3060",
		//				}
		//				_, err := categoryService.Add(&request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				update := web.CategoryUpdateRequest{
		//					ID:            1,
		//					Name:          "VGA2",
		//					Specification: "RTX 3090",
		//				}
		//				_, err = categoryService.Update(update)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				category, err := categoryService.GetByID(request.ID)
		//				Expect(category.ID).To(Equal(update.ID))
		//				Expect(category.Name).To(Equal(update.Name))
		//				Expect(category.Specification).To(Equal(update.Specification))
		//
		//				err = db.Reset(connection, "categories")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//})

		/*
			########################################################################################
											ITEM SERVICE
			########################################################################################
		*/
		//Describe("Item Service", func() {
		//	Describe("Add Item", func() {
		//		When("add item is successful", func() {
		//			It("should add item to database", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "VGA",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				item, err := itemService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(item.Name).To(Equal(request.Name))
		//				Expect(item.CategoryID).To(Equal(request.CategoryID))
		//				Expect(item.Quantity).To(Equal(request.Quantity))
		//				Expect(item.Price).To(Equal(request.Price))
		//				Expect(item.Description).To(Equal(request.Description))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("add item with blank field", func() {
		//			It("should return error", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "",
		//				}
		//				item, err := itemService.Add(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(item).To(Equal(model.Items{}))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get By ID", func() {
		//		When("get by id is success", func() {
		//			It("should return item data", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "VGA",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				_, err := itemService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				item, err := itemService.GetByID(1)
		//				Expect(item.Name).To(Equal(request.Name))
		//				Expect(item.CategoryID).To(Equal(request.CategoryID))
		//				Expect(item.Quantity).To(Equal(request.Quantity))
		//				Expect(item.Price).To(Equal(request.Price))
		//				Expect(item.Description).To(Equal(request.Description))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get by id is failed", func() {
		//			It("should return empty item data", func() {
		//				item, err := itemService.GetByID(1)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(item).To(Equal(model.Items{}))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get All", func() {
		//		When("get all data item success", func() {
		//			It("should return all item data", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "VGA",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				_, err := itemService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := web.ItemAddRequest{
		//					Name:        "PC",
		//					CategoryID:  2,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				_, err = itemService.Add(request2)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				item, err := itemService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(item).To(HaveLen(2))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get all data item empty", func() {
		//			It("should return empty item data", func() {
		//				items, err := itemService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(items).To(HaveLen(0))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Delete", func() {
		//		When("delete existing item data", func() {
		//			It("should delete data item", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "VGA",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				_, err := itemService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				err = itemService.Delete(1)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				item, err := itemService.GetByID(1)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(item).To(Equal(model.Items{}))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("delete not existing item data", func() {
		//			It("should return error", func() {
		//				err = itemService.Delete(1)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Update", func() {
		//		When("update item data blank field", func() {
		//			It("should return error", func() {
		//				request := web.ItemUpdateRequest{
		//					Name:        "",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "",
		//				}
		//				_, err := itemService.Update(request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update not existing item data", func() {
		//			It("should return error", func() {
		//				request := web.ItemUpdateRequest{
		//					ID:          1,
		//					Name:        "",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "",
		//				}
		//				_, err := itemService.Update(request)
		//				Expect(err).Should(HaveOccurred())
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("update item data is success", func() {
		//			It("should update data item", func() {
		//				request := web.ItemAddRequest{
		//					Name:        "VGA",
		//					CategoryID:  1,
		//					Quantity:    10,
		//					Price:       500.00,
		//					Description: "RTX 3060",
		//				}
		//				_, err := itemService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				update := web.ItemUpdateRequest{
		//					ID:          1,
		//					Name:        "VGN",
		//					CategoryID:  2,
		//					Quantity:    12,
		//					Price:       5000.00,
		//					Description: "RTX 3060 new",
		//				}
		//				_, err = itemService.Update(update)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				item, err := itemService.GetByID(1)
		//				Expect(item.ID).To(Equal(update.ID))
		//				Expect(item.Name).To(Equal(update.Name))
		//				Expect(item.CategoryID).To(Equal(update.CategoryID))
		//				Expect(item.Quantity).To(Equal(update.Quantity))
		//				Expect(item.Price).To(Equal(update.Price))
		//				Expect(item.Description).To(Equal(update.Description))
		//
		//				err = db.Reset(connection, "items")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//})

		/*
			########################################################################################
											ACTIVITY SERVICE
			########################################################################################
		*/
		//Describe("Activity Service", func() {
		//	Describe("Add Activity", func() {
		//		When("add activity is successful", func() {
		//			It("should add activity to database", func() {
		//				request := web.ActivityAddRequest{
		//					ItemID:        1,
		//					Action:        "PUT",
		//					QuantityChane: -5,
		//					Timestamp:     time.Now(),
		//					PerformedBy:   1,
		//				}
		//				activity, err := activityService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(activity.ItemID).To(Equal(request.ItemID))
		//				Expect(activity.Action).To(Equal(request.Action))
		//				Expect(activity.QuantityChange).To(Equal(request.QuantityChane))
		//				Expect(activity.Timestamp).To(Equal(request.Timestamp))
		//				Expect(activity.PerformedBy).To(Equal(request.PerformedBy))
		//
		//				err = db.Reset(connection, "activities")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("add activity with blank field", func() {
		//			It("should return error", func() {
		//				request := web.ActivityAddRequest{
		//					Action:        "",
		//					QuantityChane: -5,
		//					Timestamp:     time.Now(),
		//					PerformedBy:   1,
		//				}
		//				activity, err := activityService.Add(request)
		//				Expect(err).Should(HaveOccurred())
		//				Expect(activity).To(Equal(model.Activities{}))
		//
		//				err = db.Reset(connection, "activities")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//
		//	Describe("Get All", func() {
		//		When("get all data activity success", func() {
		//			It("should return all activity data", func() {
		//				request := web.ActivityAddRequest{
		//					ItemID:        1,
		//					Action:        "PUT",
		//					QuantityChane: -5,
		//					Timestamp:     time.Now(),
		//					PerformedBy:   1,
		//				}
		//				_, err := activityService.Add(request)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				request2 := web.ActivityAddRequest{
		//					ItemID:        2,
		//					Action:        "PUT",
		//					QuantityChane: 10,
		//					Timestamp:     time.Now(),
		//					PerformedBy:   1,
		//				}
		//				_, err = activityService.Add(request2)
		//				Expect(err).ShouldNot(HaveOccurred())
		//
		//				activity, err := activityService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(activity).To(HaveLen(2))
		//
		//				err = db.Reset(connection, "activities")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//
		//		When("get all data activity empty", func() {
		//			It("should return empty activity data", func() {
		//				activities, err := activityService.GetAll()
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(activities).To(HaveLen(0))
		//
		//				err = db.Reset(connection, "activities")
		//				Expect(err).ShouldNot(HaveOccurred())
		//			})
		//		})
		//	})
		//})
	})

	Describe("Controller", func() {
		Describe("User Controller", func() {
			Describe("Register", func() {
				When("user registration is success", func() {
					It("should save data user to database", func() {
						userRequest := web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "rahasia",
							Role:     "admin",
						}

						body, _ := json.Marshal(userRequest)
						writer := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
						request.Header.Set("Content-Type", "application/json")

						apiServer.ServeHTTP(writer, request)
						//err := json.Unmarshal(writer.Body.Bytes(), &)
					})
				})
			})
		})
	})
})

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
