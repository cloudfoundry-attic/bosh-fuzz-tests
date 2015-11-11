package manifest_test

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/manifest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployer", func() {
	var (
		cmdRunner *fakesys.FakeCmdRunner
		fs        *fakesys.FakeFileSystem
		deployer  Deployer
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
		cmdRunner = fakesys.NewFakeCmdRunner()
		directorInfo := bltaction.DirectorInfo{
			UUID: "fake-director-uuid",
			URL:  "fake-director-url",
		}
		boshCmd := boshsys.Command{Name: "bosh"}

		manifestFile := fakesys.NewFakeFile("cli-config-path", fs)
		fs.ReturnTempFile = manifestFile

		cliRunner := bltclirunner.NewRunner(boshCmd, cmdRunner, fs)
		cliRunner.Configure()
		renderer := NewRenderer(fs)

		parameters := bftconfig.Parameters{
			Name:              []string{"foo", "bar"},
			Instances:         []int{2, 4},
			AvailabilityZones: [][]string{[]string{"z1"}, []string{"z1", "z2"}},
		}

		inputRandomizer := NewSeededInputRandomizer(parameters, 2, 64)
		deployer = NewDeployer(cliRunner, directorInfo, renderer, inputRandomizer, fs)
	})

	It("runs deploys with generated manifests", func() {
		manifestFile := fakesys.NewFakeFile("manifest-path", fs)
		fs.ReturnTempFile = manifestFile

		err := deployer.RunDeploys()
		Expect(err).ToNot(HaveOccurred())
		Expect(fs.FileExists("manifest-path")).To(BeFalse())

		Expect(cmdRunner.RunComplexCommands).To(ConsistOf([]boshsys.Command{
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "target", "fake-director-url"},
			},
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "login", "admin", "admin"},
			},
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "deployment", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "deploy"},
			},
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "deployment", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "deploy"},
			},
		}))
	})
})
