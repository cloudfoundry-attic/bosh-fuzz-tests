package clirunner

import (
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Factory interface {
	Create() Runner
}

type factory struct {
	cliPath   string
	cmdRunner boshsys.CmdRunner
	fs        boshsys.FileSystem
}

func NewFactory(cliPath string, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) *factory {
	return &factory{
		cliPath:   cliPath,
		cmdRunner: cmdRunner,
		fs:        fs,
	}
}

func (f *factory) Create() Runner {
	cmd := bltcom.CreateCommand(f.cliPath)

	return NewRunner(cmd, f.cmdRunner, f.fs)
}
