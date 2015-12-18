// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/ljfranklin/service-canary/adapters"
	"github.com/ljfranklin/service-canary/service-factory"
)

type FakeServiceFactory struct {
	GetAllServicesStub        func() ([]adapters.Adapter, error)
	getAllServicesMutex       sync.RWMutex
	getAllServicesArgsForCall []struct{}
	getAllServicesReturns struct {
		result1 []adapters.Adapter
		result2 error
	}
}

func (fake *FakeServiceFactory) GetAllServices() ([]adapters.Adapter, error) {
	fake.getAllServicesMutex.Lock()
	fake.getAllServicesArgsForCall = append(fake.getAllServicesArgsForCall, struct{}{})
	fake.getAllServicesMutex.Unlock()
	if fake.GetAllServicesStub != nil {
		return fake.GetAllServicesStub()
	} else {
		return fake.getAllServicesReturns.result1, fake.getAllServicesReturns.result2
	}
}

func (fake *FakeServiceFactory) GetAllServicesCallCount() int {
	fake.getAllServicesMutex.RLock()
	defer fake.getAllServicesMutex.RUnlock()
	return len(fake.getAllServicesArgsForCall)
}

func (fake *FakeServiceFactory) GetAllServicesReturns(result1 []adapters.Adapter, result2 error) {
	fake.GetAllServicesStub = nil
	fake.getAllServicesReturns = struct {
		result1 []adapters.Adapter
		result2 error
	}{result1, result2}
}

var _ service_factory.ServiceFactory = new(FakeServiceFactory)
