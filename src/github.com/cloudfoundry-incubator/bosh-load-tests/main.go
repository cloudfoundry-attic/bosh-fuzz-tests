package main

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"os"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	bltenv "github.com/cloudfoundry-incubator/bosh-load-tests/environment"
	bltflow "github.com/cloudfoundry-incubator/bosh-load-tests/flow"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		println("Usage: blt path/to/config.json [path/to/state.json]")
		os.Exit(1)
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	config := bltconfig.NewConfig(fs)
	err := config.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	assetsProvider := bltassets.NewProvider(config.AssetsPath)

	logger.Debug("main", "Setting up environment")
	environmentProvider := bltenv.NewProvider(config, fs, cmdRunner, assetsProvider, logger)
	environment := environmentProvider.Get()

	logger.Debug("main", "Starting deploy")

	cliRunnerFactory := bltclirunner.NewFactory(config.CliCmd, cmdRunner, fs)

	var directorInfo bltaction.DirectorInfo
	if environment == nil {
		directorInfo, err = bltaction.NewDirectorInfo(os.Getenv("BOSH_ENVIRONMENT"), cliRunnerFactory)
	} else {
		directorInfo, err = bltaction.NewDirectorInfo(environment.DirectorURL(), cliRunnerFactory)
	}

	if err != nil {
		panic(err)
	}

	actionFactory := bltaction.NewFactory(directorInfo, fs, assetsProvider)

	prepareConfigServerFlow := bltflow.NewFlow(1, []bltflow.ActionInfo{{Name: "prepare_config_server"}}, actionFactory, cliRunnerFactory)
	err = prepareConfigServerFlow.Run(false)
	if err != nil {
		panic(err)
	}

	prepareActionFlow := bltflow.NewFlow(1, []bltflow.ActionInfo{{Name: "prepare"}}, actionFactory, cliRunnerFactory)
	err = prepareActionFlow.Run(false)
	if err != nil {
		panic(err)
	}

	if !config.UsingLegacyManifest {
		uploadCloudConfigActionFlow := bltflow.NewFlow(2, []bltflow.ActionInfo{{Name: "upload_cloud_config"}}, actionFactory, cliRunnerFactory)
		err := uploadCloudConfigActionFlow.Run(false)
		if err != nil {
			panic(err)
		}
	}

	randomizer := bltflow.NewRandomizer(actionFactory, cliRunnerFactory, fs, logger)
	if len(os.Args) == 3 {
		err := randomizer.Configure(os.Args[2])
		if err != nil {
			panic(err)
		}
	} else {
		err := randomizer.Prepare(config.Flows, config.NumberOfDeployments)
		if err != nil {
			panic(err)
		}
	}

	doneCh := make(chan error)
	for i := 0; i < config.NumberOfDeployments; i++ {
		go func(i int) {
			doneCh <- randomizer.RunFlow(i, config.UsingLegacyManifest)
		}(i)
	}

	for i := 0; i < config.NumberOfDeployments; i++ {
		err := <-doneCh
		if err != nil {
			panic(err)
		}
	}

	println("Done!")
}
