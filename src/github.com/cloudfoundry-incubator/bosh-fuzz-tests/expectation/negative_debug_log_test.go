package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DebugLog", func() {
	var (
		negativeDebugLog Expectation
	)

	BeforeEach(func() {
		negativeDebugLog = NewNegativeDebugLog("expected-string")
	})

	Context("when debug logs contain given string", func() {
		It("returns an error", func() {
			err := negativeDebugLog.Run("expected-string")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task debug logs output contains unexpected string"))
		})
	})

	Context("when debug logs do not contain given string", func() {
		It("does not return an error", func() {
			err := negativeDebugLog.Run("some-other-string")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
