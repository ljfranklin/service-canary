package service_factory_test

import (
	"github.com/ljfranklin/service-canary/adapters"
	configPkg "github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/service-factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("ServiceManager", func() {

	Describe("GetAllServices", func() {

		It("builds a service for each service specified in config", func() {
			config := &configPkg.Config{
				Logger: lagertest.NewTestLogger("ServiceFactoryTest"),
				Services: []configPkg.ServiceConfig{
					configPkg.ServiceConfig{
						Name: "my-test-db",
						Type: "mysql",
					},
				},
			}

			factory := service_factory.New(config)
			factoryServices, err := factory.GetAllServices()
			Expect(err).ToNot(HaveOccurred())

			Expect(factoryServices).To(HaveLen(len(config.Services)))

			mysqlService := factoryServices[0]
			Expect(mysqlService.Name()).To(Equal("my-test-db"))
			Expect(mysqlService).To(BeAssignableToTypeOf(&adapters.MysqlAdapter{}))
		})
	})
})
