package service_manager

import (
	"fmt"

	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/service-factory"
	"github.com/pivotal-golang/lager"
)

type serviceManager struct {
	factory service_factory.ServiceFactory
	logger  lager.Logger
}

func New(factory service_factory.ServiceFactory, config *config.Config) *serviceManager {
	return &serviceManager{
		factory: factory,
		logger:  config.Logger,
	}
}

func (m *serviceManager) RunAllInBackground() error {

	services := m.factory.GetAllServices()
	for _, service := range services {
		go func() {
			err := service.Run()
			if err != nil {
				m.logger.Error(fmt.Sprintf("Failed to run %s service with err: %s", service.Name()), err)
			}
		}()
	}

	return nil
}
