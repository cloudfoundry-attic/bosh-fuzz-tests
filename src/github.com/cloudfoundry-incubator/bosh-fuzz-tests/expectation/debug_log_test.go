package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DebugLog", func() {
	var (
		debugLog  Expectation
		cliRunner *clirunnerfakes.FakeRunner
	)

	BeforeEach(func() {
		debugLog = NewDebugLog("expected-string")
		cliRunner = &clirunnerfakes.FakeRunner{}
	})

	Context("when debug logs contain expected string", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns("expected-string", nil)
		})

		It("does not return an error", func() {
			err := debugLog.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when debug logs do not contain expected string", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns("nothing-here", nil)
		})

		It("returns an error", func() {
			err := debugLog.Run(cliRunner, "1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task debug logs output does not contain expected string"))
		})
	})
})
