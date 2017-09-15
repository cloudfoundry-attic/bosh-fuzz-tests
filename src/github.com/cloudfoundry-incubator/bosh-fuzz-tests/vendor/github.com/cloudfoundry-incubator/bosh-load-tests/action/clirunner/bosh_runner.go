package clirunner

import (
	"os"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type boshRunner struct {
	cmd       boshsys.Command
	cmdRunner boshsys.CmdRunner
	fs        boshsys.FileSystem
	env       string
}

func NewBoshRunner(cmd boshsys.Command, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) Runner {
	return &boshRunner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
		fs:        fs,
	}
}

func (r *boshRunner) SetEnv(envName string) {
	r.env = envName
}

func (r *boshRunner) RunInDirWithArgs(dir string, args ...string) error {
	cmd := r.cliCommand(args...)
	cmd.WorkingDir = dir
	_, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *boshRunner) RunWithArgs(args ...string) error {
	_, err := r.RunWithOutput(args...)
	return err
}

func (r *boshRunner) RunWithOutput(args ...string) (string, error) {
	stdOut, _, _, err := r.cmdRunner.RunComplexCommand(r.cliCommand(args...))
	if err != nil {
		return stdOut, err
	}

	return stdOut, nil
}

func (r *boshRunner) cliCommand(args ...string) boshsys.Command {
	cmd := r.cmd
	if r.env != "" {
		cmd.Args = append(cmd.Args, "-e", r.env)

		if _, found := os.LookupEnv("BOSH_CA_CERT"); !found {
			cmd.Args = append(cmd.Args, "--ca-cert", "/tmp/cert")
		}
		if _, found := os.LookupEnv("BOSH_CLIENT"); !found {
			cmd.Args = append(cmd.Args, "--client", "test", "--client-secret", "secret")
		}
	}
	cmd.Args = append(cmd.Args, "-n", "--tty")
	cmd.Args = append(cmd.Args, args...)
	return cmd
}