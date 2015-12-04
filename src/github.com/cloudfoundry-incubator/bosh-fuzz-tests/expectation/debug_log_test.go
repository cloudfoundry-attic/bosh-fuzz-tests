package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DebugLog", func() {
	var (
		debugLog Expectation
	)

	BeforeEach(func() {
		debugLog = NewDebugLog("expected-string")
	})

	Context("when debug logs contain expected string", func() {
		It("does not return an error", func() {
			err := debugLog.Run("expected-string")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when debug logs do not contain expected string", func() {
		It("returns an error", func() {
			err := debugLog.Run("nothing here")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task debug logs output does not contain expected string"))
		})
	})
})
