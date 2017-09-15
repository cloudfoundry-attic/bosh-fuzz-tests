package clirunner

import (
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"
	"github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Runner interface {
	SetEnv(envName string)
	RunInDirWithArgs(dir string, args ...string) error
	RunWithArgs(args ...string) error
	RunWithOutput(args ...string) (string, error)
}

type Factory interface {
	Create(cmd string) (Runner, error)
}

type factory struct {
	boshCliPath string
	cmdRunner   boshsys.CmdRunner
	fs          boshsys.FileSystem
}

func NewFactory(cliPath string, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) *factory {
	return &factory{
		boshCliPath: cliPath,
		cmdRunner:   cmdRunner,
		fs:          fs,
	}
}

func (f *factory) Create(cmd string) (Runner, error) {
	switch cmd {
	case "bosh":
		return NewBoshRunner(bltcom.CreateCommand(f.boshCliPath), f.cmdRunner, f.fs), nil
	case "uaac":
		return NewUaacRunner(bltcom.CreateCommand("uaac"), f.cmdRunner), nil
	}

	return nil, errors.Errorf("Unable to create runner for type '%s'", cmd)
}
