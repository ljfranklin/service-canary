package service_factory

import (
	"github.com/ljfranklin/service-canary/runner"
)

type Service interface {
	Name() string
	Type() string
	runner.Runner
}

type ServiceFactory interface {
	GetAllServices() []Service
}
