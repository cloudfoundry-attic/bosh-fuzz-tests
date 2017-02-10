package clirunner

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Runner interface {
	SetEnv(envName string)
	RunInDirWithArgs(dir string, args ...string) error
	RunWithArgs(args ...string) error
	RunWithOutput(args ...string) (string, error)
}

type runner struct {
	cmd       boshsys.Command
	cmdRunner boshsys.CmdRunner
	fs        boshsys.FileSystem
	env       string
}

func NewRunner(cmd boshsys.Command, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) Runner {
	return &runner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
		fs:        fs,
	}
}

func (r *runner) SetEnv(envName string) {
	r.env = envName
}

func (r *runner) RunInDirWithArgs(dir string, args ...string) error {
	cmd := r.cliCommand(args...)
	cmd.WorkingDir = dir
	_, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *runner) RunWithArgs(args ...string) error {
	_, err := r.RunWithOutput(args...)
	return err
}

func (r *runner) RunWithOutput(args ...string) (string, error) {
	stdOut, _, _, err := r.cmdRunner.RunComplexCommand(r.cliCommand(args...))
	if err != nil {
		return stdOut, err
	}

	return stdOut, nil
}

func (r *runner) cliCommand(args ...string) boshsys.Command {
	cmd := r.cmd
	if r.env != "" {
		cmd.Args = append(cmd.Args, "-e", r.env)
	}
	cmd.Args = append(cmd.Args, "--ca-cert", "/tmp/cert")
	cmd.Args = append(cmd.Args, "-n", "--tty")
	cmd.Args = append(cmd.Args, "--client", "admin", "--client-secret", "admin")
	cmd.Args = append(cmd.Args, args...)

	return cmd
}
