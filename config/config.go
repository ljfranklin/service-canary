package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotal-golang/lager"
)

type Config struct {
	Interval time.Duration
	Port     int
	Logger   lager.Logger
	Services []ServiceConfig
	setupErr error
}

type ServiceConfig struct {
	Name       string
	Type       string
	ConfigJSON []byte
}

type serviceJson struct {
	Name string `json:"name"`
}

func New() *Config {

	cnf := &Config{}

	cf_lager.AddFlags(flag.CommandLine)
	flag.Parse()

	cnf.Logger, _ = cf_lager.New("ServiceCanary")

	// conversion errors are defered to Validate
	runInterval, _ := strconv.Atoi(os.Getenv("RUN_INTERVAL"))
	cnf.Interval = time.Duration(runInterval) * time.Second

	cnf.Port, _ = strconv.Atoi(os.Getenv("PORT"))

	cnf.Services = []ServiceConfig{}

	// see tests for example format
	servicesJson := []byte(os.Getenv("VCAP_SERVICES"))
	var rawServices map[string][]map[string]interface{}
	if err := json.Unmarshal(servicesJson, &rawServices); err != nil {
		cnf.setupErr = fmt.Errorf("Failed to Parse VCAP_SERVICES: %s", err.Error())
	} else {
		for k, v := range rawServices {
			serviceType := k
			servicesInType := v

			for _, serviceProperties := range servicesInType {
				configJson, _ := json.Marshal(serviceProperties)
				newService := ServiceConfig{
					Name:       serviceProperties["name"].(string),
					Type:       serviceType,
					ConfigJSON: configJson,
				}
				cnf.Services = append(cnf.Services, newService)
			}
		}
	}

	return cnf
}

func (c Config) Validate() error {

	errMsg := ""
	if c.setupErr != nil {
		errMsg += fmt.Sprintf("Failed to build config: %s\n", c.setupErr.Error())
	}

	if c.Interval == 0 {
		errMsg += "RUN_INTERVAL was not present in environment\n"
	}

	if c.Port == 0 {
		errMsg += "PORT was not present in environment\n"
	}

	if len(c.Services) == 0 {
		errMsg += "Failed to parse any services from VCAP_SERVICES\n"
	}

	if len(errMsg) > 0 {
		return errors.New(errMsg)
	}
	return nil
}
