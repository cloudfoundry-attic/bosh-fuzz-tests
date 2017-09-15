package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deleteDeployment struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
}

func NewDeleteDeployment(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner, fs boshsys.FileSystem) *deleteDeployment {
	return &deleteDeployment{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (d *deleteDeployment) Execute() error {
	d.cliRunner.SetEnv(d.directorInfo.URL)

	deployWrapper := NewDeployWrapper(d.cliRunner)
	_, err := deployWrapper.RunWithDebug("-d", d.deploymentName, "delete-deployment")
	if err != nil {
		return err
	}

	return nil
}
