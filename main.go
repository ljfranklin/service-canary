package main

import (
	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/scheduler"
	"github.com/ljfranklin/service-canary/service-factory"
	"github.com/ljfranklin/service-canary/service-manager"
)

func main() {
	rootConfig := config.New()
	err := rootConfig.Validate()
	if err != nil {
		rootConfig.Logger.Fatal("Failed to validate config", err)
	}

	factory := service_factory.New(rootConfig)

	manager := service_manager.New(factory, rootConfig)

	scheduler := scheduler.New(manager, rootConfig)

	errChan := make(chan error)
	go func() {
		errChan <- scheduler.RunInBackground()
	}()

	err = <-errChan
	rootConfig.Logger.Fatal("App exited unexpectedly", err)
}
