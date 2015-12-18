package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/pivotal-golang/lager"
)

type Config struct {
	Interval time.Duration
	Logger   lager.Logger
	Services []ServiceConfig
}

type ServiceConfig struct {
	Name string
	Type string
}

type serviceJson struct {
	Name string `json:"name"`
}

func New() *Config {

	cnf := &Config{}

	cnf.Logger = lager.NewLogger("ServiceCanary")

	// conversion errors are defered to Validate
	runInterval, _ := strconv.Atoi(os.Getenv("RUN_INTERVAL"))
	cnf.Interval = time.Duration(runInterval) * time.Second

	cnf.Services = []ServiceConfig{}

	// see tests for example format
	servicesJson := []byte(os.Getenv("VCAP_SERVICES"))
	var rawServices map[string][]map[string]interface{}
	if err := json.Unmarshal(servicesJson, &rawServices); err != nil {
		cnf.Logger.Error("Failed to Parse VCAP_SERVICES", err)
	} else {
		for k, v := range rawServices {
			serviceType := k
			servicesInType := v

			for _, serviceProperties := range servicesInType {
				newService := ServiceConfig{
					Name: serviceProperties["name"].(string),
					Type: serviceType,
				}
				cnf.Services = append(cnf.Services, newService)
			}
		}
	}

	return cnf
}

func (c Config) Validate() error {
	if c.Interval == 0 {
		return errors.New("RUN_INTERVAL cannot be zero")
	}

	if len(c.Services) == 0 {
		return errors.New("Failed to parse any services from VCAP_SERVICES")
	}

	return nil
}
