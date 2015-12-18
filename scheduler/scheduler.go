package scheduler

import (
	"errors"
	"time"

	"github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/runner"
	"github.com/pivotal-golang/lager"
)

type Scheduler interface {
	RunInBackground() error
	Stop() error
}

type scheduler struct {
	runner   runner.Runner
	interval time.Duration
	logger   lager.Logger
	stopChan chan bool
	doneChan chan bool
}

func New(runner runner.Runner, config *config.Config) Scheduler {
	return &scheduler{
		runner:   runner,
		interval: config.Interval,
		logger:   config.Logger,
		stopChan: make(chan bool),
		doneChan: make(chan bool),
	}
}

func (s *scheduler) RunInBackground() error {
	s.logger.Info("Running in background...")

	// loop in background until stopChan is closed
	go func() {
		for {
			select {
			case <-s.stopChan:
				close(s.doneChan)
				return
			case <-time.After(s.interval):
				err := s.runner.Run()
				if err != nil {
					s.logger.Error("Error running command", err)
				} else {
					s.logger.Info("Successfully ran command")
				}
			}
		}
	}()

	return nil
}

func (s *scheduler) Stop() error {
	s.logger.Info("Stopping scheduler...")

	close(s.stopChan)

	timeout := 1 * time.Second
	select {
	case <-s.doneChan:
		s.logger.Info("Successfully stopped scheduler")
		return nil
	case <-time.After(timeout):
		err := errors.New("Failed to stop task after 1 second")
		s.logger.Error("Failed to stop scheduler", err)
		return err
	}
}
