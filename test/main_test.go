package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"inventory-management-system/app"
	"inventory-management-system/controller"
	model "inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"inventory-management-system/repository"
	"inventory-management-system/service"
	"net/http"
	"net/http/httptest"
	"time"
)

func hashPassword(password string) (string, error) {
	bytess, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytess), err
}

func login(mux *gin.Engine, handleRepository repository.HandlerRepository) *http.Cookie {
	pwd, _ := hashPassword("testing123")
	handleRepository.Add(&model.Users{
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
	var handlerRepository repository.HandlerRepository
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
		handlerRepository = repository.NewHandlerRepository(connection)
		userService = service.NewUserService(handlerRepository)
		categoryService = service.NewCategoryService(handlerRepository)
		itemService = service.NewItemService(handlerRepository)
		reportService = service.NewReportService(handlerRepository)
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

	Describe("HandlerRepository", func() {
		When("get by username failed", func() {
			It("should return error", func() {
				_, errResponse := userService.GetByUsername("empty")
				Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
			})
		})

		When("add data success", func() {
			It("should save data to database", func() {
				user := model.Users{
					FullName: "bangkit anom",
					Username: "bangkit",
					Password: "rahasia",
					Role:     "admin",
				}
				err := handlerRepository.Add(&user)
				Expect(err).ShouldNot(HaveOccurred())

				userDB := model.Users{}
				err = handlerRepository.GetByUsername(user.Username, &userDB)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(userDB.ID).To(Equal(uint(1)))

				category := model.Categories{
					Name: "VGA",
				}
				err = handlerRepository.Add(&category)
				Expect(err).ShouldNot(HaveOccurred())

				categoryDB := model.Categories{}
				err = handlerRepository.GetByID(1, &categoryDB)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(categoryDB.Name).To(Equal("VGA"))

				err = db.Reset(connection, "users")
				err = db.Reset(connection, "categories")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("update data by username success", func() {
			It("should update data", func() {
				user := model.Users{
					FullName: "bangkit anom",
					Username: "bangkit",
					Password: "rahasia",
					Role:     "admin",
				}
				err := handlerRepository.Add(&user)
				Expect(err).ShouldNot(HaveOccurred())

				newUser := model.Users{
					FullName: "new bangkit anom",
					Username: "bangkit",
					Password: "rahasia",
					Role:     "user",
				}

				err = handlerRepository.UpdateByUsername(newUser.Username, &newUser)
				Expect(err).ShouldNot(HaveOccurred())

				userDB := model.Users{}
				err = handlerRepository.GetByUsername(newUser.Username, &userDB)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(userDB.FullName).To(Equal(newUser.FullName))

				err = db.Reset(connection, "users")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("delete data is success", func() {
			It("should delete data", func() {
				category := model.Categories{
					Name: "VGA",
				}
				err := handlerRepository.Add(&category)
				Expect(err).ShouldNot(HaveOccurred())

				err = handlerRepository.DeleteByID(1, &model.Categories{})
				Expect(err).ShouldNot(HaveOccurred())

				categoryDB := model.Categories{}
				err = handlerRepository.GetByID(1, &categoryDB)
				Expect(err).Should(HaveOccurred())

				err = db.Reset(connection, "users")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("delete data is success", func() {
			It("should delete data", func() {
				user := model.Users{
					FullName: "bangkit anom",
					Username: "bangkit",
					Password: "rahasia",
					Role:     "admin",
				}
				err := handlerRepository.Add(&user)
				Expect(err).ShouldNot(HaveOccurred())

				err = handlerRepository.DeleteByUsername("bangkit", &model.Users{})
				Expect(err).ShouldNot(HaveOccurred())

				userDB := model.Users{}
				err = handlerRepository.GetByUsername("bangkit", &userDB)
				Expect(err).Should(HaveOccurred())

				err = db.Reset(connection, "users")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("get all data success", func() {
			It("should take all data", func() {
				user := model.Users{
					FullName: "bangkit anom",
					Username: "bangkit",
					Password: "rahasia",
					Role:     "admin",
				}
				err := handlerRepository.Add(&user)
				Expect(err).ShouldNot(HaveOccurred())

				user = model.Users{
					FullName: "bangkit anom",
					Username: "bangkit2",
					Password: "rahasia",
					Role:     "admin",
				}
				err = handlerRepository.Add(&user)
				Expect(err).ShouldNot(HaveOccurred())

				users := []model.Users{}
				err = handlerRepository.GetAll(&users)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).To(Equal(2))

				err = db.Reset(connection, "users")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
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
					err := handlerRepository.Add(&model.Users{
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
					err := handlerRepository.Add(&user)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Users{}
					err = handlerRepository.GetByUsername("bangkit", &result)
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
					err := handlerRepository.Add(&model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.Add(&model.Users{
						FullName: "Bangkit Anom",
						Username: "anom",
						Password: "rahasia",
						Role:     "user",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results := []model.Users{}
					err = handlerRepository.GetAll(&results)
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
					err := handlerRepository.Add(&user)
					Expect(err).ShouldNot(HaveOccurred())

					newUser := model.Users{}
					err = handlerRepository.GetByUsername(user.Username, &newUser)
					Expect(err).ShouldNot(HaveOccurred())

					newUser.FullName = "Anom Sedhayu"
					newUser.Password = "newpassword"
					newUser.Role = "user"
					err = handlerRepository.UpdateByUsername(newUser.Username, &newUser)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Users{}
					err = handlerRepository.GetByUsername(user.Username, &result)
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
					err := handlerRepository.Add(&user)
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.DeleteByUsername(user.Username, &model.Users{})
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Users{}
					err = handlerRepository.GetByUsername(user.Username, &result)
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
					err := handlerRepository.Add(&model.Categories{
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
					err := handlerRepository.Add(&category)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Categories{}
					err = handlerRepository.GetByID(1, &result)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(int(1)))
					Expect(result.Name).To(Equal(category.Name))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data categories table in database postgres", func() {
				It("should return all data categories", func() {
					err := handlerRepository.Add(&model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.Add(&model.Categories{
						Name: "Monitor",
					})
					Expect(err).ShouldNot(HaveOccurred())

					results := []model.Categories{}
					err = handlerRepository.GetAll(&results)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).Should(HaveLen(2))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update category data to categories table in database postgres", func() {
				It("should save new data category", func() {
					err := handlerRepository.Add(&model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newCategory := model.Categories{}
					err = handlerRepository.GetByID(1, &newCategory)
					Expect(err).ShouldNot(HaveOccurred())

					newCategory.Name = "CPU"
					err = handlerRepository.UpdateByID(1, &newCategory)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Categories{}
					err = handlerRepository.GetByID(1, &result)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.Name).To(Equal(newCategory.Name))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete category data by id to categories table in database postgres", func() {
				It("should soft delete data category by id", func() {
					err := handlerRepository.Add(&model.Categories{
						Name: "VGA",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.DeleteByID(1, model.Categories{})
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Categories{}
					err = handlerRepository.GetByID(1, &model.Categories{})
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
					err := handlerRepository.Add(&model.Items{
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
					err := handlerRepository.Add(&item)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Items{}
					err = handlerRepository.GetByID(1, &result)
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
					err := handlerRepository.Add(&model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.Add(&model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					items := []model.Items{}
					err = handlerRepository.GetAll(&items)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).Should(HaveLen(2))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("update item data to items table in database postgres", func() {
				It("should save new data item", func() {
					err := handlerRepository.Add(&model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newItem := model.Items{}
					err = handlerRepository.GetByID(1, &newItem)
					Expect(err).ShouldNot(HaveOccurred())

					newItem.Name = "VGA"
					newItem.CategoryID = 2
					newItem.Quantity = 5
					newItem.Price = 4500000.00
					newItem.Specification = "desc"
					err = handlerRepository.UpdateByID(1, &newItem)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Items{}
					err = handlerRepository.GetByID(1, &result)
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
					err := handlerRepository.Add(&model.Items{
						Name:          "VGA2",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.DeleteByID(1, &model.Items{})
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Items{}
					err = handlerRepository.GetByID(1, &result)
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Items{}))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data items table by category id in database postgres", func() {
				It("should return all data items by category id", func() {
					err := handlerRepository.Add(&model.Items{
						Name:          "VGA",
						CategoryID:    1,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.Add(&model.Items{
						Name:          "VGA2",
						CategoryID:    2,
						Quantity:      10,
						Price:         5000000.00,
						Specification: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					items := []model.Items{}
					err = handlerRepository.GetByCategoryID(2, &items)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(items).Should(HaveLen(1))

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
					err := handlerRepository.Add(&model.Activities{
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
					err := handlerRepository.Add(&model.Activities{
						ItemID:         1,
						Action:         "POST",
						QuantityChange: 5,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.Add(&model.Activities{
						ItemID:         2,
						Action:         "POST",
						QuantityChange: -2,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					results := []model.Activities{}
					err = handlerRepository.GetAll(&results)
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
					err := handlerRepository.Add(&model.Sessions{
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
					err := handlerRepository.Add(&session)
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Sessions{}
					err = handlerRepository.GetByUsername("bangkit", &result)
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
					err := handlerRepository.Add(&model.Sessions{
						Username:  "bangkit",
						Token:     "token",
						ExpiresAt: time.Now().Add(5 * time.Minute),
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = handlerRepository.DeleteByUsername("bangkit", &model.Sessions{})
					Expect(err).ShouldNot(HaveOccurred())

					result := model.Sessions{}
					err = handlerRepository.GetByUsername("bangkit", &result)
					Expect(err).Should(HaveOccurred())
					Expect(result).To(Equal(model.Sessions{}))

					err = db.Reset(connection, "sessions")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						user := model.Users{}
						err = handlerRepository.GetByUsername(request.Username, &user)
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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						errResponse = userService.Register(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						_, errResponse := userService.Login(&request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						_, errResponse := userService.Login(&request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						request2 := web.UserLoginRequest{
							Username: "bangkit",
							Password: "wrongpassword",
						}
						_, errResponse = userService.Login(&request2)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						request2 := web.UserLoginRequest{
							Username: "bangkit",
							Password: "12345678",
						}
						tokenString, errResponse := userService.Login(&request2)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						session := model.Sessions{}
						err = handlerRepository.GetByUsername(request.Username, &session)
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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						user, errResponse := userService.GetByUsername(request.Username)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(user.Username).To(Equal(request.Username))
						Expect(user.FullName).To(Equal(request.FullName))
						Expect(user.Role).To(Equal(request.Role))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get by username is failed", func() {
					It("should return empty user data", func() {
						user, errResponse := userService.GetByUsername("wrong")
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						request2 := &web.UserRegisterRequest{
							FullName: "Bangkit Anom",
							Username: "bangkit2",
							Password: "12345678",
							Role:     "admin",
						}
						errResponse = userService.Register(request2)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						users, errResponse := userService.GetAll()
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(users).To(HaveLen(2))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data users empty", func() {
					It("should return empty user data", func() {
						users, errResponse := userService.GetAll()
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						errResponse = userService.Delete(request.Username)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						user, errResponse := userService.GetByUsername(request.Username)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
						Expect(user).To(Equal(model.Users{}))

						err = db.Reset(connection, "users")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing user data", func() {
					It("should return error", func() {
						errResponse := userService.Delete("empty")
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Update(update)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Update(update)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Update(update)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := userService.Register(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						update := web.UserUpdateRequest{
							FullName: "Anom Sedhayu",
							Password: "newpassword",
							Role:     "user",
							Username: "bangkit",
						}
						errResponse = userService.Update(update)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						user, errResponse := userService.GetByUsername(request.Username)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
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
						errResponse := categoryService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
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
						errResponse := categoryService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						errResponse = categoryService.Add(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := categoryService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						category, errResponse := categoryService.GetByID(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(category.Name).To(Equal(request.Name))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get by id is failed", func() {
					It("should return empty category data", func() {
						category, errResponse := categoryService.GetByID(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
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
						errResponse := categoryService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						request2 := &web.CategoryAddRequest{
							Name: "Monitor",
						}
						errResponse = categoryService.Add(request2)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						category, errResponse := categoryService.GetAll()
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(category).To(HaveLen(2))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data category empty", func() {
					It("should return empty category data", func() {
						categories, errResponse := categoryService.GetAll()
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
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
						errResponse := categoryService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						errResponse = categoryService.Delete(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						category, errResponse := categoryService.GetByID(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
						Expect(category).To(Equal(model.Categories{}))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing category data", func() {
					It("should return error", func() {
						errResponse := categoryService.Delete(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := categoryService.Update(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := categoryService.Update(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

						err = db.Reset(connection, "categories")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update category data is success", func() {
					It("should update data category", func() {
						request := web.CategoryAddRequest{
							Name: "VGA",
						}
						errResponse := categoryService.Add(&request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						update := web.CategoryUpdateRequest{
							ID:   1,
							Name: "VGA2",
						}
						errResponse = categoryService.Update(update)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						category, errResponse := categoryService.GetByID(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
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
						categoryService.Add(&web.CategoryAddRequest{
							Name: "VGA",
						})
						request := web.ItemAddRequest{
							Name:          "RTX 3060",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse := itemService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						item := model.Items{}
						err := handlerRepository.GetByID(1, &item)
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
						categoryService.Add(&web.CategoryAddRequest{
							Name: "VGA",
						})
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse := itemService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						item, errResponse := itemService.GetByID(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
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
						item, errResponse := itemService.GetByID(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
						Expect(item).To(Equal(model.Items{}))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Get All", func() {
				When("get all data item success", func() {
					It("should return all item data", func() {
						categoryService.Add(&web.CategoryAddRequest{
							Name: "VGA",
						})
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse := itemService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						request2 := web.ItemAddRequest{
							Name:          "PC",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse = itemService.Add(request2)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						item, errResponse := itemService.GetAll()
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(item).To(HaveLen(2))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data item empty", func() {
					It("should return empty item data", func() {
						items, errResponse := itemService.GetAll()
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
						Expect(items).To(HaveLen(0))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("delete existing item data", func() {
					It("should delete data item", func() {
						categoryService.Add(&web.CategoryAddRequest{
							Name: "VGA",
						})
						request := web.ItemAddRequest{
							Name:          "VGA",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse := itemService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						errResponse = itemService.Delete(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						item, errResponse := itemService.GetByID(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
						Expect(item).To(Equal(model.Items{}))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("delete not existing item data", func() {
					It("should return error", func() {
						errResponse := itemService.Delete(1)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := itemService.Update(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

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
						errResponse := itemService.Update(request)
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))

						err = db.Reset(connection, "items")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("update item data is success", func() {
					It("should update data item", func() {
						categoryService.Add(&web.CategoryAddRequest{
							Name: "VGA",
						})
						request := web.ItemAddRequest{
							Name:          "RTX 3060",
							CategoryID:    1,
							Quantity:      10,
							Price:         500.00,
							Specification: "RTX 3060",
						}
						errResponse := itemService.Add(request)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						update := web.ItemUpdateRequest{
							ID:            1,
							Name:          "RTX 3090",
							CategoryID:    1,
							Quantity:      12,
							Price:         5000.00,
							Specification: "RTX 3060 new",
						}
						errResponse = itemService.Update(update)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))

						item, errResponse := itemService.GetByID(1)
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
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
						err := handlerRepository.Add(&request)
						Expect(err).ShouldNot(HaveOccurred())

						request2 := model.Activities{
							ItemID:         1,
							Action:         "PUT",
							QuantityChange: 5,
							Timestamp:      time.Now(),
							PerformedBy:    1,
						}
						err = handlerRepository.Add(&request2)
						Expect(err).ShouldNot(HaveOccurred())

						activity, errResponse := reportService.GetAllActivity()
						Expect(errResponse).To(Equal(web.ErrorResponse{}))
						Expect(activity).To(HaveLen(2))

						err = db.Reset(connection, "activities")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("get all data activity empty", func() {
					It("should return empty activity data", func() {
						activities, errResponse := reportService.GetAllActivity()
						Expect(errResponse).ToNot(Equal(web.ErrorResponse{}))
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
						request.AddCookie(login(apiServer, handlerRepository))

						apiServer.ServeHTTP(recorder, request)

						response := web.SuccessResponseMessage{}
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.Code).To(Equal(http.StatusCreated))
						Expect(response.Status).To(Equal("status created"))
						Expect(response.Message).To(Equal("register user success"))

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

						response := web.ErrorResponse{}
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
						request.AddCookie(login(apiServer, handlerRepository))

						apiServer.ServeHTTP(recorder, request)

						response := web.ErrorResponse{}
						err = json.Unmarshal(recorder.Body.Bytes(), &response)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(response.Code).To(Equal(http.StatusBadRequest))
						Expect(response.Status).To(Equal("status bad request"))

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
