package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DebugLog", func() {
	var (
		negativeDebugLog Expectation
		cliRunner        *clirunnerfakes.FakeRunner
	)

	BeforeEach(func() {
		negativeDebugLog = NewNegativeDebugLog("expected-string")
		cliRunner = &clirunnerfakes.FakeRunner{}
	})

	Context("when debug logs contain given string", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns("expected-string", nil)
		})
		It("returns an error", func() {
			err := negativeDebugLog.Run(cliRunner, "1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task debug logs output contains unexpected string"))
		})
	})

	Context("when debug logs do not contain given string", func() {
		BeforeEach(func() {
			cliRunner.RunWithOutputReturns("some-other-string", nil)
		})
		It("does not return an error", func() {
			err := negativeDebugLog.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
