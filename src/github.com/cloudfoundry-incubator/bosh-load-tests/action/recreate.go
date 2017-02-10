package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type recreate struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
}

func NewRecreate(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner, fs boshsys.FileSystem) *recreate {
	return &recreate{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (r *recreate) Execute() error {
	r.cliRunner.SetEnv(r.directorInfo.URL)

	deployWrapper := NewDeployWrapper(r.cliRunner)
	_, err := deployWrapper.RunWithDebug("-d", r.deploymentName, "recreate", "simple/0")
	if err != nil {
		return err
	}

	return nil
}
