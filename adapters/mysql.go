package adapters

import "github.com/ljfranklin/service-canary/config"

type MysqlAdapter struct {
	Adapter
	name string
}

func NewMysqlAdapter(serviceName string, config *config.Config) *MysqlAdapter {
	return &MysqlAdapter{
		name: serviceName,
	}
}

func (a MysqlAdapter) Name() string {
	return a.name
}

func (a *MysqlAdapter) Run() error {
	return nil
}
