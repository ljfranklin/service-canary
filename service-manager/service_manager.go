package service_manager

import (
	"fmt"

	"github.com/ljfranklin/service-canary/adapters"
	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/service-factory"
	"github.com/pivotal-golang/lager"
)

type serviceManager struct {
	factory  service_factory.ServiceFactory
	logger   lager.Logger
	services []adapters.Adapter
}

func New(factory service_factory.ServiceFactory, config *config.Config) *serviceManager {
	return &serviceManager{
		factory: factory,
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
		go func() {
			err := service.Run()
			if err != nil {
				m.logger.Error(fmt.Sprintf("Failed to run %s service with err: %s", service.Name()), err)
			}
		}()
	}

	return nil
}

// ensure this conforms the Runner interface
func (m *serviceManager) Run() error {
	return m.RunAllInBackground()
}
