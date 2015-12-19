package service_manager_test

import (
	"github.com/ljfranklin/service-canary/adapters"
	adapterFakes "github.com/ljfranklin/service-canary/adapters/fakes"
	configPkg "github.com/ljfranklin/service-canary/config"
	emitterFakes "github.com/ljfranklin/service-canary/event-emitter/fakes"
	factoryFakes "github.com/ljfranklin/service-canary/service-factory/fakes"
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

		fakeService := &adapterFakes.FakeAdapter{}
		fakeService.NameReturns("my-mysql-db")
		fakeService.RunReturns(nil)

		fakeServices := []adapters.Adapter{
			fakeService,
		}
		factory := &factoryFakes.FakeServiceFactory{}
		factory.GetAllServicesReturns(fakeServices, nil)

		emitter := &emitterFakes.FakeEmitter{}

		manager := service_manager.New(factory, emitter, config)

		err := manager.Setup()
		Expect(err).ToNot(HaveOccurred())
		for _, service := range fakeServices {
			fake := service.(*adapterFakes.FakeAdapter)
			Expect(fake.SetupCallCount()).To(Equal(1), "Expected service.Setup to be called once")
		}

		err = manager.RunAllInBackground()
		Expect(err).ToNot(HaveOccurred())

		for _, service := range fakeServices {
			fake := service.(*adapterFakes.FakeAdapter)
			Eventually(fake.RunCallCount).Should(Equal(1), "Expected service.Run to be called once")
		}
		Expect(emitter.EmitCallCount()).To(Equal(len(fakeServices)), "Expected emitter.Emit to be called once for each service")
	})
})
