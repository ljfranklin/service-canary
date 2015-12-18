package config_test

import (
	"os"

	configPkg "github.com/ljfranklin/service-canary/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("Validate", func() {

		Context("when all env variable are present", func() {

			BeforeEach(func() {
				os.Setenv("RUN_INTERVAL", "10")
				os.Setenv("VCAP_SERVICES", `
{
  "p-mysql": [
    {
			"name": "my-mysql-db"
    }
  ]
}`)
			})

			AfterEach(func() {
				os.Unsetenv("RUN_INTERVAL")
				os.Unsetenv("VCAP_SERVICES")
			})

			It("returns nil", func() {
				config := configPkg.New()
				err := config.Validate()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when env variables are absent", func() {

			BeforeEach(func() {
				os.Unsetenv("RUN_INTERVAL")
			})

			It("returns an error", func() {
				config := configPkg.New()
				err := config.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RUN_INTERVAL"))
			})
		})
	})

	Describe("Services", func() {
		It("parses VCAP_SERVICES into structs", func() {
			os.Setenv("VCAP_SERVICES", `
{
  "p-mysql": [
    {
			"name": "my-mysql-db"
    }
  ]
}`)

			config := configPkg.New()
			Expect(config.Services).To(Equal([]configPkg.ServiceConfig{
				configPkg.ServiceConfig{
					Name: "my-mysql-db",
					Type: "p-mysql",
				},
			}))
		})
	})
})
