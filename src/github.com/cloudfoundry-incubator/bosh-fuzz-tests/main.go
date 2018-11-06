package main

import (
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bftanalyzer "github.com/cloudfoundry-incubator/bosh-fuzz-tests/analyzer"
	bftconfig "github.com/cloudfoundry-incubator/bosh-fuzz-tests/config"
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftdeployment "github.com/cloudfoundry-incubator/bosh-fuzz-tests/deployment"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	bftparam "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter"
	bftnetwork "github.com/cloudfoundry-incubator/bosh-fuzz-tests/parameter/network"
	bftvariables "github.com/cloudfoundry-incubator/bosh-fuzz-tests/variables"

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

	defer func() {
		cmd := exec.Command("bash", "-c", `
		ps aux \
		| egrep "nats-server|bin/bosh-director|bosh-config-server-executable|bosh-agent" \
		| grep -v grep \
		| awk '{print $2}' \
		| xargs kill
		`)
		cmd.Run()
	}()

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

	setup(logger)
	envConfig.CpiConfig.PreferredCpiAPIVersion = testConfig.Parameters.CpiAPIVersion[rand.Intn(len(testConfig.Parameters.CpiAPIVersion))]

	assetsProvider := bltassets.NewProvider(envConfig.AssetsPath)

	logger.Debug("main", "Setting up environment")

	environmentProvider := bltenv.NewProvider(envConfig, fs, cmdRunner, assetsProvider, logger)
	environment := environmentProvider.Get()

	if !testConfig.GenerateManifestOnly {
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
	if testConfig.GenerateManifestOnly {
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

	cliRunner, err := cliRunnerFactory.Create("bosh")
	if err != nil {
		panic(err)
	}

	uaaRunner, err := cliRunnerFactory.Create("uaac")
	if err != nil {
		panic(err)
	}

	if !testConfig.GenerateManifestOnly {
		logger.Debug("main", "Preparing to deploy")
		preparer := bftdeployment.NewPreparer(directorInfo, cliRunner, fs, assetsProvider)
		err = preparer.Prepare()
		if err != nil {
			panic(err)
		}
	}

	logger.Debug("main", "Starting deploy")
	renderer := bftdeployment.NewRenderer(fs)

	nameGenerator := bftnamegen.NewNameGenerator()
	decider := bftdecider.NewRandomDecider()

	ipPoolProvider := bftnetwork.NewIpPoolProvider()
	networksAssigner := bftnetwork.NewAssigner(testConfig.Parameters.Networks, nameGenerator, ipPoolProvider, decider, logger)
	parameterProvider := bftparam.NewParameterProvider(testConfig.Parameters, nameGenerator, decider, networksAssigner, logger)
	inputGenerator := bftdeployment.NewInputGenerator(testConfig.Parameters, parameterProvider, testConfig.NumberOfConsequentDeploys, nameGenerator, decider, logger)
	analyzer := bftanalyzer.NewAnalyzer(logger)
	errandGenerator := bftdeployment.NewErrandStepGenerator()

	varsRandomizer := bftvariables.DefaultNumberRandomizer{}
	varsPathBuilder := bftvariables.NewPathBuilder()
	varsPathPicker := bftvariables.NewPathPicker(varsRandomizer)
	varsPlaceholderPlanter := bftvariables.NewPlaceholderPlanter(nameGenerator)

	sprinkler := bftvariables.NewSprinkler(
		testConfig.Parameters,
		fs,
		varsRandomizer,
		varsPathBuilder,
		varsPathPicker,
		varsPlaceholderPlanter,
		nameGenerator,
	)

	deployer := bftdeployment.NewDeployer(
		cliRunner,
		uaaRunner,
		directorInfo,
		renderer,
		inputGenerator,
		[]bftdeployment.StepGenerator{errandGenerator},
		analyzer,
		sprinkler,
		fs,
		logger,
		testConfig.GenerateManifestOnly,
	)
	err = deployer.RunDeploys()
	if err != nil {
		panic(err)
	}

	println("Done!")
}

func setup(logger boshlog.Logger) {
	var seed int64
	if len(os.Args) == 3 {
		seed, _ = strconv.ParseInt(os.Args[2], 10, 64)
	} else {
		seed = time.Now().Unix()
	}

	logger.Info("main", "Seeding with %d", seed)
	rand.Seed(seed)
}
