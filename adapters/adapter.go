package adapters

import "github.com/ljfranklin/service-canary/runner"

type Adapter interface {
	Name() string
	runner.Runner
}
