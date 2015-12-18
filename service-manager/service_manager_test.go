package service_manager_test

import (
	configPkg "github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/service-factory"
	"github.com/ljfranklin/service-canary/service-factory/fakes"
	"github.com/ljfranklin/service-canary/service-manager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("ServiceManager", func() {

	It("calls Run on each Verifier in Factory", func() {

		config := &configPkg.Config{
			Logger: lagertest.NewTestLogger("ServiceManagerTest"),
		}

		fakeService := &fakes.FakeService{}
		fakeService.NameReturns("my-mysql-db")
		fakeService.TypeReturns("mysql")
		fakeService.RunReturns(nil)

		fakeServices := []service_factory.Service{
			fakeService,
		}
		factory := &fakes.FakeServiceFactory{}
		factory.GetAllServicesReturns(fakeServices)

		manager := service_manager.New(factory, config)

		err := manager.RunAllInBackground()
		Expect(err).ToNot(HaveOccurred())

		for _, service := range fakeServices {
			fake := service.(*fakes.FakeService)
			Eventually(fake.RunCallCount).Should(Equal(1), "Expected service.Run to be called at least once")
		}
	})
})
