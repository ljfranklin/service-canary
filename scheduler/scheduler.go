package scheduler

import (
	"errors"
	"time"

	"github.com/ljfranklin/service-canary/runner"
)

type Scheduler interface {
	RunInBackground() error
	Stop() error
}

type scheduler struct {
	runner   runner.Runner
	interval time.Duration
	stopChan chan bool
	doneChan chan bool
}

func New(runner runner.Runner, interval time.Duration) Scheduler {
	return &scheduler{
		runner:   runner,
		interval: interval,
		stopChan: make(chan bool),
		doneChan: make(chan bool),
	}
}

func (s *scheduler) RunInBackground() error {

	// loop in background until stopChan is closed
	go func() {
		for {
			select {
			case <-s.stopChan:
				close(s.doneChan)
				return
			case <-time.After(s.interval):
				s.runner.Run()
			}
		}
	}()

	return nil
}

func (s *scheduler) Stop() error {

	close(s.stopChan)

	timeout := 1 * time.Second
	select {
	case <-s.doneChan:
		return nil
	case <-time.After(timeout):
		return errors.New("Failed to stop task after 1 second")
	}
}
