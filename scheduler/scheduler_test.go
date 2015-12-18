package scheduler_test

import (
	"time"

	configPkg "github.com/ljfranklin/service-canary/config"
	"github.com/ljfranklin/service-canary/runner/fakes"
	schedulerPkg "github.com/ljfranklin/service-canary/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Scheduler", func() {

	Describe("#Run", func() {

		It("calls the command at the given interval", func() {

			config := &configPkg.Config{
				Interval: 10 * time.Millisecond,
				Logger:   lagertest.NewTestLogger("SchedulerTest"),
			}

			runner := &fakes.FakeRunner{}
			scheduler := schedulerPkg.New(runner, config)

			err := scheduler.RunInBackground()
			Expect(err).ToNot(HaveOccurred())

			Eventually(runner.RunCallCount).Should(BeNumerically(">", 1), "Expected Run to be called more than once")

			err = scheduler.Stop()
			Expect(err).ToNot(HaveOccurred())

			currCallCount := runner.RunCallCount()
			Consistently(runner.RunCallCount).Should(Equal(currCallCount), "Expected Run to not be called after Stop")
		})
	})
})
