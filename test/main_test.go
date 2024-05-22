package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"inventory-management-system/controller"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"inventory-management-system/service"
	"net/http"
	"net/http/httptest"

	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/app"
	model "inventory-management-system/model/domain"
)

func hashPassword(password string) (string, error) {
	bytess, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytess), err
}

func login(mux *gin.Engine, userRepository repository.UserRepository) *http.Cookie {
	pwd, _ := hashPassword("testing123")
	userRepository.Add(model.Users{
		FullName: "testing",
		Username: "testing",
		Password: pwd,
	})

	loginR := web.UserLoginRequest{
		Username: "testing",
		Password: "testing123",
	}

	body, _ := json.Marshal(loginR)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(w, r)

	var cookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session_token" {
			cookie = c
		}
	}

	return cookie
}

var _ = Describe("Digital Inventory Management API", func() {
	var apiServer *gin.Engine
	var userRepository repository.UserRepository
	var categoryRepository repository.CategoryRepository
	var itemRepository repository.ItemRepository
	var reportRepository repository.ReportRepository
	var sessionRepository repository.SessionRepository
	var userService service.UserService
	var categoryService service.CategoryService
	var itemService service.ItemService
	var reportService service.ReportService
	var userController controller.UserController

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

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)

		validate := validator.New()
		userRepository = repository.NewUserRepository(connection)
		categoryRepository = repository.NewCategoryRepository(connection)
		itemRepository = repository.NewItemRepository(connection)
		reportRepository = repository.NewReportRepository(connection)
		sessionRepository = repository.NewSessionRepository(connection)
		userService = service.NewUserService(userRepository, sessionRepository)
		categoryService = service.NewCategoryService(categoryRepository)
		itemService = service.NewItemService(itemRepository, reportRepository)
		reportService = service.NewReportService(reportRepository)
		userController = controller.NewUserController(userService, validate)

		apiServer = gin.New()
		apiServer = app.UserRouter(apiServer, userController)

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

	Describe("Repository", func() {

		/*
			########################################################################################
											USER REPOSITORY
			########################################################################################
		*/
		Describe("User Repository", func() {
			When("add new user to users table in database postgres", func() {
				It("should save data user to users table in database postgres", func() {
					err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "users")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get user by username is success", func() {
				It("should return data user", func() {
					user := model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					}
					err := userRepository.Add(user)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := userRepository.GetByUsername("bangkit")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(uint(1)))
					Expect(result.FullName).To(Equal(user.FullName))
					Expect(result.Username).To(Equal(user.Username))
					Expect(result.Password).To(Equal(user.Password))
					Expect(result.Role).To(Equal(user.Role))

					err = db.Reset(connection, "users")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data users table in database postgres", func() {
				It("should return all data users", func() {
					err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "anom",
						Password: "rahasia",
						Role:     "user",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := userRepository.GetAll()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(2))

					err = db.Reset(connection, "users")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update user data to users table in database postgres", func() {
				It("should save new data user", func() {
					user := model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					}
					err := userRepository.Add(user)
					Expect(err).ShouldNot(HaveOccurred())

					newUser, err := userRepository.GetByUsername(user.Username)
					Expect(err).ShouldNot(HaveOccurred())

					newUser.FullName = "Anom Sedhayu"
					newUser.Password = "newpassword"
					newUser.Role = "user"
					err = userRepository.Update(newUser)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := userRepository.GetByUsername(user.Username)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.FullName).To(Equal(newUser.FullName))
					Expect(result.Username).To(Equal(newUser.Username))
					Expect(result.Password).To(Equal(newUser.Password))
					Expect(result.Role).To(Equal(newUser.Role))

					err = db.Reset(connection, "users")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete user data by username to users table in database postgres", func() {
				It("should soft delete data user by username", func() {
					user := model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					}
					err := userRepository.Add(user)
					Expect(err).ShouldNot(HaveOccurred())

					err = userRepository.Delete(user.Username)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := userRepository.GetByUsername(user.Username)
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Users{}))

					err = db.Reset(connection, "users")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		/*
			########################################################################################
											CATEGORY REPOSITORY
			########################################################################################
		*/
		Describe("Category Repository", func() {
			When("add new category to categories table in database postgres", func() {
				It("should save data category to categories table in database postgres", func() {
					err := categoryRepository.Add(model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get category by id is success", func() {
				It("should return data category", func() {
					category := model.Categories{
						Name: "VGA",
					}
					err := categoryRepository.Add(category)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(int(1)))
					Expect(result.Name).To(Equal(category.Name))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data categories table in database postgres", func() {
				It("should return all data categories", func() {
					err := categoryRepository.Add(model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = categoryRepository.Add(model.Categories{
						Name: "Monitor",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := categoryRepository.GetAll()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(2))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update category data to categories table in database postgres", func() {
				It("should save new data category", func() {
					err := categoryRepository.Add(model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newCategory, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())

					newCategory.Name = "CPU"
					err = categoryRepository.Update(newCategory)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.Name).To(Equal(newCategory.Name))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete category data by id to categories table in database postgres", func() {
				It("should soft delete data category by id", func() {
					err := categoryRepository.Add(model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = categoryRepository.Delete(1)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepository.GetByID(1)
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Categories{}))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		/*
			########################################################################################
											ITEM REPOSITORY
			########################################################################################
		*/
		Describe("Item Repository", func() {
			When("add new item to items table in database postgres", func() {
				It("should save data item to items table in database postgres", func() {
					err := itemRepository.Add(model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get item by item id is success", func() {
				It("should return data item", func() {
					item := model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					}
					err := itemRepository.Add(item)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(int(1)))
					Expect(result.Name).To(Equal(item.Name))
					Expect(result.CategoryID).To(Equal(item.CategoryID))
					Expect(result.Quantity).To(Equal(item.Quantity))
					Expect(result.Price).To(Equal(item.Price))
					Expect(result.Specification).To(Equal(item.Specification))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data items table in database postgres", func() {
				It("should return all data items", func() {
					err := itemRepository.Add(model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = itemRepository.Add(model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := itemRepository.GetAll()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(2))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update item data to items table in database postgres", func() {
				It("should save new data item", func() {
					err := itemRepository.Add(model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newItem, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())

					newItem.Name = "VGA"
					newItem.CategoryID = 2
					newItem.Quantity = 5
					newItem.Price = 4500000.00
					newItem.Specification = "desc"
					err = itemRepository.Update(newItem)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(int(1)))
					Expect(result.Name).To(Equal(newItem.Name))
					Expect(result.CategoryID).To(Equal(newItem.CategoryID))
					Expect(result.Quantity).To(Equal(newItem.Quantity))
					Expect(result.Price).To(Equal(newItem.Price))
					Expect(result.Specification).To(Equal(newItem.Specification))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete item data by item id to items table in database postgres", func() {
				It("should soft delete data item by id", func() {
					err := itemRepository.Add(model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = itemRepository.Delete(1)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := itemRepository.GetByItemID(1)
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Items{}))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data items table by category id in database postgres", func() {
				It("should return all data items by category id", func() {
					err := itemRepository.Add(model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = itemRepository.Add(model.Items{
						Name:          "VGA2",
						CategoryID:    2,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := itemRepository.GetByCategoryID(2)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(1))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		/*
			########################################################################################
											ACTIVITY REPOSITORY
			########################################################################################
		*/
		Describe("Item Repository", func() {
			When("add new activity to activities table in database postgres", func() {
				It("should save data activity to activities table in database postgres", func() {
					err := reportRepository.AddActivity(model.Activities{
						ItemID:         1,
						Action:         "POST",
						QuantityChange: 5,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "activities")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data activities table in database postgres", func() {
				It("should return all data activities", func() {
					err := reportRepository.AddActivity(model.Activities{
						ItemID:         1,
						Action:         "POST",
						QuantityChange: 5,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = reportRepository.AddActivity(model.Activities{
						ItemID:         2,
						Action:         "POST",
						QuantityChange: -2,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := reportRepository.GetAllActivity()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(2))

					err = db.Reset(connection, "activities")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		/*
			########################################################################################
											SESSION REPOSITORY
			########################################################################################
		*/
		Describe("Session Repository", func() {
			When("add new session to sessions table in database postgres", func() {
				It("should save data session to sessions table in database postgres", func() {
					err := sessionRepository.Add(model.Sessions{
						Username:  "bangkit",
						Token:     "token",
						ExpiresAt: time.Now().Add(5 * time.Minute),
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "sessions")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get session by username is success", func() {
				It("should return data session", func() {
					session := model.Sessions{
						Username:  "bangkit",
						Token:     "token",
						ExpiresAt: time.Now().Add(5 * time.Minute),
					}
					err := sessionRepository.Add(session)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := sessionRepository.GetByUsername("bangkit")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(uint(1)))
					Expect(result.Username).To(Equal(session.Username))
					Expect(result.Token).To(Equal(session.Token))

					err = db.Reset(connection, "sessions")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete session data by username to sessions table in database postgres", func() {
				It("should soft delete data session by username", func() {
					err := sessionRepository.Add(model.Sessions{
						Username:  "bangkit",
						Token:     "token",
						ExpiresAt: time.Now().Add(5 * time.Minute),
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = sessionRepository.Delete("bangkit")
					Expect(err).ShouldNot(HaveOccurred())

					result, err := sessionRepository.GetByUsername("bangkit")
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Sessions{}))

					err = db.Reset(connection, "sessions")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	//====================================================================================================
	//====================================================================================================
	//====================================================================================================

	Describe("Service", func() {

		/*
			########################################################################################
											USER SERVICE
			########################################################################################
		*/
		Describe("User Service", func() {
			Describe("Register", func() {
				When("register is successful", func() {
					It("should register user", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						user, err := userRepository.GetByUsername(request.Username)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(user.FullName).To(Equal(request.FullName))
						Expect(user.Username).To(Equal(request.Username))
						Expect(true).To(Equal(checkPasswordHash(request.Password, user.Password)))
						Expect(user.Role).To(Equal(request.Role))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("register with duplicate username", func() {
					It("should return error", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						err = userService.Register(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Login", func() {
				When("login with blank field", func() {
					It("should return error", func() {
						request := web.UserLoginRequest{
							Username: "",
							Password: "",
						}
						_, err := userService.Login(&request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("login with wrong username", func() {
					It("should return error", func() {
						request := web.UserLoginRequest{
							Username: "wrongusername",
							Password: "wrongpassword",
						}
						_, err := userService.Login(&request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("login with wrong password", func() {
					It("should return error", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := web.UserLoginRequest{
							Username: "bangkit",
							Password: "wrongpassword",
						}
						_, err = userService.Login(&request2)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("login is successful", func() {
					It("should create session and create token", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := web.UserLoginRequest{
							Username: "bangkit",
							Password: "12345678",
						}
						tokenString, err := userService.Login(&request2)
						Expect(err).ShouldNot(HaveOccurred())

						session, err := sessionRepository.GetByUsername(request.Username)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(session.Token).ToNot(Equal(tokenString))
						Expect(session.Username).To(Equal(request.Username))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get By Username", func() {
				When("get by username is success", func() {
					It("should return user data", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						user, err := userService.GetByUsername(request.Username)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(user.Username).To(Equal(request.Username))
						Expect(user.FullName).To(Equal(request.FullName))
						Expect(user.Role).To(Equal(request.Role))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get by username is failed", func() {
					It("should return empty user data", func() {
						user, err := userService.GetByUsername("wrong")
						Expect(err).Should(HaveOccurred())
						Expect(user).To(Equal(model.Users{}))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get All", func() {
				When("get all data users success", func() {
					It("should return all user data", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit2",
							Password: "12345678",
							Role:     "admin",
						}
						err = userService.Register(request2)
						Expect(err).ShouldNot(HaveOccurred())

						users, err := userService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(users).To(HaveLen(2))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data users empty", func() {
					It("should return empty user data", func() {
						users, err := userService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(users).To(HaveLen(0))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("delete existing user data", func() {
					It("should delete data user", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						err = userService.Delete(request.Username)
						Expect(err).ShouldNot(HaveOccurred())

						user, err := userService.GetByUsername(request.Username)
						Expect(err).Should(HaveOccurred())
						Expect(user).To(Equal(model.Users{}))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing user data", func() {
					It("should return error", func() {
						err = userService.Delete("empty")
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Update", func() {
				When("update user data blank field", func() {
					It("should return error", func() {
						update := web.UserUpdateRequest{
							FullName: "",
							Password: "",
							Role:     "",
							Username: "",
						}
						err = userService.Update(update)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update user data validation error", func() {
					It("should return error", func() {
						update := web.UserUpdateRequest{
							FullName: "Anom",
							Password: "123",
							Role:     "admin",
							Username: "bangkit",
						}
						err = userService.Update(update)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update not existing user data", func() {
					It("should return error", func() {
						update := web.UserUpdateRequest{
							FullName: "Anoman",
							Password: "123123123",
							Role:     "admin",
							Username: "bangkit",
						}
						err = userService.Update(update)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update user data is success", func() {
					It("should update data user", func() {
						request := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "12345678",
							Role:     "admin",
						}
						err := userService.Register(request)
						Expect(err).ShouldNot(HaveOccurred())

						update := web.UserUpdateRequest{
							FullName: "Anom Sedhayu",
							Password: "newpassword",
							Role:     "user",
							Username: "bangkit",
						}
						err = userService.Update(update)
						Expect(err).ShouldNot(HaveOccurred())

						user, err := userService.GetByUsername(request.Username)
						Expect(user.Username).To(Equal(update.Username))
						Expect(user.Role).To(Equal(update.Role))
						Expect(user.FullName).To(Equal(update.FullName))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})

		/*
			########################################################################################
											CATEGORY SERVICE
			########################################################################################
		*/
		Describe("Category Service", func() {
			Describe("Add Category", func() {
				When("add category is successful", func() {
					It("should add category to database", func() {
						request := &web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())
						Expect("VGA").To(Equal(request.Name))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("add category with duplicate id", func() {
					It("should return error", func() {
						request := &web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						err = categoryService.Add(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get By ID", func() {
				When("get by id is success", func() {
					It("should return category data", func() {
						request := &web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						category, err := categoryService.GetByID(1)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(category.Name).To(Equal(request.Name))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get by id is failed", func() {
					It("should return empty category data", func() {
						category, err := categoryService.GetByID(1)
						Expect(err).Should(HaveOccurred())
						Expect(category).To(Equal(model.Categories{}))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get All", func() {
				When("get all data category success", func() {
					It("should return all category data", func() {
						request := &web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := &web.CategoryAddRequest{
							Name: "Monitor",
						}
						err = categoryService.Add(request2)
						Expect(err).ShouldNot(HaveOccurred())

						category, err := categoryService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(category).To(HaveLen(2))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data category empty", func() {
					It("should return empty category data", func() {
						categories, err := categoryService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(categories).To(HaveLen(0))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("delete existing category data", func() {
					It("should delete data user", func() {
						request := &web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						err = categoryService.Delete(1)
						Expect(err).ShouldNot(HaveOccurred())

						category, err := categoryService.GetByID(1)
						Expect(err).Should(HaveOccurred())
						Expect(category).To(Equal(model.Categories{}))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing category data", func() {
					It("should return error", func() {
						err = categoryService.Delete(1)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Update", func() {
				When("update category data blank field", func() {
					It("should return error", func() {
						request := web.CategoryUpdateRequest{
							ID:   1,
							Name: "",
						}
						err = categoryService.Update(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update not existing category data", func() {
					It("should return error", func() {
						request := web.CategoryUpdateRequest{
							ID:   1,
							Name: "VGA",
						}
						err = categoryService.Update(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update category data is success", func() {
					It("should update data category", func() {
						request := web.CategoryAddRequest{
							Name: "VGA",
						}
						err := categoryService.Add(&request)
						Expect(err).ShouldNot(HaveOccurred())

						update := web.CategoryUpdateRequest{
							ID:   1,
							Name: "VGA2",
						}
						err = categoryService.Update(update)
						Expect(err).ShouldNot(HaveOccurred())

						category, err := categoryService.GetByID(1)
						Expect(category.ID).To(Equal(update.ID))
						Expect(category.Name).To(Equal(update.Name))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})

		/*
			########################################################################################
											ITEM SERVICE
			########################################################################################
		*/
		Describe("Item Service", func() {
			Describe("Add Item", func() {
				When("add item is successful", func() {
					It("should add item to database", func() {
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err := itemService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						item, err := itemRepository.GetByItemID(1)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(item.Name).To(Equal(request.Name))
						Expect(item.CategoryID).To(Equal(request.CategoryID))
						Expect(item.Quantity).To(Equal(request.Quantity))
						Expect(item.Price).To(Equal(request.Price))
						Expect(item.Specification).To(Equal(request.Specification))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get By ID", func() {
				When("get by id is success", func() {
					It("should return item data", func() {
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err := itemService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						item, err := itemService.GetByID(1)
						Expect(item.Name).To(Equal(request.Name))
						Expect(item.CategoryID).To(Equal(request.CategoryID))
						Expect(item.Quantity).To(Equal(request.Quantity))
						Expect(item.Price).To(Equal(request.Price))
						Expect(item.Specification).To(Equal(request.Specification))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get by id is failed", func() {
					It("should return empty item data", func() {
						item, err := itemService.GetByID(1)
						Expect(err).Should(HaveOccurred())
						Expect(item).To(Equal(model.Items{}))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get All", func() {
				When("get all data item success", func() {
					It("should return all item data", func() {
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err := itemService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := web.ItemAddRequest{
							Name:          "PC",
							CategoryID:    2,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err = itemService.Add(request2)
						Expect(err).ShouldNot(HaveOccurred())

						item, err := itemService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(item).To(HaveLen(2))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data item empty", func() {
					It("should return empty item data", func() {
						items, err := itemService.GetAll()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(items).To(HaveLen(0))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("delete existing item data", func() {
					It("should delete data item", func() {
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err := itemService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						err = itemService.Delete(1)
						Expect(err).ShouldNot(HaveOccurred())

						item, err := itemService.GetByID(1)
						Expect(err).Should(HaveOccurred())
						Expect(item).To(Equal(model.Items{}))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing item data", func() {
					It("should return error", func() {
						err = itemService.Delete(1)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Update", func() {
				When("update item data blank field", func() {
					It("should return error", func() {
						request := web.ItemUpdateRequest{
							Name:          "",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "",
						}
						err := itemService.Update(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update not existing item data", func() {
					It("should return error", func() {
						request := web.ItemUpdateRequest{
							ID:            1,
							Name:          "",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "",
						}
						err := itemService.Update(request)
						Expect(err).Should(HaveOccurred())

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update item data is success", func() {
					It("should update data item", func() {
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						err := itemService.Add(request)
						Expect(err).ShouldNot(HaveOccurred())

						update := web.ItemUpdateRequest{
							ID:            1,
							Name:          "VGN",
							CategoryID:    2,
							Quantity:      12,
							Price:         5000.00,
							Specification: "RTX 3060 new",
						}
						err = itemService.Update(update)
						Expect(err).ShouldNot(HaveOccurred())

						item, err := itemService.GetByID(1)
						Expect(item.ID).To(Equal(update.ID))
						Expect(item.Name).To(Equal(update.Name))
						Expect(item.CategoryID).To(Equal(update.CategoryID))
						Expect(item.Quantity).To(Equal(update.Quantity))
						Expect(item.Price).To(Equal(update.Price))
						Expect(item.Specification).To(Equal(update.Specification))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})

		/*
			########################################################################################
											ACTIVITY SERVICE
			########################################################################################
		*/
		Describe("Activity Service", func() {
			Describe("Get All", func() {
				When("get all data activity success", func() {
					It("should return all activity data", func() {
						request := model.Activities{
							ItemID:         1,
							Action:         "PUT",
							QuantityChange: 5,
							Timestamp:      time.Now(),
							PerformedBy:    1,
						}
						err := reportRepository.AddActivity(request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := model.Activities{
							ItemID:         1,
							Action:         "PUT",
							QuantityChange: 5,
							Timestamp:      time.Now(),
							PerformedBy:    1,
						}
						err = reportRepository.AddActivity(request2)
						Expect(err).ShouldNot(HaveOccurred())

						activity, err := reportService.GetAllActivity()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(activity).To(HaveLen(2))

						err = db.Reset(connection, "activities")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data activity empty", func() {
					It("should return empty activity data", func() {
						activities, err := reportService.GetAllActivity()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(activities).To(HaveLen(0))

						err = db.Reset(connection, "activities")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("Controller", func() {
		Describe("User Controller", func() {
			Describe("Register", func() {
				When("user registration is success", func() {
					It("should save data user to database", func() {
						userRequest := web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "rahasia1",
							Role:     "admin",
						}

						body, err := json.Marshal(userRequest)
						Expect(err).ShouldNot(HaveOccurred())

						recorder := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
						request.Header.Set("Content-Type", "application/json")
						request.AddCookie(login(apiServer, userRepository))

						apiServer.ServeHTTP(recorder, request)

						response := web.Created{}
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.Code).To(Equal(http.StatusCreated))
						Expect(response.Status).To(Equal("status created"))
						Expect(response.Message).To(Equal("registration successful"))

						err = db.Reset(connection, "users")
						err = db.Reset(connection, "sessions")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("user registration without sending cookie", func() {
					It("should return unauthorized", func() {
						userRequest := web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit",
							Password: "rahasia1",
							Role:     "admin",
						}

						body, err := json.Marshal(userRequest)
						Expect(err).ShouldNot(HaveOccurred())

						recorder := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
						request.Header.Set("Content-Type", "application/json")

						apiServer.ServeHTTP(recorder, request)

						response := web.Unauthorized{}
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.Code).To(Equal(http.StatusUnauthorized))
						Expect(response.Status).To(Equal("status unauthorized"))
						Expect(response.Message).To(Equal("session token is empty"))

						err = db.Reset(connection, "users")
						err = db.Reset(connection, "sessions")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("user registration with blank field", func() {
					It("should return error", func() {
						userRequest := web.UserRegisterRequest{
							FullName: "",
							Username: "",
							Password: "",
							Role:     "",
						}

						body, err := json.Marshal(userRequest)
						Expect(err).ShouldNot(HaveOccurred())

						recorder := httptest.NewRecorder()
						request := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
						request.Header.Set("Content-Type", "application/json")
						request.AddCookie(login(apiServer, userRepository))

						apiServer.ServeHTTP(recorder, request)

						response := web.BadRequestResponse{}
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.Code).To(Equal(http.StatusBadRequest))
						Expect(response.Status).To(Equal("status bad request"))
						Expect(response.Message).To(Equal("validation error"))

						err = db.Reset(connection, "users")
						err = db.Reset(connection, "sessions")
						Expect(err).ShouldNot(HaveOccurred())
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
