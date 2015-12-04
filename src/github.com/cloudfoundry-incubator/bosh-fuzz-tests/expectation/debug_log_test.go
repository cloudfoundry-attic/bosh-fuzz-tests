package expectation_test

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DebugLog", func() {
	var (
		cmdRunner *fakesys.FakeCmdRunner
		debugLog  Expectation
	)

	BeforeEach(func() {
		fs := fakesys.NewFakeFileSystem()
		cmdRunner = fakesys.NewFakeCmdRunner()
		boshCmd := boshsys.Command{Name: "bosh"}
		fs.ReturnTempFile = fakesys.NewFakeFile("cli-config-path", fs)

		cliRunner := bltclirunner.NewRunner(boshCmd, cmdRunner, fs)
		cliRunner.Configure()

		debugLog = NewDebugLog("expected-string", cliRunner)
	})

	Context("when debug logs contain expected string", func() {
		BeforeEach(func() {
			cmdRunner.AddCmdResult("bosh -n -c cli-config-path task 15 --debug", fakesys.FakeCmdResult{
				Stdout: "expected-string",
			})
		})

		It("does not return an error", func() {
			err := debugLog.Run("15")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when debug logs do not contain expected string", func() {
		BeforeEach(func() {
			cmdRunner.AddCmdResult("bosh -n -c cli-config-path task 15 --debug", fakesys.FakeCmdResult{
				Stdout: "nothing here",
			})
		})

		It("returns an error", func() {
			err := debugLog.Run("15")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Task debug logs output does not contain expected string"))
		})
	})
})
