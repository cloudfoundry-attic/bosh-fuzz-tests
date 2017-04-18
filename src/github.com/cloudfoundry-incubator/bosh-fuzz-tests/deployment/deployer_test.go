package deployment_test

import (
	"math/rand"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	fakebftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider/fakes"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	bftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"
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
			Name: "fake-director",
			UUID: "fake-director-uuid",
			URL:  "fake-director-url",
		}
		boshCmd := boshsys.Command{Name: "bosh"}

		cliRunner := bltclirunner.NewRunner(boshCmd, cmdRunner, fs)
		renderer := NewRenderer(fs)

		cmdRunner.AddCmdResult("bosh -e fake-director-url --ca-cert /tmp/cert -n --tty --client admin --client-secret admin -d foo-deployment deploy manifest-path", fakesys.FakeCmdResult{
			Stdout: `Task 888`,
		})

		cmdRunner.AddCmdResult("bosh -e fake-director-url --ca-cert /tmp/cert -n --tty --client admin --client-secret admin -d foo-deployment deploy manifest-path", fakesys.FakeCmdResult{
			Stdout: `Task 999`,
		})

		parameters := bftconfig.Parameters{
			NameLength:                 []int{5, 10},
			Instances:                  []int{2, 4},
			AvailabilityZones:          [][]string{[]string{"z1"}, []string{"z1", "z2"}},
			PersistentDiskDefinition:   []string{"disk_type", "disk_pool"},
			PersistentDiskSize:         []int{0, 100},
			NumberOfJobs:               []int{1, 2},
			MigratedFromCount:          []int{0},
			VmTypeDefinition:           []string{"vm_type"},
			StemcellDefinition:         []string{"name"},
			StemcellVersions:           []string{"1"},
			Templates:                  [][]string{[]string{"simple"}},
			NumberOfCompilationWorkers: []int{3},
			Canaries:                   []int{3},
			MaxInFlight:                []int{3},
			Serial:                     []string{"true"},
			NumOfCloudProperties:       []int{2},
		}

		networks := [][]string{[]string{"manual_with_static", "manual_with_auto"}}

		logger := boshlog.NewLogger(boshlog.LevelNone)
		rand.Seed(64)

		nameGenerator := bftnamegen.NewNameGenerator()
		decider := &fakebftdecider.FakeDecider{}

		ipPoolProvider := bftnetwork.NewIpPoolProvider()
		networkAssigner := bftnetwork.NewAssigner(networks, nameGenerator, ipPoolProvider, decider, logger)
		parameterProvider := bftparam.NewParameterProvider(parameters, nameGenerator, decider, networkAssigner, logger)
		inputGenerator := NewInputGenerator(parameters, parameterProvider, 2, nameGenerator, decider, logger)
		analyzer := bftanalyzer.NewAnalyzer(logger)
		deployer = NewDeployer(cliRunner, directorInfo, renderer, inputGenerator, analyzer, fs, logger, false)
	})

	It("runs deploys with generated manifests", func() {
		fs.ReturnTempFile = fakesys.NewFakeFile("manifest-path", fs)

		err := deployer.RunDeploys()
		Expect(err).ToNot(HaveOccurred())
		Expect(fs.FileExists("config-path")).To(BeFalse())

		Expect(cmdRunner.RunComplexCommands).To(ConsistOf([]boshsys.Command{
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "update-cloud-config", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "-d", "foo-deployment", "deploy", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "task", "888", "--debug"},
			},
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "update-cloud-config", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "-d", "foo-deployment", "deploy", "manifest-path"},
			},
			{
				Name: "bosh",
				Args: []string{"-e", "fake-director-url", "--ca-cert", "/tmp/cert", "-n", "--tty", "--client", "admin", "--client-secret", "admin", "task", "999", "--debug"},
			},
		}))
	})
})
