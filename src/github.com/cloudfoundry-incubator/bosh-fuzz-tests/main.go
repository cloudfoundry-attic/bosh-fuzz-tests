package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdeployment "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	bltenv "github.com/cloudfoundry-incubator/bosh-load-tests/environment"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		println("Usage: bft path/to/config.json seed")
		os.Exit(1)
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	testConfig := bftconfig.NewConfig(fs)
	err := testConfig.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	envConfig := bltconfig.NewConfig(fs)
	err = envConfig.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	assetsProvider := bltassets.NewProvider(envConfig.AssetsPath)

	logger.Debug("main", "Setting up environment")

	environmentProvider := bltenv.NewProvider(envConfig, fs, cmdRunner, assetsProvider)
	environment := environmentProvider.Get()

	if !envConfig.GenerateManifestOnly {
		err = environment.Setup()
		if err != nil {
			panic(err)
		}
		defer environment.Shutdown()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		environment.Shutdown()
		os.Exit(1)
	}()

	cliRunnerFactory := bltclirunner.NewFactory(envConfig.CliCmd, cmdRunner, fs)

	var directorInfo bltaction.DirectorInfo
	if envConfig.GenerateManifestOnly {
		directorInfo = bltaction.DirectorInfo{
			UUID: "blah",
			URL:  "xxx",
		}
	} else {
		directorInfo, err = bltaction.NewDirectorInfo(environment.DirectorURL(), cliRunnerFactory)
	}

	if err != nil {
		panic(err)
	}
	cliRunner := cliRunnerFactory.Create()
	cliRunner.Configure()
	defer cliRunner.Clean()

	if !envConfig.GenerateManifestOnly {
		logger.Debug("main", "Preparing to deploy")
		preparer := bftdeployment.NewPreparer(directorInfo, cliRunner, fs, assetsProvider)
		err = preparer.Prepare()
		if err != nil {
			panic(err)
		}
	}

	logger.Debug("main", "Starting deploy")
	renderer := bftdeployment.NewRenderer(fs)

	var jobsRandomizer bftdeployment.JobsRandomizer
	var networksAssigner bftdeployment.NetworksAssigner
	nameGenerator := bftdeployment.NewNameGenerator()
	ipPoolProvider := bftdeployment.NewIpPoolProvider()
	staticIpDecider := bftdeployment.NewRandomDecider()

	if len(os.Args) == 3 {
		seed, _ := strconv.ParseInt(os.Args[2], 10, 64)
		jobsRandomizer = bftdeployment.NewSeededJobsRandomizer(testConfig.Parameters, testConfig.NumberOfConsequentDeploys, seed, nameGenerator, logger)
		networksAssigner = bftdeployment.NewSeededNetworksAssigner(testConfig.Parameters.Networks, nameGenerator, ipPoolProvider, staticIpDecider, seed)
	} else {
		jobsRandomizer = bftdeployment.NewJobsRandomizer(testConfig.Parameters, testConfig.NumberOfConsequentDeploys, nameGenerator, logger)
		networksAssigner = bftdeployment.NewNetworksAssigner(testConfig.Parameters.Networks, nameGenerator, ipPoolProvider, staticIpDecider)
	}

	deployer := bftdeployment.NewDeployer(cliRunner, directorInfo, renderer, jobsRandomizer, networksAssigner, fs, envConfig.GenerateManifestOnly)
	err = deployer.RunDeploys()
	if err != nil {
		panic(err)
	}

	println("Done!")
}
