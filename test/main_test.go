package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"inventory-management-system/app"
	"inventory-management-system/model"
	"inventory-management-system/repository"
	"time"
)

var _ = Describe("Digital Inventory Management API", func() {
	var userRepository repository.UserRepository
	var categoryRepository repository.CategoryRepository
	var itemRepository repository.ItemRepository
	var activityRepository repository.ActivityRepository
	var sessionRepository repository.SessionRepository

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

	userRepository = repository.NewUserRepository(connection)
	categoryRepository = repository.NewCategoryRepository(connection)
	itemRepository = repository.NewItemRepository(connection)
	activityRepository = repository.NewActivityRepository(connection)
	sessionRepository = repository.NewSessionRepository(connection)

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

	Describe("Repository", func() {

		/*
			########################################################################################
											USER REPOSITORY
			########################################################################################
		*/
		Describe("User Repository", func() {
			When("add new user to users table in database postgres", func() {
				It("should save data user to users table in database postgres", func() {
					_, err := userRepository.Add(model.Users{
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
					user, err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
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
					_, err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
					Expect(err).ShouldNot(HaveOccurred())

					_, err = userRepository.Add(model.Users{
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
					user, err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newUser, err := userRepository.GetByUsername(user.Username)
					Expect(err).ShouldNot(HaveOccurred())

					newUser.FullName = "Anom Sedhayu"
					newUser.Password = "newpassword"
					newUser.Role = "user"
					_, err = userRepository.Update(newUser)
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
					user, err := userRepository.Add(model.Users{
						FullName: "Bangkit Anom",
						Username: "bangkit",
						Password: "rahasia",
						Role:     "admin",
					})
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
					_, err := categoryRepository.Add(model.Categories{
						Name:          "VGA",
						Specification: "RTX-3060, RAM 4 GB",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get category by id is success", func() {
				It("should return data category", func() {
					category, err := categoryRepository.Add(model.Categories{
						Name:          "VGA",
						Specification: "RTX-3060, RAM 4 GB",
					})
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(uint(1)))
					Expect(result.Name).To(Equal(category.Name))
					Expect(result.Specification).To(Equal(category.Specification))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data categories table in database postgres", func() {
				It("should return all data categories", func() {
					_, err := categoryRepository.Add(model.Categories{
						Name:          "VGA",
						Specification: "RTX-3060, RAM 4 GB",
					})
					Expect(err).ShouldNot(HaveOccurred())

					_, err = categoryRepository.Add(model.Categories{
						Name:          "Monitor",
						Specification: "14 in Ajhua",
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
					_, err := categoryRepository.Add(model.Categories{
						Name:          "VGA",
						Specification: "RTX-3060, RAM 4 GB",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newCategory, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())

					newCategory.Name = "CPU"
					newCategory.Specification = "4 core"
					_, err = categoryRepository.Update(newCategory)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepository.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.Name).To(Equal(newCategory.Name))
					Expect(result.Specification).To(Equal(newCategory.Specification))

					err = db.Reset(connection, "categories")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete category data by id to categories table in database postgres", func() {
				It("should soft delete data category by id", func() {
					_, err := categoryRepository.Add(model.Categories{
						Name:          "VGA",
						Specification: "RTX-3060, RAM 4 GB",
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
					_, err := itemRepository.Add(model.Items{
						Name:        "VGA",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get item by item id is success", func() {
				It("should return data item", func() {
					item, err := itemRepository.Add(model.Items{
						Name:        "VGA",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					result, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(uint(1)))
					Expect(result.Name).To(Equal(item.Name))
					Expect(result.CategoryID).To(Equal(item.CategoryID))
					Expect(result.Quantity).To(Equal(item.Quantity))
					Expect(result.Price).To(Equal(item.Price))
					Expect(result.Description).To(Equal(item.Description))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("get all data items table in database postgres", func() {
				It("should return all data items", func() {
					_, err := itemRepository.Add(model.Items{
						Name:        "VGA",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					_, err = itemRepository.Add(model.Items{
						Name:        "VGA2",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
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
					_, err := itemRepository.Add(model.Items{
						Name:        "VGA2",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					newItem, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())

					newItem.Name = "VGA"
					newItem.CategoryID = 2
					newItem.Quantity = 5
					newItem.Price = 4500000.00
					newItem.Description = "desc"
					_, err = itemRepository.Update(newItem)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := itemRepository.GetByItemID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.ID).To(Equal(uint(1)))
					Expect(result.Name).To(Equal(newItem.Name))
					Expect(result.CategoryID).To(Equal(newItem.CategoryID))
					Expect(result.Quantity).To(Equal(newItem.Quantity))
					Expect(result.Price).To(Equal(newItem.Price))
					Expect(result.Description).To(Equal(newItem.Description))

					err = db.Reset(connection, "items")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("delete item data by item id to items table in database postgres", func() {
				It("should soft delete data item by id", func() {
					_, err := itemRepository.Add(model.Items{
						Name:        "VGA2",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
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
					_, err := itemRepository.Add(model.Items{
						Name:        "VGA",
						CategoryID:  1,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
					})
					Expect(err).ShouldNot(HaveOccurred())

					_, err = itemRepository.Add(model.Items{
						Name:        "VGA2",
						CategoryID:  2,
						Quantity:    10,
						Price:       5000000.00,
						Description: "",
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
					_, err := activityRepository.Add(model.Activities{
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
					_, err := activityRepository.Add(model.Activities{
						ItemID:         1,
						Action:         "POST",
						QuantityChange: 5,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					_, err = activityRepository.Add(model.Activities{
						ItemID:         2,
						Action:         "POST",
						QuantityChange: -2,
						PerformedBy:    1,
					})
					Expect(err).ShouldNot(HaveOccurred())

					results, err := activityRepository.GetAll()
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
					_, err := sessionRepository.Add(model.Sessions{
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
					session, err := sessionRepository.Add(model.Sessions{
						Username:  "bangkit",
						Token:     "token",
						ExpiresAt: time.Now().Add(5 * time.Minute),
					})
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
					_, err := sessionRepository.Add(model.Sessions{
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
})
