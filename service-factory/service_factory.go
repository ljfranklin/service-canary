package service_factory

import (
	"fmt"

	"github.com/ljfranklin/service-canary/adapters"
	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"
)

type ServiceFactory interface {
	GetAllServices() []adapters.Adapter
}

type serviceFactory struct {
	config *config.Config
	logger lager.Logger
}

func New(config *config.Config) ServiceFactory {
	return &serviceFactory{
		config: config,
		logger: config.Logger,
	}
}

func (f *serviceFactory) GetAllServices() []adapters.Adapter {

	services := []adapters.Adapter{}
	for _, serviceConfig := range f.config.Services {
		switch serviceConfig.Type {
		case "mysql":
			adapter := adapters.NewMysqlAdapter(serviceConfig.Name, f.config)
			services = append(services, adapter)
		default:
			err := fmt.Errorf("Unknown service type '%s'", serviceConfig.Type)
			f.logger.Error("Failed to build service adapter", err)
		}
	}

	return services
}
