package expectation_test

import (
	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	"github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner/clirunnerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExistingInstanceDebugLog", func() {
	var (
		existingInstanceDebugLog Expectation
		cliRunner                *clirunnerfakes.FakeRunner
	)

	BeforeEach(func() {
		existingInstanceDebugLog = NewExistingInstanceDebugLog("stemcell_changed?", "etcd")
		cliRunner = &clirunnerfakes.FakeRunner{}
	})

	Context("when debug logs contain expected 'stemcell_changed?' line", func() {
		BeforeEach(func() {
			debugLog := `
			Existing desired instance 'etcd/0' in az 'z1' with active vm
			stemcell_changed? changed FROM: version: 1 TO: version: 2 on etcd/c42ab873-6f46-4273-be13-1286ba96464c (0)
			`
			cliRunner.RunWithOutputReturns(debugLog, nil)
		})

		It("does not return an error", func() {
			err := existingInstanceDebugLog.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when debug logs do not contain expected 'stemcell_changed?' line", func() {
		BeforeEach(func() {
			debugLog := `
			Existing desired instance 'etcd/0' in az 'z1' with active vm
			`
			cliRunner.RunWithOutputReturns(debugLog, nil)
		})
		It("returns an error", func() {
			err := existingInstanceDebugLog.Run(cliRunner, "1")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when debug logs have no active vm", func() {
		BeforeEach(func() {
			debugLog := `
			Existing desired instance 'etcd/0' in az 'z1' with no active vm
			`
			cliRunner.RunWithOutputReturns(debugLog, nil)
		})
		It("returns no error", func() {
			err := existingInstanceDebugLog.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when there is no existing instance", func() {
		BeforeEach(func() {
			debugLog := `
			nothing here
			`
			cliRunner.RunWithOutputReturns(debugLog, nil)
		})
		It("does not return an error", func() {
			err := existingInstanceDebugLog.Run(cliRunner, "1")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("only checks the specified instance for expectation", func() {
		debugLog := `
		Existing desired instance 'another/0'
		`
		cliRunner.RunWithOutputReturns(debugLog, nil)

		err := existingInstanceDebugLog.Run(cliRunner, "1")
		Expect(err).ToNot(HaveOccurred())
	})
})
