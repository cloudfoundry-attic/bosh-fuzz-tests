package clirunner

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Runner interface {
	Configure() error
	Clean() error
	TargetAndLogin(target string) error
	RunInDirWithArgs(dir string, args ...string) error
	RunWithArgs(args ...string) error
	RunWithOutput(args ...string) (string, error)
}

type runner struct {
	configPath string
	cmd        boshsys.Command
	cmdRunner  boshsys.CmdRunner
	fs         boshsys.FileSystem
}

func NewRunner(cmd boshsys.Command, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) Runner {
	return &runner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
		fs:        fs,
	}
}

func (r *runner) Configure() error {
	configFile, err := r.fs.TempFile("bosh-config")
	if err != nil {
		return err
	}
	r.configPath = configFile.Name()
	return nil
}

func (r *runner) Clean() error {
	if r.configPath == "" {
		return nil
	}

	return r.fs.RemoveAll(r.configPath)
}

func (r *runner) TargetAndLogin(target string) error {
	err := r.RunWithArgs("target", target)
	if err != nil {
		return err
	}

	err = r.RunWithArgs("login", "admin", "admin")
	if err != nil {
		return err
	}

	return nil
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
	cmd.Args = append(cmd.Args, "-n", "-c", r.configPath)
	cmd.Args = append(cmd.Args, args...)

	return cmd
}
