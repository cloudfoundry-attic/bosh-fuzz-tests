package clirunner

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type uaacRunner struct {
	cmd       boshsys.Command
	cmdRunner boshsys.CmdRunner
	env       string
}

func NewUaacRunner(cmd boshsys.Command, cmdRunner boshsys.CmdRunner) Runner {
	return &uaacRunner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
	}
}

func (r *uaacRunner) SetEnv(envName string) {
	r.env = envName
}

func (r *uaacRunner) RunInDirWithArgs(dir string, args ...string) error {
	cmd := r.cliCommand(args...)
	cmd.WorkingDir = dir
	_, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *uaacRunner) RunWithArgs(args ...string) error {
	_, err := r.RunWithOutput(args...)
	return err
}

func (r *uaacRunner) RunWithOutput(args ...string) (string, error) {
	stdOut, _, _, err := r.cmdRunner.RunComplexCommand(r.cliCommand(args...))
	if err != nil {
		return stdOut, err
	}

	return stdOut, nil
}

func (r *uaacRunner) cliCommand(args ...string) boshsys.Command {
	cmd := r.cmd
	cmd.Args = append(cmd.Args, args...)
	return cmd
}
