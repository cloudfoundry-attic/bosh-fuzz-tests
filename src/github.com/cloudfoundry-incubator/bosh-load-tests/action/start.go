package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type start struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
}

func NewStart(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner, fs boshsys.FileSystem) *start {
	return &start{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (r *start) Execute() error {
	manifestPath, err := r.fs.TempFile("manifest-test")
	if err != nil {
		return err
	}
	defer r.fs.RemoveAll(manifestPath.Name())

	err = r.cliRunner.RunWithArgs("download", "manifest", r.deploymentName, manifestPath.Name())
	if err != nil {
		return err
	}

	err = r.cliRunner.RunWithArgs("deployment", manifestPath.Name())
	if err != nil {
		return err
	}

	deployWrapper := NewDeployWrapper(r.cliRunner)
	_, err = deployWrapper.RunWithDebug("start", "simple", "0")
	if err != nil {
		return err
	}

	return nil
}
