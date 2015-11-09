package main

import (
	"os"
	"os/signal"
	"syscall"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	// bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	// bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	// bltflow "github.com/cloudfoundry-incubator/bosh-load-tests/flow"

	bltenv "github.com/cloudfoundry-incubator/bosh-load-tests/environment"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		println("Usage: bft path/to/config.json [path/to/state.json]")
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
	environmentProvider := bltenv.NewProvider(config, fs, cmdRunner, assetsProvider)
	environment := environmentProvider.Get()
	err = environment.Setup()
	if err != nil {
		panic(err)
	}
	defer environment.Shutdown()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		environment.Shutdown()
		os.Exit(1)
	}()
}
