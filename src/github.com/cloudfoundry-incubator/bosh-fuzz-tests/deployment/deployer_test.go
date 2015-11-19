package deployment_test

import (
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"

	. "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

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

		fs.ReturnTempFile = fakesys.NewFakeFile("cli-config-path", fs)

		cliRunner := bltclirunner.NewRunner(boshCmd, cmdRunner, fs)
		cliRunner.Configure()
		renderer := NewRenderer(fs)

		parameters := bftconfig.Parameters{
			NameLength:               []int{5, 10},
			Instances:                []int{2, 4},
			AvailabilityZones:        [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition: []string{"disk_type", "disk_pool"},
			PersistentDiskSize:       []int{0, 100},
			NumberOfJobs:             []int{1, 2},
			MigratedFromCount:        []int{0},
		}

		logger := boshlog.NewLogger(boshlog.LevelNone)
		nameGenerator := NewNameGenerator()
		jobsRandomizer := NewSeededJobsRandomizer(parameters, 2, 64, nameGenerator, logger)
		networksAssigner := NewSeededNetworksAssigner(5)
		deployer = NewDeployer(cliRunner, directorInfo, renderer, jobsRandomizer, networksAssigner, fs)
	})

	It("runs deploys with generated manifests", func() {
		fs.ReturnTempFile = fakesys.NewFakeFile("manifest-path", fs)

		err := deployer.RunDeploys()
		Expect(err).ToNot(HaveOccurred())
		Expect(fs.FileExists("config-path")).To(BeFalse())

		Expect(cmdRunner.RunComplexCommands).To(ConsistOf([]boshsys.Command{
			{
				Name: "bosh",
				Args: []string{"-n", "-c", "cli-config-path", "update", "cloud-config", "manifest-path"},
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
				Args: []string{"-n", "-c", "cli-config-path", "update", "cloud-config", "manifest-path"},
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
