package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExistingInstanceDebugLog", func() {
	var (
		existingInstanceDebugLog Expectation
	)

	BeforeEach(func() {
		existingInstanceDebugLog = NewExistingInstanceDebugLog("stemcell_changed?")
	})

	Context("when debug logs contain expected string for existing instance", func() {
		It("does not return an error", func() {
			debugLog := `
			Existing desired instance 'foobar/0' in az 'z1'
			stemcell_changed? changed FROM: version: 1 TO: version: 2 on instance foobar/0
			`
			err := existingInstanceDebugLog.Run(debugLog)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when debug logs do not contain expected string", func() {
		It("returns an error", func() {
			debugLog := `
			Existing desired instance 'etcd/0'
			`
			err := existingInstanceDebugLog.Run(debugLog)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when there is no existing instance", func() {
		It("does not return an error", func() {
			debugLog := `
			nothing here
			`
			err := existingInstanceDebugLog.Run(debugLog)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
