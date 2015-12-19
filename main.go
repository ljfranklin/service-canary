package main

import (
	"fmt"
	"net/http"

	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/event-emitter"
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

	emitter := event_emitter.NewDataDogEmitter(rootConfig)

	manager := service_manager.New(factory, emitter, rootConfig)

	scheduler := scheduler.New(manager, rootConfig)

	err = scheduler.RunInBackground()
	if err != nil {
		rootConfig.Logger.Fatal("Failed to run scheduler", err)
	}

	// bind to PORT so CF sees the app as healthy
	http.ListenAndServe(fmt.Sprintf(":%d", rootConfig.Port), nil)
}
