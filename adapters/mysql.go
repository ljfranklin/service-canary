package adapters

import (
	"github.com/ljfranklin/service-canary/config"
	"github.com/pivotal-golang/lager"
)

type MysqlAdapter struct {
	Adapter
	name   string
	logger lager.Logger
}

func NewMysqlAdapter(serviceName string, config *config.Config) *MysqlAdapter {
	return &MysqlAdapter{
		name:   serviceName,
		logger: config.Logger,
	}
}

func (a MysqlAdapter) Name() string {
	return a.name
}

func (a *MysqlAdapter) Setup() error {
	a.logger.Info("Setting up mysql adapter...")
	return nil
}

func (a *MysqlAdapter) Run() error {
	a.logger.Info("Running mysql adapter...")
	return nil
}
