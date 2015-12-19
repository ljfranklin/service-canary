package config_test

import (
	"os"

	configPkg "github.com/ljfranklin/service-canary/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	const configJson = `
{
  "p-mysql": [
    {
			"name": "my-mysql-db"
    }
  ]
}`

	Describe("Validate", func() {

		Context("when all env variable are present", func() {

			BeforeEach(func() {
				os.Setenv("RUN_INTERVAL", "10")
				os.Setenv("PORT", "8081")
				os.Setenv("DATADOG_API_KEY", "fakekey")
				os.Setenv("VCAP_SERVICES", configJson)
			})

			AfterEach(func() {
				os.Unsetenv("RUN_INTERVAL")
				os.Unsetenv("PORT")
				os.Unsetenv("DATADOG_API_KEY")
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
			os.Setenv("VCAP_SERVICES", configJson)

			config := configPkg.New()
			Expect(config.Services).To(Equal([]configPkg.ServiceConfig{
				configPkg.ServiceConfig{
					Name:       "my-mysql-db",
					Type:       "p-mysql",
					ConfigJSON: []byte(`{"name":"my-mysql-db"}`),
				},
			}))
		})
	})
})
