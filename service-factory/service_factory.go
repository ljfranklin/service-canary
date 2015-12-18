package service_factory

import (
	"fmt"

	"github.com/ljfranklin/service-canary/adapters"
	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"
)

type ServiceFactory interface {
	GetAllServices() ([]adapters.Adapter, error)
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

func (f *serviceFactory) GetAllServices() ([]adapters.Adapter, error) {

	services := []adapters.Adapter{}
	for _, serviceConfig := range f.config.Services {
		switch serviceConfig.Type {
		case "p-mysql":
			adapter := adapters.NewMysqlAdapter(serviceConfig.Name, f.config)
			services = append(services, adapter)
		default:
			err := fmt.Errorf("Unknown service type '%s'", serviceConfig.Type)
			return nil, err
		}
	}

	return services, nil
}
