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

func (s *start) Execute() error {
	s.cliRunner.SetEnv(s.directorInfo.URL)

	deployWrapper := NewDeployWrapper(s.cliRunner)
	_, err := deployWrapper.RunWithDebug("-d", s.deploymentName, "start", "simple/0")
	if err != nil {
		return err
	}

	return nil
}
