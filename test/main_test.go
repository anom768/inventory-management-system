package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"inventory-management-system/app"
	"inventory-management-system/model"
	"inventory-management-system/repository"
)

var _ = Describe("Digital Inventory Management API", func() {
	var userRepository repository.UserRepository

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

	BeforeEach(func() {
		//err = connection.Migrator().DropTable("users")
		Expect(err).ShouldNot(HaveOccurred())

		err := connection.AutoMigrate(&model.Users{})
		Expect(err).ShouldNot(HaveOccurred())

		err = db.Reset(connection, "users")
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Repository", func() {
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
	})
})
