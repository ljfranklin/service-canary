package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/pivotal-golang/lager"
)

type Config struct {
	Interval time.Duration
	Logger   lager.Logger
}

func New() *Config {

	cnf := &Config{}
	// conversion errors are defered to Validate
	runInterval, _ := strconv.Atoi(os.Getenv("RUN_INTERVAL"))

	cnf.Interval = time.Duration(runInterval) * time.Second

	cnf.Logger = lager.NewLogger("ServiceCanary")

	return cnf
}

func (c Config) Validate() error {
	if c.Interval == 0 {
		return errors.New("RUN_INTERVAL cannot be zero")
	}

	return nil
}
