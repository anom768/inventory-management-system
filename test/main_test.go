package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"inventory-management-system/app"
	"inventory-management-system/model"
)

var _ = Describe("Digital Inventory Management API", func() {
	postgres := app.NewDB()
	credential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "inventory_test",
		Port:         5432,
		Schema:       "public",
	}
	_, err := postgres.Connect(&credential)
	Expect(err).ShouldNot(HaveOccurred())
})
