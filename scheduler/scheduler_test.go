package scheduler_test

import (
	"time"

	"github.com/ljfranklin/service-canary/runner/fakes"
	schedulerPkg "github.com/ljfranklin/service-canary/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scheduler", func() {

	Describe("#Run", func() {

		It("calls the command at the given interval", func() {

			runner := &fakes.FakeRunner{}
			interval := 10 * time.Millisecond
			scheduler := schedulerPkg.New(runner, interval)

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
