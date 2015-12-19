package service_manager

import (
	"fmt"

	"github.com/ljfranklin/service-canary/adapters"
	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/event-emitter"
	"github.com/ljfranklin/service-canary/service-factory"
	"github.com/pivotal-golang/lager"
)

type serviceManager struct {
	factory  service_factory.ServiceFactory
	emitter  event_emitter.Emitter
	logger   lager.Logger
	services []adapters.Adapter
}

func New(factory service_factory.ServiceFactory, emitter event_emitter.Emitter, config *config.Config) *serviceManager {
	return &serviceManager{
		factory: factory,
		emitter: emitter,
		logger:  config.Logger,
	}
}

func (m *serviceManager) Setup() error {
	var err error
	m.services, err = m.factory.GetAllServices()

	if err != nil {
		return fmt.Errorf("Failed to Setup Factory: %s", err.Error())
	}

	for _, service := range m.services {
		err := service.Setup()
		if err != nil {
			return fmt.Errorf("Failed to setup '%s': %s", service.Name(), err.Error())
		}
	}
	return nil
}

func (m *serviceManager) RunAllInBackground() error {
	for _, service := range m.services {
		tags := map[string]string{}
		go func() {
			err := service.Run()
			result := 1
			if err != nil {
				m.logger.Error(fmt.Sprintf("Failed to run %s service", service.Name()), err)
				result = 0
			}

			err = m.emitter.Emit(service.Name(), result, tags)
			if err != nil {
				m.logger.Error("Failed to emit event to datadog", err)
			}
		}()
	}

	return nil
}

// ensure this conforms the Runner interface
func (m *serviceManager) Run() error {
	return m.RunAllInBackground()
}
