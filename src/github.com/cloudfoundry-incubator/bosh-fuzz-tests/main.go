package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftdeployment "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"

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

	var seed int64
	if len(os.Args) == 3 {
		seed, _ = strconv.ParseInt(os.Args[2], 10, 64)
	} else {
		seed = time.Now().Unix()
	}

	logger.Info("main", "Seeding with %d", seed)
	rand.Seed(seed)

	nameGenerator := bftnamegen.NewNameGenerator()
	decider := bftdecider.NewRandomDecider()

	ipPoolProvider := bftdeployment.NewIpPoolProvider()
	parameterProvider := bftparam.NewParameterProvider(testConfig.Parameters, nameGenerator, decider, logger)
	inputGenerator := bftdeployment.NewInputGenerator(testConfig.Parameters, parameterProvider, testConfig.NumberOfConsequentDeploys, nameGenerator, logger)
	networksAssigner := bftdeployment.NewNetworksAssigner(testConfig.Parameters.Networks, nameGenerator, ipPoolProvider, decider)

	deployer := bftdeployment.NewDeployer(cliRunner, directorInfo, renderer, inputGenerator, networksAssigner, fs, envConfig.GenerateManifestOnly)
	err = deployer.RunDeploys()
	if err != nil {
		panic(err)
	}

	println("Done!")
}
